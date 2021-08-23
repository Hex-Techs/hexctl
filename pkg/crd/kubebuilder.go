package crd

import (
	"fmt"

	"github.com/Hex-Techs/hexctl/pkg/common/exec"
)

const (
	option = "--skip-go-version-check"
)

func Init(owner, repo string, gvk *GVK) {
	projectInit(repo)
	cmd := fmt.Sprintf("kubebuilder init --domain %s --owner %s %s", gvk.Domain, owner, option)
	exec.RunCommand(cmd)
}

func CreateAPI(gvk *GVK) {
	cmd := fmt.Sprintf("yes | kubebuilder create api --group %s --version %s --kind %s --force=%v --namespaced=%v",
		gvk.Group, gvk.Version, gvk.Kind, gvk.Force, gvk.UseNamespace)
	exec.RunCommand(cmd)
}

func CreateWebhook(gvk *GVK) {
	cmd := fmt.Sprintf("kubebuilder create webhook --group %s --version %s --kind %s --defaulting",
		gvk.Group, gvk.Version, gvk.Kind)
	exec.RunCommand(cmd)
}
