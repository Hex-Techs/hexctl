package crd

import "github.com/spf13/cobra"

func NewCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Scaffold a Kubernetes API",
		Example: `	# Create a Foos API with Group: foo, Version: v1alpha1 and Kind: Foo 
		- hexctl crd create api --group foo --version v1alpha1 --kind Foo`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(NewApiCommand())

	return cmd
}

func NewApiCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api",
		Short: "Scaffold a Kubernetes API",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	return cmd
}
