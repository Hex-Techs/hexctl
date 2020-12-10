package cluster

import (
	"fmt"
	"strings"

	"github.com/Hex-Techs/hexctl/pkg/cluster/network"
	"github.com/Hex-Techs/hexctl/pkg/output"
	"github.com/Hex-Techs/hexctl/pkg/utils"
)

func (kc *KubernetesCluster) setRepo() {
	taskOutput("Configure repo mirror")
	var cmd1 string
	cmd2 := `sudo yum-config-manager --add-repo https://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo`
	if kc.ncmd.CN {
		cmd1 = fmt.Sprintf("sudo echo '%s' > /tmp/kubernetes.repo && sudo mv /tmp/kubernetes.repo /etc/yum.repos.d/kuberetes.repo", alikuberepo)
	} else {
		cmd1 = fmt.Sprintf("sudo echo '%s' > /tmp/kubernetes.repo && sudo mv /tmp/kubernetes.repo /etc/yum.repos.d/kuberetes.repo", googlekuberepo)
	}
	if kc.Type == "local" {
		utils.RunCommand(localFormat(kc.node, cmd1))
		utils.RunCommand(localFormat(kc.node, cmd2))
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
		output.Errorf("run command with error: %v, command: %s\n", err, cmd2)
	} else {
		output.Infoln(s)
	}
}

