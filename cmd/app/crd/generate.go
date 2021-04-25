package crd

import "github.com/spf13/cobra"

func NewGenerateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Invokes a specific generator",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(NewK8SCommand())
	cmd.AddCommand(NewOpenApiCommand())

	return cmd
}

func NewK8SCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "k8s",
		Short: "generate k8s clients whit the API",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	return cmd
}

func NewOpenApiCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "openapi",
		Short: "generate k8s openAPI whit the API",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	return cmd
}
