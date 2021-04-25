package crd

import "github.com/spf13/cobra"

func NewBuildCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build",
		Short: "build image for the project",
		Long:  ``,
		Example: `	#Build image
		- hexctl crd build examlpe.com/hex-techs/image:v0.0.0`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	return cmd
}
