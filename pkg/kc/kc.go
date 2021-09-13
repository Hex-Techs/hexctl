// Package kc manage local kubeconfig file
package kc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"time"

	"github.com/Hex-Techs/hexctl/pkg/common/display"
	"github.com/Hex-Techs/hexctl/pkg/common/exec"
	"github.com/Hex-Techs/hexctl/pkg/common/file"
	"github.com/ghodss/yaml"
	"github.com/gookit/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	kyaml "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Ls show the context list
func Ls(kubeconfig string) {
	d := defaultKubeConfig(kubeconfig)
	cfg := getContent(d)

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name", "Endpoint"})
	for _, c := range cfg.Clusters {
		if cfg.CurrentContext == c.Name {
			t.AppendRow([]interface{}{"* " + c.Name, c.Cluster.Server})
		} else {
			t.AppendRow([]interface{}{"  " + c.Name, c.Cluster.Server})
		}
	}
	t.Render()
}

// Switch switch context for kubeconfig
func Switch(kubeconfig string) {
	d := defaultKubeConfig(kubeconfig)
	cfg := getContent(d)

	items := []string{}
	for _, v := range cfg.Contexts {
		items = append(items, v.Name)
	}
	context := display.SelectUI("Select the kubeconfig Context", len(items), items)
	if context == "" {
		return
	}
	cfg.CurrentContext = context
	file.Write(convert(cfg), d)
}

// Show show the context
func Show(kubeconfig string) {
	d := defaultKubeConfig(kubeconfig)
	// cmd := fmt.Sprintf("kubectl config current-context --kubeconfig %s", d)
	// exec.RunCommand(cmd)
	cfg := getContent(d)
	var ns string
	for _, c := range cfg.Contexts {
		if c.Name == cfg.CurrentContext {
			ns = c.Context.Namespace
		}
	}
	if ns == "" {
		ns = "default"
	}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Clustr", "Namespace"})
	t.AppendRow([]interface{}{cfg.CurrentContext, ns})
	t.Render()
}

// Delete delete a context from kubeconfig
func Delete(kubeconfig string) {
	d := defaultKubeConfig(kubeconfig)
	cfg := getContent(d)

	items := []string{}
	for _, v := range cfg.Contexts {
		items = append(items, v.Name)
	}
	context := display.SelectUI("Select the kubeconfig Context which do you want delete", len(items), items)
	ctx, nctx := generateContext(context, cfg.Contexts)
	_, u := generateAuth(ctx.Context.AuthInfo, cfg.AuthInfos)
	_, c := generateCluster(ctx.Context.Cluster, cfg.Clusters)
	cfg.Contexts = nctx
	cfg.AuthInfos = u
	cfg.Clusters = c
	file.Write(convert(cfg), d)
}

func generateContext(name string, cts []Context) (*Context, []Context) {
	r := cts
	for k, v := range cts {
		if v.Name == name {
			r = append(r[:k], r[k+1:]...)
			return &v, r
		}
	}
	return nil, cts
}

func generateAuth(name string, auth []AuthInfo) (*AuthInfo, []AuthInfo) {
	r := auth
	for k, v := range auth {
		if v.Name == name {
			r = append(r[:k], r[k+1:]...)
			return &v, r
		}
	}
	return nil, nil
}

func generateCluster(name string, cluster []Cluster) (*Cluster, []Cluster) {
	r := cluster
	for k, v := range cluster {
		if v.Name == name {
			r = append(r[:k], r[k+1:]...)
			return &v, r
		}
	}
	return nil, nil
}

// Namespace switch default work namespace
func Namespace(kubeconfig, namespace string) {
	if len(namespace) == 0 {
		namespace = "default"
	}
	d := defaultKubeConfig(kubeconfig)

	config, err := clientcmd.BuildConfigFromFlags("", d)
	if err != nil {
		cobra.CheckErr(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		cobra.CheckErr(err)
	}
	nss, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		cobra.CheckErr(err)
	}

	items := []string{}
	for _, v := range nss.Items {
		items = append(items, v.Name)
	}
	ns := display.SelectUI("Select the Namespace", len(items), items)
	if ns == "" {
		ns = namespace
	}
	// use kubectl switch work namespace
	cmd := fmt.Sprintf("kubectl config set-context --kubeconfig %s --current --namespace %s", d, ns)
	exec.RunCommand(cmd)
}

// Merge merge kubeconfig
func Merge(src, kubeconfig string) {
	d := defaultKubeConfig(kubeconfig)
	srcConfig := getContent(src)
	dstConfig := getContent(d)

	// generate a context name by time
	name := fmt.Sprint(time.Now().Unix())

	for index, v := range srcConfig.Contexts {
		v.Name = fmt.Sprintf("%s-%s%d", v.Name, name, index)
		u := v.Context.AuthInfo
		c := v.Context.Cluster
		for _, av := range srcConfig.AuthInfos {
			if av.Name == u {
				av.Name = v.Name
				dstConfig.AuthInfos = append(dstConfig.AuthInfos, av)
			}
		}
		for _, cv := range srcConfig.Clusters {
			if cv.Name == c {
				cv.Name = v.Name
				dstConfig.Clusters = append(dstConfig.Clusters, cv)
			}
		}
		v.Context.Cluster = v.Name
		v.Context.AuthInfo = v.Name
		dstConfig.Contexts = append(dstConfig.Contexts, v)
	}
	file.Write(convert(dstConfig), d)
}

// GetContext get a context want
func GetContext(kubeconfig string) {
	d := defaultKubeConfig(kubeconfig)
	cfg := getContent(d)

	items := []string{}
	for _, v := range cfg.Contexts {
		items = append(items, v.Name)
	}
	context := display.SelectUI("Select the kubeconfig Context which do you want to manifest", len(items), items)
	ctx, _ := generateContext(context, cfg.Contexts)
	u, _ := generateAuth(ctx.Context.AuthInfo, cfg.AuthInfos)
	c, _ := generateCluster(ctx.Context.Cluster, cfg.Clusters)
	var get KubeConfig
	get.Contexts = []Context{*ctx}
	get.AuthInfos = []AuthInfo{*u}
	get.Clusters = []Cluster{*c}
	get.CurrentContext = ctx.Name
	get.APIVersion = "v1"
	get.Kind = "Config"
	fmt.Println(convert(&get))
}

func defaultKubeConfig(kubeconfig string) string {
	if kubeconfig == "" {
		u, err := user.Current()
		cobra.CheckErr(err)
		return fmt.Sprintf("%s/.kube/config", u.HomeDir)
	}
	return kubeconfig
}

func getContent(kubeconfig string) *KubeConfig {
	f := file.Read(kubeconfig)
	cfg := KubeConfig{}
	obj := &unstructured.Unstructured{}
	var err error
	dec := kyaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	dec.Decode([]byte(f), nil, obj)

	s := bytes.NewBuffer(nil)
	enc := json.NewEncoder(s)
	enc.Encode(obj)
	err = json.Unmarshal(s.Bytes(), &cfg)
	if err != nil {
		color.Red.Println(err)
	}
	return &cfg
}

func convert(cfg *KubeConfig) string {
	d, err := json.Marshal(cfg)
	cobra.CheckErr(err)
	c, err := yaml.JSONToYAML(d)
	cobra.CheckErr(err)
	return string(c)
}
