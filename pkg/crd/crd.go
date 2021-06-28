// Package crd help to create a operator controller by kubeilder or code-generate
package crd

import (
	"github.com/Hex-Techs/hexctl/pkg/common/exec"
)

// GVK this is the infomation of your CRD
type GVK struct {
	Domain       string
	Group        string
	Version      string
	Kind         string
	UseNamespace bool
	UseStatus    bool
	// attempt to create resource even if it already exists
	Force bool
}

type WorkOption struct {
	Options   string
	Generated string
	API       string
	Group     string
	Version   string
}

func projectInit(repo string) {
	if repo != "" {
		cmd := "go mod init " + repo
		exec.RunCommand(cmd)
	}
}
