package kc

import (
	"fmt"
	"os"
	"os/user"

	"github.com/Fize/n/pkg/output"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var DefaultKubeconfig string

// init default kubeconfig
func init() {
	u, err := user.Current()
	if err != nil {
		output.Errorln(err)
		os.Exit(1)
	}
	DefaultKubeconfig = fmt.Sprintf("%s/.kube/config", u.HomeDir)
}

type KubeConfig struct {
	Kind           string     `json:"kind,omitempty"`
	APIVersion     string     `json:"apiVersion,omitempty"`
	Clusters       []Cluster  `json:"clusters"`
	AuthInfos      []AuthInfo `json:"users"`
	Contexts       []Context  `json:"contexts"`
	CurrentContext string     `json:"current-context"`
}

type Cluster struct {
	Name    string                `json:"name"`
	Cluster *clientcmdapi.Cluster `json:"cluster"`
}

type AuthInfo struct {
	Name string                 `json:"name"`
	User *clientcmdapi.AuthInfo `json:"user"`
}

type Context struct {
	Name    string                `json:"name"`
	Context *clientcmdapi.Context `json:"context"`
}
