// Package kc manage local kubeconfig file
package kubeconfig

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"os/user"
	"path"
	"strings"
	"time"

	"github.com/Hex-Techs/hexctl/pkg/common/display"
	"github.com/jedib0t/go-pretty/v6/table"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

const (
	// 默认 kubeconfig 文件
	defulatKubeConfigPath = "/.kube/config"
)

// NewKCMgr create a new kubeconfig manager, p is the special path of kubeconfig file
func NewKCMgr(p string) (*kubeconfigMgr, error) {
	u, err := user.Current()
	if err != nil {
		return nil, err
	}
	kcm := &kubeconfigMgr{}
	if p == "" {
		p = path.Join(u.HomeDir, defulatKubeConfigPath)
	}
	kcm.path = p
	kcm.config, err = clientcmd.LoadFromFile(p)
	return kcm, err
}

type kubeconfigMgr struct {
	// kubeconfig file path
	path string
	// cluster context
	clusterSelector string
	// which namespace to use for the current context
	namespaceSelector string
	// config content
	config *clientcmdapi.Config
}

// MergeContext merge kubeconfig file
func (k *kubeconfigMgr) MergeContext(src string) error {
	srcConfig, err := clientcmd.LoadFromFile(src)
	if err != nil {
		return err
	}
	// context name has unix timestamp
	tmp := fmt.Sprint(time.Now().Unix())
	for idx, v := range srcConfig.Contexts {
		name := fmt.Sprintf("%s-%s", idx, tmp)
		k.config.AuthInfos[name] = srcConfig.AuthInfos[v.AuthInfo]
		k.config.Clusters[name] = srcConfig.Clusters[v.Cluster]
		v.Cluster = name
		v.AuthInfo = name
		k.config.Contexts[name] = v
	}
	return clientcmd.WriteToFile(*k.config, k.path)
}

// DeleteContext delete a existed context
func (k *kubeconfigMgr) DeleteContext() error {
	items := k.readContexts()
	d := display.NewTerminalDisplay("选择要删除的集群配置", len(items), items...)
	context := d.Select()
	if len(context) == 0 {
		return nil
	}
	d = display.NewTerminalDisplay(fmt.Sprintf("确定删除 %s 集群配置?", context), 0)
	if d.Confirm() {
		ctx := k.config.Contexts[context]
		delete(k.config.AuthInfos, ctx.AuthInfo)
		delete(k.config.Clusters, ctx.Cluster)
		delete(k.config.Contexts, context)
		if k.config.CurrentContext == context {
			for i := range k.config.Contexts {
				k.config.CurrentContext = i
				break
			}
		}
		err := clientcmd.WriteToFile(*k.config, k.path)
		if err != nil {
			return err
		}
		fmt.Printf("删除 %s 配置成功\n", context)
	}
	return nil
}

// RenameContext rename a existed context
func (k *kubeconfigMgr) RenameContext() error {
	items := k.readContexts()
	d := display.NewTerminalDisplay("选择要修改名称的 context", len(items), items...)
	context := d.Select()
	if len(context) == 0 {
		return nil
	}
	inputReader := bufio.NewReader(os.Stdin)
	fmt.Printf("输入 context 名称: ")
	name, _ := inputReader.ReadString('\n')
	name = strings.Replace(name, "\n", "", -1)
	ctx := k.config.Contexts[context]
	cls := k.config.Clusters[ctx.Cluster]
	delete(k.config.Contexts, context)
	delete(k.config.Clusters, ctx.Cluster)
	ctx.Cluster = name
	k.config.Clusters[name] = cls
	k.config.Contexts[name] = ctx
	if k.config.CurrentContext == context {
		k.config.CurrentContext = name
	}
	err := clientcmd.WriteToFile(*k.config, k.path)
	if err != nil {
		return err
	}
	fmt.Printf("成功将 context 名称从 %s 修改到 %s\n", context, name)
	return nil
}

// GetContext get a context from kubeconfig
// bf is the byte format output, used by []byte
func (k *kubeconfigMgr) GetContext(dst string, bf bool) error {
	items := k.readContexts()
	d := display.NewTerminalDisplay("选择要查看的集群配置", len(items), items...)
	context := d.Select()
	if len(context) == 0 {
		return nil
	}
	get := clientcmdapi.NewConfig()
	a := k.config.Contexts[context].AuthInfo
	c := k.config.Contexts[context].Cluster
	get.Contexts[context] = k.config.Contexts[context]
	get.AuthInfos[a] = k.config.AuthInfos[a]
	get.Clusters[c] = k.config.Clusters[c]
	get.CurrentContext = context
	if dst != "" {
		return clientcmd.WriteToFile(*get, dst)
	}
	b, err := clientcmd.Write(*get)
	if err != nil {
		return err
	}
	if bf {
		fmt.Printf("\n%s\n", base64.StdEncoding.EncodeToString(b))
		return nil
	}
	fmt.Printf("\n%s\n", string(b))
	return nil
}

// ShowCurrentContext show current context
func (k *kubeconfigMgr) ShowCurrentContext() {
	ns := "default"
	for i, c := range k.config.Contexts {
		if i == k.config.CurrentContext {
			ns = c.Namespace
			break
		}
	}
	d := display.NewTerminalDisplay("", 0)
	header := table.Row{"集群", "命名空间"}
	d.Table(header, []interface{}{k.config.CurrentContext, ns})
}

// ListContext show the context list
func (k *kubeconfigMgr) ListContext() {
	d := display.NewTerminalDisplay("", 0)
	header := table.Row{"当前选中", "集群", "API-Server 地址"}
	rows := [][]interface{}{}
	for i, c := range k.config.Clusters {
		current := ""
		if k.config.CurrentContext == i {
			current = "*"
		}
		rows = append(rows, []interface{}{current, i, c.Server})
	}
	d.Table(header, rows...)
}

// Switch switch context for kubeconfig
func (k *kubeconfigMgr) SwitchContext(cascade bool) error {
	items := k.readContexts()
	d := display.NewTerminalDisplay("选择集群", len(items), items...)
	k.clusterSelector = d.Select()
	if len(k.clusterSelector) == 0 {
		return nil
	}
	k.config.CurrentContext = k.clusterSelector
	if err := clientcmd.WriteToFile(*k.config, k.path); err != nil {
		return err
	}
	if cascade {
		return k.SwitchNamespace()
	}
	return nil
}

// SwitchNamespace 切换默认的namespace
func (k *kubeconfigMgr) SwitchNamespace() error {
	items, err := k.getNamespacesFromCluster()
	if err != nil {
		return err
	}
	d := display.NewTerminalDisplay("选择命名空间", len(items), items...)
	k.namespaceSelector = d.Select()
	if len(k.namespaceSelector) == 0 {
		return nil
	}
	for i, v := range k.config.Contexts {
		if i == k.config.CurrentContext {
			v.Namespace = k.namespaceSelector
		}
	}
	return clientcmd.WriteToFile(*k.config, k.path)
}

// read context and return the context name
func (k *kubeconfigMgr) readContexts() []string {
	items := []string{}
	for i := range k.config.Contexts {
		items = append(items, i)
	}
	return items
}

// 根据当前 cluster 获取 namespace
func (k *kubeconfigMgr) getNamespacesFromCluster() ([]string, error) {
	config, err := clientcmd.BuildConfigFromFlags("", k.path)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	var timeout int64 = 5
	nss, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{TimeoutSeconds: &timeout})
	if err != nil {
		return nil, err
	}
	items := []string{}
	for _, v := range nss.Items {
		items = append(items, v.Name)
	}
	return items, nil
}
