package crd

import (
	"runtime"

	"github.com/Hex-Techs/hexctl/internal/version"
	"github.com/Hex-Techs/hexctl/pkg/output"

	"github.com/spf13/cobra"
)

const (
	Logo = `  ___ ___                ___________           .__
 /   |   \   ____ ___  __\__    ___/___   ____ |  |__
/    ~    \_/ __ \\  \/  / |    |_/ __ \_/ ___\|  |  \
\    Y    /\  ___/ >    <  |    |\  ___/\  \___|   Y  \
 \___|_  /  \___  >__/\_ \ |____| \___  >\___  >___|  /
       \/       \/      \/            \/     \/     \/  
________________________________________________________`
)

func NewCrdVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Prints the version of hexctl crd",
		Run: func(cmd *cobra.Command, args []string) {
			output.Progressf("%s \n \n {hexctl version: %q, crd version: %q, kubernetes version: %q, commit: %q, go version: %q, GOOS: %q, GOARCH: %q}\n",
				Logo, version.GitVersion, version.CrdVersion, version.KubernetesVersion, version.GitCommit, runtime.Version(), runtime.GOOS, runtime.GOARCH)
		},
	}

	return cmd
}
