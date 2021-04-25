package crd

import "github.com/spf13/cobra"

func NewCrdCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "crd",
		Short: "crd is a command for building kubernetes extensions and tools",
		Long: `Development kit for building Kubernetes extensions and tools.

	Provides libraries and tools to create new projects, APIs and controllers.
	Includes tools for packaging artifacts into an installer container.

	Typical project lifecycle:

	- initialize a project:

	  hexctl crd init --domain example.com --license apache2 --owner "The Hex-Techs authors"

	- create one or more a new resource APIs and add your code to them:

	  hexctl crd create api --group <group> --version <version> --kind <Kind>

	After the scaffold is written, api will run make on the project.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(NewInitCommand())
	cmd.AddCommand(NewCreateCommand())
	cmd.AddCommand(NewBuildCommand())
	cmd.AddCommand(NewPushCommand())
	cmd.AddCommand(NewGenerateCommand())
	cmd.AddCommand(NewRunCommand())
	cmd.AddCommand(NewCrdVersionCommand())

	return cmd
}
