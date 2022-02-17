// Package kc manage local kubeconfig file
package kc

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"time"

	"github.com/Hex-Techs/hexctl/pkg/common/display"
	"github.com/gookit/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// Ls show the context list
func Ls(kubeconfig string) {
	d := defaultKubeConfig(kubeconfig)
	cfg, err := clientcmd.LoadFromFile(d)
	cobra.CheckErr(err)

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name", "Server"})
	for i, c := range cfg.Clusters {
		if cfg.CurrentContext == i {
			t.AppendRow([]interface{}{"* " + i, c.Server})
		} else {
			t.AppendRow([]interface{}{"  " + i, c.Server})
		}
	}
	t.Render()
}

// Switch switch context for kubeconfig
func Switch(kubeconfig, cluster string, ns bool) {
	d := defaultKubeConfig(kubeconfig)
	cfg, err := clientcmd.LoadFromFile(d)
	cobra.CheckErr(err)
	var context string
	if len(cluster) == 0 {
		items := []string{}
		for i := range cfg.Contexts {
			items = append(items, i)
		}
		context = display.Select("Select the kubeconfig Context", len(items), items)
		if context == "" {
			return
		}
	} else {
		for idx := range cfg.Contexts {
			if idx == cluster {
				context = idx
			}
		}
	}
	cfg.CurrentContext = context
	err = clientcmd.WriteToFile(*cfg, d)
	cobra.CheckErr(err)
	color.Green.Println("switch context to", context)
	if ns {
		Namespace(kubeconfig, "")
	}
	Show(kubeconfig)
}

// Show show the context
func Show(kubeconfig string) {
	d := defaultKubeConfig(kubeconfig)
	cfg, err := clientcmd.LoadFromFile(d)
	cobra.CheckErr(err)
	var ns string
	for i, c := range cfg.Contexts {
		if i == cfg.CurrentContext {
			ns = c.Namespace
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
	cfg, err := clientcmd.LoadFromFile(d)
	cobra.CheckErr(err)

	items := []string{}
	for i := range cfg.Contexts {
		items = append(items, i)
	}
	context := display.Select("Select the kubeconfig Context which do you want to delete", len(items), items)
	if len(context) == 0 {
		return
	}
	if display.Confirm(fmt.Sprintf("Do you want to Delete %s", context)) {
		ctx := cfg.Contexts[context]
		delete(cfg.AuthInfos, ctx.AuthInfo)
		delete(cfg.Clusters, ctx.Cluster)
		delete(cfg.Contexts, context)
		if cfg.CurrentContext == context {
			for i := range cfg.Contexts {
				cfg.CurrentContext = i
				break
			}
		}
		err = clientcmd.WriteToFile(*cfg, d)
		cobra.CheckErr(err)
		color.Green.Printf("Delete %s success\n", context)
	}
}

// Namespace switch default work namespace
func Namespace(kubeconfig, namespace string) {
	var ns string
	d := defaultKubeConfig(kubeconfig)
	if len(namespace) == 0 {
		config, err := clientcmd.BuildConfigFromFlags("", d)
		if err != nil {
			cobra.CheckErr(err)
		}
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			cobra.CheckErr(err)
		}
		var timeout int64 = 5
		nss, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{TimeoutSeconds: &timeout})
		if err != nil {
			cobra.CheckErr(err)
		}
		items := []string{}
		for _, v := range nss.Items {
			items = append(items, v.Name)
		}
		ns = display.Select("Select the Namespace", len(items), items)
		if len(ns) == 0 {
			return
		}
	} else {
		ns = namespace
	}
	c, err := clientcmd.LoadFromFile(d)
	cobra.CheckErr(err)
	for i, v := range c.Contexts {
		if i == c.CurrentContext {
			v.Namespace = ns
		}
	}
	err = clientcmd.WriteToFile(*c, d)
	cobra.CheckErr(err)
}

// Merge merge kubeconfig
func Merge(src, kubeconfig string) {
	d := defaultKubeConfig(kubeconfig)
	srcConfig, err := clientcmd.LoadFromFile(src)
	cobra.CheckErr(err)
	dstConfig, err := clientcmd.LoadFromFile(d)
	cobra.CheckErr(err)

	// generate a context name by time
	tmp := fmt.Sprint(time.Now().Unix())

	for idx, v := range srcConfig.Contexts {
		name := fmt.Sprintf("%s-%s", idx, tmp)
		dstConfig.AuthInfos[name] = srcConfig.AuthInfos[v.AuthInfo]
		dstConfig.Clusters[name] = srcConfig.Clusters[v.Cluster]
		v.Cluster = name
		v.AuthInfo = name
		dstConfig.Contexts[name] = v
	}
	err = clientcmd.WriteToFile(*dstConfig, d)
	cobra.CheckErr(err)
}

// GetContext get a context want
func GetContext(kubeconfig string) {
	d := defaultKubeConfig(kubeconfig)
	cfg, err := clientcmd.LoadFromFile(d)
	cobra.CheckErr(err)

	items := []string{}
	for i := range cfg.Contexts {
		items = append(items, i)
	}
	context := display.Select("Select the kubeconfig Context which do you want to manifest", len(items), items)
	if len(context) == 0 {
		return
	}
	get := clientcmdapi.NewConfig()
	a := cfg.Contexts[context].AuthInfo
	c := cfg.Contexts[context].Cluster

	get.Contexts[context] = cfg.Contexts[context]
	get.AuthInfos[a] = cfg.AuthInfos[a]
	get.Clusters[c] = cfg.Clusters[c]
	b, err := clientcmd.Write(*get)
	cobra.CheckErr(err)
	fmt.Printf("\n%s\n", string(b))
}

func defaultKubeConfig(kubeconfig string) string {
	if kubeconfig == "" {
		u, err := user.Current()
		cobra.CheckErr(err)
		return fmt.Sprintf("%s/.kube/config", u.HomeDir)
	}
	return kubeconfig
}
