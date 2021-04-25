package crd

import "github.com/spf13/cobra"

func NewPushCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "push",
		Short: "push image to docker harbor",
		Long:  ``,
		Example: `	#Push image, it will push image that you made.
		- hexctl crd push`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	return cmd
}
