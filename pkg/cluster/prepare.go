package cluster

import (
	"fmt"

	"github.com/Fize/n/pkg/output"
	"github.com/Fize/n/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta1"
)

const (
	googlekuberepo = `[kubernetes]
name=Kubernetes
baseurl=https://packages.cloud.google.com/yum/repos/kubernetes-el7-x86_64
enabled=1
gpgcheck=1
repo_gpgcheck=1 gpgkey=https://packages.cloud.google.com/yum/doc/yum-key.gpg https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg`
	huaweikuberepo = `[kubernetes]
name=Kubernetes
baseurl=https://mirrors.huaweicloud.com/kubernetes/yum/repos/kubernetes-el7-x86_64
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://mirrors.huaweicloud.com/kubernetes/yum/doc/yum-key.gpg https://mirrors.huaweicloud.com/kubernetes/yum/doc/rpm-package-key.gpg`
)

// Run run Command
func Run(ncmd *ClusterCommand) {
	handleWorkDir()
	renderVagrantFile(ncmd)
	startVirtualMachine()
	// ncmd.Port = DefaultPort
	// if ncmd.Master != "" {
	// k := &KubernetesCluster{}
	// k.initMaster(ncmd.APIServer, DefaultPort, ncmd.Key, ncmd.Password, ncmd.Iface)
	// k.initConfig(ncmd)

	// taskOutput("init a new kubernetes cluster")

	// setHostName(k.MasterOption)
	// setKubeRepo(k.MasterOption, ncmd.Repo)
	// swapoff(k.MasterOption)
	// createDir(k.MasterOption)
	// formatDisk(k.MasterOption, ncmd.Volume)
	// mount(k.MasterOption, ncmd.Volume)

	// taskOutput("init a new kubernetes cluster success!")
	// } else {
	// k := &KubernetesCluster{}
	// k.initConfig(ncmd)
	// taskOutput("join a kubernetes cluster")

	// for _, r := range ncmd.Node {
	// k.initNode(r, DefaultPort, ncmd.Key, ncmd.Password, ncmd.Iface)
	// }

	// var stop chan bool

	// for _, r := range k.NodeOption {
	// go func(r *utils.RemoteOption, stop chan bool) {
	// setHostName(r)
	// setKubeRepo(r, ncmd.Repo)
	// swapoff(r)
	// createDir(r)
	// formatDisk(r, ncmd.Volume)
	// mount(r, ncmd.Volume)
	// stop <- true
	// }(r, stop)
	// }

	// for {
	// if len(stop) == len(k.NodeOption) {
	// break
	// }
	// }
	// taskOutput("join a kubernetes cluster success!")
	// }
}

func (kc *KubernetesCluster) initMaster(host, port, key, password, netDevice string) {
	remote, err := utils.NewRemoteOption(host, port, DefaultUser, password, key, netDevice)
	if err != nil {
		output.Fatalln(err)
	}
	kc.MasterOption = remote
}

func (kc *KubernetesCluster) initNode(host, port, key, password, netDevice string) {
	// taskOutput("init a new kubernetes cluster")
	remote, err := utils.NewRemoteOption(host, port, DefaultPort, password, key, netDevice)
	if err != nil {
		output.Fatalln(err)
	}
	kc.NodeOption = append(kc.NodeOption, remote)
}

func (kc *KubernetesCluster) initConfig(ncmd *ClusterCommand) {
	if ncmd.ServicePortRange == "" {
		ncmd.ServicePortRange = DefaultServicePortRange
	}
	kc.ClusterConfig = &v1beta1.ClusterConfiguration{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "kubeadm.k8s.io/v1beta1",
			Kind:       "ClusterConfiguration",
		},
		KubernetesVersion:    "v" + kubeVersion,
		ControlPlaneEndpoint: ncmd.APIServer,
		APIServer: v1beta1.APIServer{
			ControlPlaneComponent: v1beta1.ControlPlaneComponent{
				ExtraArgs: map[string]string{
					"max-mutating-requests-inflight": "600",
					"max-requests-inflight":          "1000",
					"service-node-port-range":        ncmd.ServicePortRange,
				},
			},
			CertSANs: append(ncmd.CertHost, ncmd.Master),
		},
		Etcd: v1beta1.Etcd{
			Local: &v1beta1.LocalEtcd{
				ServerCertSANs: append(ncmd.CertHost, ncmd.Master),
				PeerCertSANs:   append(ncmd.CertHost, ncmd.Master),
			},
		},
		ControllerManager: v1beta1.ControlPlaneComponent{
			ExtraArgs: map[string]string{
				"node-cidr-mask-size": "24",
			},
		},
		Networking: v1beta1.Networking{
			PodSubnet:     ncmd.PodCIDR,
			ServiceSubnet: ncmd.ServiceCIDR,
		},
		ImageRepository: ncmd.Repo,
	}
	// kc.KubeProxyConfig = &v1alpha1.KubeProxyConfiguration{
	// Mode: "ipvs",
	// }
}

func setHostName(r *utils.RemoteOption) {
	r.Command.Cmd = fmt.Sprintf(`hostnamectl set-hostname $(ifconfig %s | grep inet | awk -F " " '{print $2}')`, r.NetDevice)
	if _, err := r.RunCommand(); err != nil {
		taskFatal(r, "set hostname with error", err)
	} else {
		taskSuccess(r, "set hostname success")
	}
}

