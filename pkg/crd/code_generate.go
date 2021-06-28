package crd

import (
	"github.com/Hex-Techs/hexctl/pkg/common/exec"
	"github.com/Hex-Techs/hexctl/pkg/common/file"
	"github.com/Hex-Techs/hexctl/pkg/crd/templates"
)

const (
	generators = "client,lister,informer"
	tools      = "hack/tools.go"
	script     = "hack/update-codegen.sh"
)

func Generate(gvk *GVK, wd *WorkOption) {
	toolsInit()
	workOption := &WorkOption{
		Options:   generators,
		Generated: wd.Generated,
		API:       wd.API,
		Group:     gvk.Group,
		Version:   gvk.Version,
	}
	scriptInit(workOption)
	// TODO: run script
}

func toolsInit() {
	file.WriteByTemp("tools", templates.ToolsTemp, tools, nil)
}

func scriptInit(opts *WorkOption) {
	file.WriteByTemp("script", templates.UpdateCodeGenTemp, script, opts)
	cmd := "chmod +x " + script
	exec.RunCommand(cmd)
}
