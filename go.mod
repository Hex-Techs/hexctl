module github.com/Hex-Techs/hexctl

go 1.14

require (
	github.com/fsnotify/fsnotify v1.4.7
	github.com/ghodss/yaml v1.0.0
	github.com/gizak/termui/v3 v3.1.0
	github.com/gookit/color v1.4.2
	github.com/jedib0t/go-pretty/v6 v6.2.4
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/kopoli/go-terminal-size v0.0.0-20170219200355-5c97524c8b54
	github.com/lunixbochs/vtclean v1.0.0 // indirect
	github.com/manifoldco/promptui v0.7.0
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pkg/sftp v1.11.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/viper v1.7.0
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/apimachinery v0.18.6
	k8s.io/client-go v0.18.6
)

replace k8s.io/api v0.0.0 => k8s.io/api v0.18.6

replace k8s.io/apiextensions-apiserver v0.0.0 => k8s.io/apiextensions-apiserver v0.18.6

replace k8s.io/apimachinery v0.0.0 => k8s.io/apimachinery v0.18.6

replace k8s.io/apiserver v0.0.0 => k8s.io/apiserver v0.18.6

replace k8s.io/cli-runtime v0.0.0 => k8s.io/cli-runtime v0.18.6

replace k8s.io/client-go v0.0.0 => k8s.io/client-go v0.18.6

replace k8s.io/cloud-provider v0.0.0 => k8s.io/cloud-provider v0.18.6

replace k8s.io/cluster-bootstrap v0.0.0 => k8s.io/cluster-bootstrap v0.18.6

replace k8s.io/code-generator v0.0.0 => k8s.io/code-generator v0.18.6

replace k8s.io/component-base v0.0.0 => k8s.io/component-base v0.18.6

replace k8s.io/cri-api v0.0.0 => k8s.io/cri-api v0.18.6

replace k8s.io/csi-translation-lib v0.0.0 => k8s.io/csi-translation-lib v0.18.6

replace k8s.io/kube-aggregator v0.0.0 => k8s.io/kube-aggregator v0.18.6

replace k8s.io/kube-controller-manager v0.0.0 => k8s.io/kube-controller-manager v0.18.6

replace k8s.io/kube-proxy v0.0.0 => k8s.io/kube-proxy v0.18.6

replace k8s.io/kube-scheduler v0.0.0 => k8s.io/kube-scheduler v0.18.6

replace k8s.io/kubectl v0.0.0 => k8s.io/kubectl v0.18.6

replace k8s.io/kubelet v0.0.0 => k8s.io/kubelet v0.18.6

replace k8s.io/legacy-cloud-providers v0.0.0 => k8s.io/legacy-cloud-providers v0.18.6

replace k8s.io/metrics v0.0.0 => k8s.io/metrics v0.18.6

replace k8s.io/sample-apiserver v0.0.0 => k8s.io/sample-apiserver v0.18.6

replace github.com/spf13/cobra v1.0.0 => github.com/spf13/cobra v1.1.3