func setKubeRepo(r *utils.RemoteOption, repo string) {
	if repo != "" {
		r.Command.Cmd = fmt.Sprintf("echo '%s' > /etc/yum.repos.d/kuberetes.repo", huaweikuberepo)
	} else {
		r.Command.Cmd = fmt.Sprintf("echo '%s' > /etc/yum.repos.d/kuberetes.repo", googlekuberepo)
	}
	if _, err := r.RunCommand(); err != nil {
		taskFatal(r, "set kubernetes repo with error", err)
	} else {
		taskSuccess(r, "set kubernetes repo success")
	}
	r.Command.Cmd = `curl -s -o /etc/yum.repos.d/docker-ce.repo https://mirrors.huaweicloud.com/docker-ce/linux/centos/docker-ce.repo \
&& sed -i 's+download.docker.com+mirrors.huaweicloud.com/docker-ce+' /etc/yum.repos.d/docker-ce.repo`
	if _, err := r.RunCommand(); err != nil {
		taskFatal(r, "set docker-ce repo with error", err)
	} else {
		taskSuccess(r, "set docker-ce repo success")
	}
}

func swapoff(r *utils.RemoteOption) {
	r.Command.Cmd = "swapoff -a && sysctl -w vm.swappiness=0"
	if _, err := r.RunCommand(); err != nil {
		taskError(r, "disable swap with error", err)
	} else {
		taskSuccess(r, "disable swap success")
	}
	r.Command.Cmd = "cp /etc/fstab /etc/fstab.bak"
	if _, err := r.RunCommand(); err != nil {
		taskError(r, "backup /etc/fstab with error", err)
	} else {
		taskSuccess(r, "backup /etc/fstab success")
	}
	r.Command.Cmd = "sed -i '/swap/d' /etc/fstab"
	if _, err := r.RunCommand(); err != nil {
		taskError(r, "modify /etc/fstab with error", err)
	} else {
		taskSuccess(r, "modify /etc/fstab success")
	}
}

func secretoff(r *utils.RemoteOption) {
	r.Command.Cmd = `setenforce 0 \
&& sed -i 's/^SELINUX=enforcing$/SELINUX=disabled/' /etc/selinux/config`
	if _, err := r.RunCommand(); err != nil {
		taskError(r, "disable selinux with error", err)
	} else {
		taskSuccess(r, "disable selinux success")
	}
	r.Command.Cmd = `iptables -t nat -F && iptables -F && systemctl disable firewalld && systemctl stop firewalld`
	if _, err := r.RunCommand(); err != nil {
		taskError(r, "disable iptables or firewalld with error", err)
	} else {
		taskSuccess(r, "disable iptables or firewalld success")
	}
}

func createDir(r *utils.RemoteOption) {
	r.Command.Cmd = "mkdir /var/lib/{container,docker,kubelet} && mkdir /apkget"
	if _, err := r.RunCommand(); err != nil {
		taskError(r, "create directory /var/lib/container /var/lib/docker /var/lib/kubelet /apkget with error", err)
	} else {
		taskSuccess(r, "create directory success /var/lib/container /var/lib/docker /var/lib/kubelet /apkget")
	}
}

func formatDisk(r *utils.RemoteOption, disk string) {
	r.Command.Cmd = fmt.Sprintf(`fdisk -S 56 %s <<-EOF
n
p
1


wq
EOF`, disk)
	if _, err := r.RunCommand(); err != nil {
		taskError(r, "partition disk with error", err)
	} else {
		taskSuccess(r, "partition disk success")
	}
	r.Command.Cmd = fmt.Sprintf("mkfs.xfs -m -f reflink=1 %s1", disk)
	if _, err := r.RunCommand(); err != nil {
		taskError(r, "format disk by xfs with error", err)
	} else {
		taskSuccess(r, "format disk by xfs success")
	}
}

func mount(r *utils.RemoteOption, disk string) {
	r.Command.Cmd = fmt.Sprintf("mount -t xfs %s1 /var/lib/container", disk)
	if _, err := r.RunCommand(); err != nil {
		taskError(r, "mount disk with error", err)
	} else {
		taskSuccess(r, "mount disk success")
	}
	r.Command.Cmd = "mkdir /var/lib/container/{docker,kubelet,apkget}"
	if _, err := r.RunCommand(); err != nil {
		taskError(r, "create directory /var/lib/container/docker /var/lib/container/kubelet /var/lib/container/apkget with error", err)
	} else {
		taskSuccess(r, "create directory /var/lib/container/docker /var/lib/container/kubelet /var/lib/container/apkget success")
	}
	r.Command.Cmd = fmt.Sprintf(`echo '%s/ /var/lib/container  xfs defaults 0 0' >> /etc/fstab \
&& echo '/var/lib/container/kubelet /var/lib/kubelet none defaults,bind 0 0' >> /etc/fstab \
&& echo '/var/lib/container/docker /var/lib/docker none defaults,bind 0 0' >> /etc/fstab \
&& echo '/var/lib/container/apkget /apkget none defaults,bind 0 0' >> /etc/fstab \
&& mount -a`, disk)
	if _, err := r.RunCommand(); err != nil {
		taskError(r, "mount directory with error", err)
	} else {
		taskSuccess(r, "mount directory success")
	}
}

func taskStart(r *utils.RemoteOption, msg string) {
	output.Notef("[%s:%s] %s: ", r.Host, r.Port, msg)
}

func taskSuccess(r *utils.RemoteOption, msg string) {
	output.Progressf("[%s:%s] ", r.Host, r.Port)
	output.Successln(msg)
}

func taskFatal(r *utils.RemoteOption, msg string, err error) {
	output.Progressf("[%s:%s] %s: ", r.Host, r.Port, msg)
	output.Errorln(err)
}

func taskError(r *utils.RemoteOption, msg string, err error) {
	output.Progressf("[%s:%s] %s: ", r.Host, r.Port, msg)
	output.Errorln(err)
}

func taskOutput(msg string) {
	output.Notef("-------------------- %s --------------------\n", msg)
}
