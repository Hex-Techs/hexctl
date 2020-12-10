package cluster

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Hex-Techs/n/pkg/output"
	"github.com/Hex-Techs/n/pkg/utils"
)

func (kc *KubernetesCluster) argsCheck() {
	if err := validateCIDR(kc.ncmd.PodCIDR); err != nil {
		output.Fatalln("pod network cidr:", err)
	}
	if err := validateCIDR(kc.ncmd.ServiceCIDR); err != nil {
		output.Fatalln("service cidr:", err)
	}
	if kc.ncmd.ServicePortRange == "" {
		kc.ncmd.ServicePortRange = defaultServicePortRange
	}
	if kc.ncmd.DockerVersion == "" {
		kc.ncmd.DockerVersion = defaultDockerVersion
	}
	if kc.ncmd.KubeVersion == "" {
		kc.ncmd.KubeVersion = defaultKubeVersion
	}
}

func (kc *KubernetesCluster) swapoff() {
	if kc.Type == "local" {
		if err := os.Chdir(homeDir()); err != nil {
			output.Fatalln(err)
		}
	}
	taskOutput("Disable swap")
	cmd1 := "sudo swapoff -a && sudo sysctl -w vm.swappiness=0"
	cmd2 := "sudo cp /etc/fstab /etc/fstab.bak"
	cmd3 := "sudo sed -i '/swap/d' /etc/fstab"
	if kc.Type == "local" {
		utils.RunCommand(localFormat(kc.node, cmd1))
		utils.RunCommand(localFormat(kc.node, cmd2))
		utils.RunCommand(localFormat(kc.node, cmd3))
		return
	}
	setSSHSession(kc)
	kc.Option.Command.Cmd = cmd1
	s, err := kc.Option.RunCommand()
	if err != nil {
		output.Errorf("run command with error: %v, command: %s\n", err, cmd1)
	} else {
		output.Infoln(s)
	}
	setSSHSession(kc)
	kc.Option.Command.Cmd = cmd2
	s, err = kc.Option.RunCommand()
	if err != nil {
		output.Errorf("run command with error: %v, command: %s\n", err, cmd1)
	} else {
		output.Infoln(s)
	}
	setSSHSession(kc)
	kc.Option.Command.Cmd = cmd3
	s, err = kc.Option.RunCommand()
	if err != nil {
		output.Errorf("run command with error: %v, command: %s\n", err, cmd1)
	} else {
		output.Infoln(s)
	}
}

func (kc *KubernetesCluster) installPackage() {
	taskOutput("Install packages")
	cmd := "sudo yum makecache fast && sudo yum install -y ipvsadm bind-utils net-tools yum-utils"
	if kc.Type == "local" {
		utils.RunCommand(localFormat(kc.node, cmd))
		return
	}
	setSSHSession(kc)
	kc.Option.Command.Cmd = cmd
	s, err := kc.Option.RunCommand()
	if err != nil {
		output.Errorf("run command with error: %v, command: %s\n", err, cmd)
	} else {
		output.Infoln(s)
	}
}

func (kc *KubernetesCluster) setKernelConfig() {
	taskOutput("Configure kernel")
	cmd := fmt.Sprintf(`sudo echo '%s' > /tmp/89-k8s.conf && sudo mv /tmp/89-k8s.conf /etc/sysctl.d && sudo sysctl --system`, k8sKernelConfig)
	if kc.Type == "local" {
		utils.RunCommand(localFormat(kc.node, cmd))
		return
	}
	setSSHSession(kc)
	kc.Option.Command.Cmd = cmd
	s, err := kc.Option.RunCommand()
	if err != nil {
		output.Errorf("run command with error: %v, command: %s\n", err, cmd)
	} else {
		output.Infoln(s)
	}
}

func validateCIDR(cidr string) error {
	if cidr == "" {
		return nil
	}
	cidrSlice := strings.Split(cidr, "/")
	if len(cidrSlice) != 2 {
		return fmt.Errorf("invalid cidr")
	}
	segment := cidrSlice[0]
	netmask, err := strconv.Atoi(cidrSlice[1])
	if err != nil {
		return err
	}
	if netmask >= 32 || netmask < 8 {
		return fmt.Errorf("invalid netmask")
	}
	if err := utils.ValidataIP(segment); err != nil {
		return err
	}
	return nil
}
