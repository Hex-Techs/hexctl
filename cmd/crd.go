package cmd

import (
	"github.com/Hex-Techs/hexctl/pkg/crd"
	"github.com/spf13/cobra"
)

var opts crd.GVK

var crdCmd = &cobra.Command{
	Use:              "crd",
	TraverseChildren: true,
	Short:            "crd is a command scaffold for kubernetes extensions",
	Long: `Provides libraries and tools to create new projects, APIs and controllers.
Includes tools for packaging artifacts into an installer container.

Typical project lifecycle:

- initialize a project:

  hexctl crd init api --domain example.com

- create one or more a new resource APIs and add your code to them:

  hexctl crd generate api --group <group> --version <version> --kind <Kind>

- create one or more a new controller for APIs and add your code to them:

  hexctl crd generate controller --kind <Kind>

After the scaffold is written, api will run make on the project.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init a Kubernetes extensions resource",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var apiInitCmd = &cobra.Command{
	Use:   "api",
	Short: "init kubernetes extensions project api",
	Example: `	# Init API
	- hexctl crd init api --domain example.com`,
	Run: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(crd.InitAPI(&opts))
	},
}

var ctrlInitCmd = &cobra.Command{
	Use:   "ctrl",
	Short: "init kubernetes extensions project controller",
	Example: `	# Init Controller
	- hexctl crd init ctrl --domain example.com`,
	Run: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(crd.InitController(&opts))
	},
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Scaffold a Kubernetes extensions resource",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Scaffold a Kubernetes API",
	Example: `	# Create a Foos API with Group: foo, Version: v1alpha1 and Kind: Foo
	- hexctl crd generate api --group foo --version v1alpha1 --kind Foo`,
	Run: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(crd.CreateAPI(&opts))
	},
}

var ctrlCmd = &cobra.Command{
	Use:   "ctrl",
	Short: "Scaffold a Kubernetes CRD controller",
	Example: `	# Create a Foos API with Group: foo, Version: v1alpha1 and Kind: Foo
	- hexctl crd generate controller --kind Foo`,
	Run: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(crd.CreateController(&opts))
	},
}

func init() {
	apiCmd.Flags().BoolVarP(&opts.UseNamespace, "namespaced", "n", true, "resource is namespaced")
	apiCmd.Flags().BoolVarP(&opts.UseStatus, "status", "s", true, "resource is has status")

	crdCmd.PersistentFlags().StringVarP(&opts.Domain, "domain", "d", "", "domain for groups (required)")
	crdCmd.PersistentFlags().StringVarP(&opts.Group, "group", "g", "", "resource group (required)")
	crdCmd.PersistentFlags().StringVarP(&opts.Version, "version", "v", "", "resouce version (required)")
	crdCmd.PersistentFlags().StringVarP(&opts.Kind, "kind", "k", "", "resource kind (required)")
	crdCmd.PersistentFlags().BoolVarP(&opts.Force, "force", "f", false, "attempt to create resource even if it already exists")

	crdCmd.MarkFlagRequired("domain")
	crdCmd.MarkFlagRequired("group")
	crdCmd.MarkFlagRequired("version")
	crdCmd.MarkFlagRequired("kind")

	initCmd.AddCommand(apiInitCmd)
	initCmd.AddCommand(ctrlInitCmd)

	generateCmd.AddCommand(apiCmd)
	generateCmd.AddCommand(ctrlCmd)

	crdCmd.AddCommand(initCmd)
	crdCmd.AddCommand(generateCmd)

	rootCmd.AddCommand(crdCmd)
}
