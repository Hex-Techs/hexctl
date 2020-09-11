package cluster

import (
	"fmt"

	"github.com/Fize/n/pkg/output"
	"github.com/Fize/n/pkg/utils"
)

const (
	defaultUser             = "root"
	defaultPort             = "22"
	defaultServicePortRange = "30000-32767"
)

// ClusterCommand n cluster command struct
type ClusterCommand struct {
	Endpoint         string
	Token            string
	UnSafe           bool
	CAHash           string
	Ignore           []string
	PodCIDR          string
	ServiceCIDR      string
	IPVS             bool
	ServicePortRange string
	Repo             string
	Password         string
	Key              string
	Iface            string
	Port             string
	KubeVersion      string
	DockerVersion    string
	Type             string
	IP               string
	CN               bool
}

// KubernetesCluster kubernetes cluster object
type KubernetesCluster struct {
	Option *utils.RemoteOption
	ncmd   ClusterCommand
	Type   string
	node   string
	method string
}

func NewKubernetesCluster(ncmd *ClusterCommand, tp, node, method string) *KubernetesCluster {
	kc := &KubernetesCluster{
		ncmd:   *ncmd,
		Type:   tp,
		node:   node,
		method: method,
	}
	var err error
	if tp == "local" {
		ncmd.Iface = "eth1"
	}
	if ncmd.IP != "" {
		kc.Option, err = utils.NewRemoteOption(ncmd.IP, defaultPort, defaultUser, ncmd.Password, ncmd.Key, ncmd.Iface)
		if err != nil {
			output.Fatalf("ssh %s:%s with error: %v\n", ncmd.IP, defaultPort, err)
		}
	}
	return kc
}

// ifconfig eth1 | grep inet | grep -v inet6 | awk -F " " '{print $2}'
func taskOutput(msg string) {
	output.Notef("============================== %s ==============================\n", msg)
}

func localFormat(node, cmd string) string {
	return fmt.Sprintf("vagrant ssh %s -c \"%s\"", node, cmd)
}
