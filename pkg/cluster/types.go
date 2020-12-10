package cluster

import (
	"fmt"
	"os"

	"github.com/Hex-Techs/hexctl/pkg/output"
	"github.com/Hex-Techs/hexctl/pkg/utils"
	"golang.org/x/crypto/ssh/terminal"
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
	User             string
	Repo             string
	Password         string
	Key              string
	Iface            string
	Port             string
	KubeVersion      string
	DockerVersion    string
	Type             string
	CertSANs         []string
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
	return kc
}

func setSSHSession(kc *KubernetesCluster) {
	if kc.ncmd.User == "" {
		kc.ncmd.User = defaultUser
	}
	option, err := utils.NewRemoteOption(kc.ncmd.IP, defaultPort, kc.ncmd.User, kc.ncmd.Password, kc.ncmd.Key, kc.ncmd.Iface)
	if err != nil {
		output.Fatalf("ssh %s:%s with error: %v\n", kc.ncmd.IP, defaultPort, err)
	}
	kc.Option = option
	var command utils.Command
	option.Command = &command
	kc.Option = option
}

// ifconfig eth1 | grep inet | grep -v inet6 | awk -F " " '{print $2}'
func taskOutput(msg string) {
	fd := int(os.Stdin.Fd())
	length, _, err := terminal.GetSize(fd)
	if err != nil {
		output.Errorln(err)
		return
	}
	if length > 500 {
		length = 500
	}
	format := fmt.Sprintf("TASK: [%s] ", msg)
	mius := length - len(format)
	if mius <= 0 {
		output.Progressln(format)
	} else {
		for i := 0; i < mius; i++ {
			format = format + "*"
		}
		output.Progressln(format)
	}
}

func localFormat(node, cmd string) string {
	return fmt.Sprintf("vagrant ssh %s -c \"%s\"", node, cmd)
}
