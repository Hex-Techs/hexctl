package crd

import "github.com/spf13/cobra"

func NewInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "init kubernetes extensions project",
		Example: `	# Init project with license
		- hexctl crd init --domain example.com --license apache2 --owner "The Hex-Techs authors"
	
		# Init project without license
		- hexctl crd init --domain example.com`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	return cmd
}