func (kc *KubernetesCluster) installClusterPackage() {
	taskOutput("Install docker and kubernetes packages")
	cmd := fmt.Sprintf("sudo yum install -y docker-ce-%s kubeadm-%s kubectl-%s kubelet-%s", kc.ncmd.DockerVersion, kc.ncmd.KubeVersion, kc.ncmd.KubeVersion, kc.ncmd.KubeVersion)
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

func (kc *KubernetesCluster) setDockerConfig() {
	taskOutput("Configure docker")
	var d string
	if kc.Type == "local" {
		d = strings.Replace(dockerConfig, "\"", "\\\"", -1)
	} else {
		d = dockerConfig
	}
	cmd := fmt.Sprintf(`sudo mkdir -p /etc/docker && sudo echo '%s' > /tmp/daemon.json && sudo mv /tmp/daemon.json /etc/docker/daemon.json`, d)
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

func (kc *KubernetesCluster) enableService() {
	taskOutput("Enable and start docker, kubelet")
	cmd := `sudo systemctl daemon-reload && sudo systemctl enable --now docker.service && sudo systemctl enable kubelet.service`
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

func (kc *KubernetesCluster) setNodeIP() {
	if kc.Type == "local" {
		taskOutput("set node ip")
		cmd := `sudo sed -i \"s#^KUBELET_EXTRA_ARGS=#KUBELET_EXTRA_ARGS=\"--node-ip=\$(ifconfig eth1 | grep netmask | awk -F' ' '{print \$2}')\"#\" /etc/sysconfig/kubelet`
		utils.RunCommand(localFormat(kc.node, cmd))
	}
}

func (kc *KubernetesCluster) initCluster() {
	taskOutput("Init cluster")
	var cmd string
	if kc.ncmd.Repo != "" {
		if kc.Type == "local" {
			cmd = fmt.Sprintf(`sudo kubeadm init --apiserver-advertise-address=\$(ifconfig eth1 | grep netmask | awk -F' ' '{print \$2}') --kubernetes-version=%s --pod-network-cidr=%s --image-repository=%s --service-cidr=%s`,
				defaultKubernetesVersion, kc.ncmd.PodCIDR, kc.ncmd.Repo, kc.ncmd.ServiceCIDR)
		} else {
			cmd = fmt.Sprintf(`sudo kubeadm init --kubernetes-version=%s --pod-network-cidr=%s --image-repository=%s --service-cidr=%s`,
				defaultKubernetesVersion, kc.ncmd.PodCIDR, kc.ncmd.Repo, kc.ncmd.ServiceCIDR)
		}
	} else {
		if kc.Type == "local" {
			cmd = fmt.Sprintf(`sudo kubeadm init --apiserver-advertise-address=\$(ifconfig eth1 | grep netmask | awk -F' ' '{print \$2}') --kubernetes-version=%s --pod-network-cidr=%s --service-cidr=%s`,
				defaultKubernetesVersion, kc.ncmd.PodCIDR, kc.ncmd.ServiceCIDR)
		} else {
			cmd = fmt.Sprintf(`sudo kubeadm init --kubernetes-version=%s --pod-network-cidr=%s --service-cidr=%s`,
				defaultKubernetesVersion, kc.ncmd.PodCIDR, kc.ncmd.ServiceCIDR)
		}
	}
	if len(kc.ncmd.CertSANs) != 0 {
		if len(kc.ncmd.CertSANs) == 1 {
			cmd = fmt.Sprintf("%s --apiserver-cert-extra-sans %s", cmd, kc.ncmd.CertSANs[0])
		} else {
			cmd = fmt.Sprintf("%s --apiserver-cert-extra-sans %s", cmd, strings.Join(kc.ncmd.CertSANs, ","))
		}
	}
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

func (kc *KubernetesCluster) joinCluster() {
	taskOutput("Join cluster")
	cmd := fmt.Sprintf("sudo kubeadm join %s --token %s --discovery-token-unsafe-skip-ca-verification", kc.ncmd.Endpoint, kc.ncmd.Token)
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

func (kc *KubernetesCluster) setKubeConfig() {
	taskOutput("Set kubeconfig")
	var cmd string
	if kc.Type == "local" {
		cmd = `mkdir -p \$HOME/.kube && sudo cp -i /etc/kubernetes/admin.conf \$HOME/.kube/config && sudo chown \$(id -u):\$(id -g) \$HOME/.kube/config`
	} else {
		cmd = `mkdir -p $HOME/.kube && sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config && sudo chown $(id -u):$(id -g) $HOME/.kube/config`
	}
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

func (kc *KubernetesCluster) setFlannel() {
	taskOutput("Set flannel network cni")
	cm := network.GenerateFlannelCM(kc.ncmd.PodCIDR)
	cmd1 := fmt.Sprintf(`cat <<EOF | kubectl apply -f -
%s
EOF`, network.Psp)
	var cmd2 string
	if kc.Type == "local" {
		cr := network.GenerateFlannelCR("\\\"\\\"")
		cmd2 = fmt.Sprintf(`cat <<EOF | kubectl apply -f -
%s
EOF`, cr)
	} else {
		cr := network.GenerateFlannelCR("''")
		cmd2 = fmt.Sprintf(`cat <<EOF | kubectl apply -f -
%s
EOF`, cr)
	}
	cmd3 := fmt.Sprintf(`cat <<EOF | kubectl apply -f -
%s
EOF`, network.Crb)
	cmd4 := fmt.Sprintf(`cat <<EOF | kubectl apply -f -
%s
EOF`, network.Sa)
	cmd5 := fmt.Sprintf(`cat <<EOF | kubectl apply -f -
%s
EOF`, cm)
	cmd6 := fmt.Sprintf(`cat <<EOF | kubectl apply -f -
%s
EOF`, network.Ds)
	if kc.Type == "local" {
		utils.RunCommand(localFormat(kc.node, cmd1))
		utils.RunCommand(localFormat(kc.node, cmd2))
		utils.RunCommand(localFormat(kc.node, cmd3))
		utils.RunCommand(localFormat(kc.node, cmd4))
		utils.RunCommand(localFormat(kc.node, cmd5))
		utils.RunCommand(localFormat(kc.node, cmd6))
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
		output.Errorf("run command with error: %v, command: %s\n", err, cmd2)
	} else {
		output.Infoln(s)
	}
	setSSHSession(kc)
	kc.Option.Command.Cmd = cmd3
	s, err = kc.Option.RunCommand()
	if err != nil {
		output.Errorf("run command with error: %v, command: %s\n", err, cmd3)
	} else {
		output.Infoln(s)
	}
	setSSHSession(kc)
	kc.Option.Command.Cmd = cmd4
	s, err = kc.Option.RunCommand()
	if err != nil {
		output.Errorf("run command with error: %v, command: %s\n", err, cmd4)
	} else {
		output.Infoln(s)
	}
	setSSHSession(kc)
	kc.Option.Command.Cmd = cmd5
	s, err = kc.Option.RunCommand()
	if err != nil {
		output.Errorf("run command with error: %v, command: %s\n", err, cmd5)
	} else {
		output.Infoln(s)
	}
	setSSHSession(kc)
	kc.Option.Command.Cmd = cmd6
	s, err = kc.Option.RunCommand()
	if err != nil {
		output.Errorf("run command with error: %v, command: %s\n", err, cmd6)
	} else {
		output.Infoln(s)
	}
}
