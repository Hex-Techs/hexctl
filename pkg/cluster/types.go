package cluster

import (
	"github.com/Fize/n/pkg/utils"
	"k8s.io/kube-proxy/config/v1alpha1"
	"k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta1"
)

const (
	DefaultUser             = "root"
	DefaultPort             = "22"
	DefaultServicePortRange = "30000-32767"
)

// ClusterCommand n cluster command struct
type ClusterCommand struct {
	Master           string
	Node             []string
	CertHost         []string
	Repo             string
	Volume           string
	APIServer        string
	Token            string
	UnSafe           bool
	CAHash           string
	Ignore           []string
	PodCIDR          string
	ServiceCIDR      string
	IPVS             bool
	ServicePortRange string
	Password         string
	Key              string
	Iface            string
	ControlPlane     bool
	Port             string
}

// KubernetesCluster kubernetes cluster object
type KubernetesCluster struct {
	MasterOption    *utils.RemoteOption
	NodeOption      []*utils.RemoteOption
	ClusterConfig   *v1beta1.ClusterConfiguration
	KubeProxyConfig *v1alpha1.KubeProxyConfiguration
}
