package cluster

const (
	VagrantFileTmpl = `$script = <<-SCRIPT
sudo setenforce 0
sudo sed -i 's/SELINUX=enforcing/SELINUX=disabled/g' /etc/selinux/config
echo "{{ .Pub }}" >> .ssh/anthorized_keys
chmod 0600 .ssh/authorized_keys
SCRIPT

Vagrant.configure("2") do |config|
  {{- range $n := .Num }}
  config.vm.define :n{{- $n }} do |n{{- $n }}|
    n{{- $n }}.vm.provider "virtualbox" do |v|
      v.customize ["modifyvm", :id, "--name", "n{{- $n }}", "--cpus", {{ $.Cpu }}, "--memory", {{ $.Memory }}]
    end
    n{{- $n }}.vm.box = "centos/7"
    n{{- $n }}.vm.hostname = "n{{- $n }}"
    n{{- $n }}.vm.network "public_network", ip: "172.168.1.{{- $n }}"
    n{{- $n }}.vm.provision "shell", inline: $script
  end
  {{- end }}
end`
)
