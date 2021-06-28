// Package kc manage local kubeconfig file
package kc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/user"
	"time"

	"github.com/Hex-Techs/hexctl/pkg/common/exec"
	"github.com/Hex-Techs/hexctl/pkg/common/file"
	"github.com/Hex-Techs/hexctl/pkg/display"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	kyaml "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
)

// Switch switch context for kubeconfig
func Switch(kubeconfig string) {
	d := defaultKubeConfig(kubeconfig)
	cfg := getContent(d)

	items := []string{}
	for _, v := range cfg.Clusters {
		items = append(items, v.Name)
	}
	context := display.SelectUI("Select the kubeconfig Context", items)
	if len(context) != 0 {
		cfg.CurrentContext = context
		file.Write(convert(cfg), d)
	}
}

// Show show the context
func Show(kubeconfig string) {
	d := defaultKubeConfig(kubeconfig)
	cmd := fmt.Sprintf("kubectl config current-context --kubeconfig %s", d)
	exec.RunCommand(cmd)
}

// Namespace switch default work namespace
func Namespace(kubeconfig, namespace string) {
	d := defaultKubeConfig(kubeconfig)
	// use kubectl switch work namespace
	cmd := fmt.Sprintf("kubectl config set-context --kubeconfig %s --current --namespace %s", d, namespace)
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
		name = fmt.Sprintf("%s%d", name, index)
		v.Name = name

		u := v.Context.AuthInfo
		c := v.Context.Cluster
		for _, av := range srcConfig.AuthInfos {
			if av.Name == u {
				av.Name = name
				dstConfig.AuthInfos = append(dstConfig.AuthInfos, av)
			}
		}
		for _, cv := range srcConfig.Clusters {
			if cv.Name == c {
				cv.Name = name
				dstConfig.Clusters = append(dstConfig.Clusters, cv)
			}
		}
		v.Context.Cluster = name
		v.Context.AuthInfo = name
		dstConfig.Contexts = append(dstConfig.Contexts, v)
	}
	file.Write(convert(dstConfig), d)
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

	dec := kyaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	_, _, err := dec.Decode([]byte(f), nil, obj)

	s := bytes.NewBuffer(nil)
	enc := json.NewEncoder(s)
	enc.Encode(obj)
	err = json.Unmarshal([]byte(s.String()), &cfg)
	if err != nil {
		display.Errorln(err)
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
