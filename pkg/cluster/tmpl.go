package cluster

const (
	defaultKubeVersion       = "1.19.4-0"
	defaultDockerVersion     = "19.03.14-3.el7"
	defaultKubernetesVersion = "v1.19.4"
	VagrantFileTmpl          = `$script = <<-SCRIPT
sudo setenforce 0
sudo sed -i 's/SELINUX=enforcing/SELINUX=disabled/g' /etc/selinux/config
SCRIPT

Vagrant.configure("2") do |config|
  {{- range $n := .Num }}
  config.vm.define :n{{- $n }} do |n{{- $n }}|
    n{{- $n }}.vm.provider "virtualbox" do |v|
      v.customize ["modifyvm", :id, "--name", "n{{- $n }}", "--cpus", {{ $.Cpu }}, "--memory", {{ $.Memory }}]
    end
    n{{- $n }}.vm.box = "centos/7"
    n{{- $n }}.vm.hostname = "n{{- $n }}"
	n{{- $n }}.vm.network :private_network, ip: "172.20.1.10{{- $n }}"
    n{{- $n }}.vm.provision "shell", inline: $script
  end
  {{- end }}
end`
	googlekuberepo = `[kubernetes]
name=Kubernetes
baseurl=https://packages.cloud.google.com/yum/repos/kubernetes-el7-x86_64
enabled=1
gpgcheck=0
repo_gpgcheck=1
gpgkey=https://packages.cloud.google.com/yum/doc/yum-key.gpg https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg`
	alikuberepo = `[kubernetes]
name=Kubernetes
baseurl=https://mirrors.aliyun.com/kubernetes/yum/repos/kubernetes-el7-x86_64/
enabled=1
gpgcheck=0
repo_gpgcheck=1
gpgkey=https://mirrors.aliyun.com/kubernetes/yum/doc/yum-key.gpg https://mirrors.aliyun.com/kubernetes/yum/doc/rpm-package-key.gpg`
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
