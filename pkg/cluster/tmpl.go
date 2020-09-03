package cluster

const (
	VagrantFileTmpl = `Vagrant.configure("2") do |config|
  {{- range $n := .Num }}
  config.vm.define :n{{- $n }} do |n{{- $n }}|
    n{{- $n }}.vm.provider "virtualbox" do |v|
      v.customize ["modifyvm", :id, "--name", "n{{- $n }}", "--cpus", {{ $.Cpu }}, "--memory", {{ $.Memory }}]
    end
    n{{- $n }}.vm.box = "centos/7"
    n{{- $n }}.vm.hostname = "n{{- $n }}"
    n{{- $n }}.vm.network "public_network", ip: "192.168.1.{{- $n }}"
  end
  {{- end }}
end`
)
