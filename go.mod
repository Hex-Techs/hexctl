module github.com/Hex-Techs/hexctl

go 1.19

require (
	github.com/Masterminds/sprig v2.22.0+incompatible
	github.com/gookit/color v1.4.2
	github.com/jedib0t/go-pretty/v6 v6.2.4
	github.com/manifoldco/promptui v0.7.0
	github.com/pkg/sftp v1.11.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.7.0
	golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83
	golang.org/x/text v0.3.6
	k8s.io/apimachinery v0.22.4
	k8s.io/client-go v0.22.4
)

require (
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/go-logr/logr v0.4.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.5 // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/googleapis/gnostic v0.5.5 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/imdario/mergo v0.3.5 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/juju/ansiterm v0.0.0-20180109212912-720a0952cc2a // indirect
	github.com/kr/fs v0.1.0 // indirect
	github.com/lunixbochs/vtclean v1.0.0 // indirect
	github.com/magiconair/properties v1.8.1 // indirect
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/mapstructure v1.1.2 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/pelletier/go-toml v1.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cast v1.3.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	github.com/xo/terminfo v0.0.0-20210125001918-ca9a967f8778 // indirect
	golang.org/x/net v0.0.0-20210520170846-37e1c6afe023 // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	golang.org/x/sys v0.0.0-20210616094352-59db8d763f22 // indirect
	golang.org/x/term v0.0.0-20210220032956-6a3ed077a48d // indirect
	golang.org/x/time v0.0.0-20210723032227-1f47c861a9ac // indirect
	google.golang.org/appengine v1.6.5 // indirect
	google.golang.org/protobuf v1.26.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.51.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	k8s.io/api v0.22.4 // indirect
	k8s.io/klog/v2 v2.9.0 // indirect
	k8s.io/utils v0.0.0-20210819203725-bdf08cb9a70a // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.1.2 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)

replace k8s.io/api v0.0.0 => k8s.io/api v0.22.4

replace k8s.io/apiextensions-apiserver v0.0.0 => k8s.io/apiextensions-apiserver v0.22.4

replace k8s.io/apimachinery v0.0.0 => k8s.io/apimachinery v0.22.4

replace k8s.io/apiserver v0.0.0 => k8s.io/apiserver v0.22.4

replace k8s.io/cli-runtime v0.0.0 => k8s.io/cli-runtime v0.22.4

replace k8s.io/client-go v0.0.0 => k8s.io/client-go v0.22.4

replace k8s.io/cloud-provider v0.0.0 => k8s.io/cloud-provider v0.22.4

replace k8s.io/cluster-bootstrap v0.0.0 => k8s.io/cluster-bootstrap v0.22.4

replace k8s.io/code-generator v0.0.0 => k8s.io/code-generator v0.22.4

replace k8s.io/component-base v0.0.0 => k8s.io/component-base v0.22.4

replace k8s.io/cri-api v0.0.0 => k8s.io/cri-api v0.22.4

replace k8s.io/csi-translation-lib v0.0.0 => k8s.io/csi-translation-lib v0.22.4

replace k8s.io/kube-aggregator v0.0.0 => k8s.io/kube-aggregator v0.22.4

replace k8s.io/kube-controller-manager v0.0.0 => k8s.io/kube-controller-manager v0.22.4

replace k8s.io/kube-proxy v0.0.0 => k8s.io/kube-proxy v0.22.4

replace k8s.io/kube-scheduler v0.0.0 => k8s.io/kube-scheduler v0.22.4

replace k8s.io/kubectl v0.0.0 => k8s.io/kubectl v0.22.4

replace k8s.io/kubelet v0.0.0 => k8s.io/kubelet v0.22.4

replace k8s.io/legacy-cloud-providers v0.0.0 => k8s.io/legacy-cloud-providers v0.22.4

replace k8s.io/metrics v0.0.0 => k8s.io/metrics v0.22.4

replace k8s.io/sample-apiserver v0.0.0 => k8s.io/sample-apiserver v0.22.4

replace github.com/spf13/cobra v1.0.0 => github.com/spf13/cobra v1.1.3
