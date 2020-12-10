package kc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/Hex-Techs/n/pkg/output"
	"github.com/Hex-Techs/n/pkg/utils"
	"github.com/manifoldco/promptui"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
)

func SelectUI(clusters []Cluster) int {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F63C {{ .Name | cyan }}",
		Inactive: "  {{ .Name | cyan }} ",
		Selected: "\U0001F638 {{ .Name | red | cyan }}",
	}

	searcher := func(input string, index int) bool {
		cluster := clusters[index]
		name := strings.Replace(strings.ToLower(cluster.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "Select Kube Context",
		Items:     clusters,
		Templates: templates,
		Size:      4,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()
	if err != nil {
		if err.Error() != "^C" {
			output.Errorf("Prompt failed %v\n", err)
			os.Exit(1)
		}
	}
	return i
}

func GetKubeConfig() KubeConfig {
	f := utils.Read(Kubeconfig)
	cfg := KubeConfig{}
	obj := &unstructured.Unstructured{}

	dec := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	_, _, err := dec.Decode([]byte(f), nil, obj)

	s := bytes.NewBuffer(nil)
	enc := json.NewEncoder(s)
	enc.Encode(obj)
	err = json.Unmarshal([]byte(s.String()), &cfg)
	if err != nil {
		output.Errorln(err)
	}
	return cfg
}

func initKubeconfig(kubeconfig string) {
	u, err := user.Current()
	if err != nil {
		output.Errorln(err)
		os.Exit(1)
	}
	if kubeconfig == "" {
		Kubeconfig = fmt.Sprintf("%s/.kube/config", u.HomeDir)
	} else {
		Kubeconfig = kubeconfig
	}
}
