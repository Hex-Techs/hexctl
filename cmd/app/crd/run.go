package crd

import "github.com/spf13/cobra"

func NewRunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run an custom project in a variety of environments",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	return cmd
}
