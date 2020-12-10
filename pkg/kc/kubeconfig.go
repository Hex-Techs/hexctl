package kc

import (
	"encoding/json"
	"os"

	"github.com/Hex-Techs/n/pkg/output"
	"github.com/Hex-Techs/n/pkg/utils"
	"github.com/ghodss/yaml"
)

func SwitchConfig(kubeconfig string) {
	initKubeconfig(kubeconfig)
	cfg := GetKubeConfig()
	i := SelectUI(cfg.Clusters)
	cfg.CurrentContext = cfg.Clusters[i].Name
	d, err := json.Marshal(cfg)
	if err != nil {
		output.Errorln(err)
		os.Exit(1)
	}
	c, err := yaml.JSONToYAML(d)
	if err != nil {
		output.Errorln(err)
		os.Exit(1)
	}
	utils.Write(string(c), Kubeconfig)
}

func Show(kubeconfig string) {
	initKubeconfig(kubeconfig)
	cfg := GetKubeConfig()
	output.Successln(cfg.CurrentContext)
}
