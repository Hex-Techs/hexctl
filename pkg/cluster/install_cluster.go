package cluster

import (
	"encoding/json"
	"fmt"

	"github.com/Fize/n/pkg/utils"
	"github.com/ghodss/yaml"
)

const (
	kubeVersion     = "1.17.4"
	dockerVersion   = "18.09.9-3.el7"
	k8sKernelConfig = `net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
#sysctls for k8s node config
net.ipv4.tcp_slow_start_after_idle=0
net.core.rmem_max=16777216
fs.inotify.max_user_watches=524288
kernel.softlockup_all_cpu_backtrace=1
kernel.softlockup_panic=1
fs.file-max=2097152
fs.inotify.max_user_instances=8192
fs.inotify.max_queued_events=16384
vm.max_map_count=262144
net.core.netdev_max_backlog=16384
net.ipv4.tcp_wmem=4096 12582912 16777216
net.core.wmem_max=16777216
net.core.somaxconn=32768
net.ipv4.ip_forward=1
net.ipv4.tcp_max_syn_backlog=8096
net.bridge.bridge-nf-call-iptables=1
net.ipv4.tcp_rmem=4096 12582912 16777216

# arp cache
net.ipv4.neigh.default.gc_thresh1 = 80000
net.ipv4.neigh.default.gc_thresh2 = 90000
net.ipv4.neigh.default.gc_thresh3 = 100000`
	dockerConfig = `{
	"exec-opts": ["native.cgroupdriver=systemd"],
	"log-driver": "json-file",
	"log-opts": {
		"max-size": "100m",
		"max-file": "10"
	},
	"bip": "169.254.123.1/24",
	"oom-score-adjust": -1000,
	"ipv6": false,
	"storage-driver": "overlay2",
	"storage-opts":["overlay2.override_kernel_check=true"]
}`
)

type kubeadmConfig struct {
}

func installPackage(r *utils.RemoteOption) {
	r.Command.Cmd = "yum makecache fast && yum install -y nfs-utils ipvsadm pciutils bind-utils"
	if _, err := r.RunCommand(); err != nil {
		taskError(r, "install package with error", err)
	} else {
		taskSuccess(r, "install package success")
	}
	r.Command.Cmd = fmt.Sprintf("yum install -y docker-ce-%s kubeadm-%s kubectl-%s kubelet-%s", dockerVersion, kubeVersion, kubeVersion, kubeVersion)
	if _, err := r.RunCommand(); err != nil {
		taskFatal(r, "install kubernetes package with error", err)
	} else {
		taskSuccess(r, "install kubernetes package success")
	}
}

func setKernelConfig(r *utils.RemoteOption) {
	r.Command.Cmd = fmt.Sprintf(`echo '%s' > /etc/sysctl.d/99-k8s.conf && sysctl --system`, k8sKernelConfig)
	if _, err := r.RunCommand(); err != nil {
		taskError(r, "set kernel config with error", err)
	} else {
		taskSuccess(r, "set kernel config success")
	}
}

func setDockerConfig(r *utils.RemoteOption) {
	r.Command.Cmd = fmt.Sprintf(`mkdir /etc/docker && echo '%s' > /etc/docker/daemon.json && sysctl --system`, dockerConfig)
	if _, err := r.RunCommand(); err != nil {
		taskError(r, "set docker config with error", err)
	} else {
		taskSuccess(r, "set docker config success")
	}
}

func enableService(r *utils.RemoteOption) {
	r.Command.Cmd = `systemctl daemon-reload && systemctl enable --now docker.service && systemctl enable kubelet.service`
	if _, err := r.RunCommand(); err != nil {
		taskError(r, "enable docker and kubelet with error", err)
	} else {
		taskSuccess(r, "enable docker and kubelet success")
	}
}

func setKubeadmConfig(r *utils.RemoteOption, c *KubernetesCluster) {
	b, err := json.Marshal(c.ClusterConfig)
	if err != nil {
		taskFatal(r, "set kubeadm configuration with error", err)
	}
	yb, err := yaml.JSONToYAML(b)
	if err != nil {
		taskFatal(r, "set kubeadm configuration with error", err)
	}
	r.Command.Cmd = fmt.Sprintf("echo '%s' > /tmp/.kubeadm.conf", string(yb))
	if _, err := r.RunCommand(); err != nil {
		taskFatal(r, "set kubeadm configuration with error", err)
	} else {
		taskSuccess(r, "set kubeadm configuration success")
	}
}

func initCluster(r *utils.RemoteOption) {
	r.Command.Cmd = "kubeadm init --config=/tmp/.kubeadm.conf"
	if result, err := r.RunCommand(); err != nil {
		taskFatal(r, "init cluster with error", err)
	} else {
		taskSuccess(r, fmt.Sprintf("init cluster success: \n%s", result))
	}
}
