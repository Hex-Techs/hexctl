package kc

import (
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// KubeConfig kubeconfig object
type KubeConfig struct {
	Kind           string     `json:"kind,omitempty"`
	APIVersion     string     `json:"apiVersion,omitempty"`
	Clusters       []Cluster  `json:"clusters"`
	AuthInfos      []AuthInfo `json:"users"`
	Contexts       []Context  `json:"contexts"`
	CurrentContext string     `json:"current-context"`
}

// Cluster k8s cluster object
type Cluster struct {
	Name    string                `json:"name"`
	Cluster *clientcmdapi.Cluster `json:"cluster"`
}

// AuthInfo k8sconfig user object
type AuthInfo struct {
	Name string                 `json:"name"`
	User *clientcmdapi.AuthInfo `json:"user"`
}

// Context k8s context object
type Context struct {
	Name    string                `json:"name"`
	Context *clientcmdapi.Context `json:"context"`
}
