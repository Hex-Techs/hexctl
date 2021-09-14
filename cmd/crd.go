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

  hexctl crd init --domain example.com

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
	Short: "init kubernetes extensions project",
	Example: `	# Init project
	- hexctl crd init --domain example.com`,
	Run: func(cmd *cobra.Command, args []string) {
		crd.Init(&opts)
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
		crd.CreateAPI(&opts)
	},
}

var ctrlCmd = &cobra.Command{
	Use:   "controller",
	Short: "Scaffold a Kubernetes CRD controller",
	Example: `	# Create a Foos API with Group: foo, Version: v1alpha1 and Kind: Foo
	- hexctl crd generate controller --kind Foo`,
	Run: func(cmd *cobra.Command, args []string) {
		crd.CreateController(&opts)
	},
}

func init() {
	initCmd.Flags().StringVarP(&opts.Domain, "domain", "d", "", "domain for groups (reqquired)")

	apiCmd.Flags().StringVarP(&opts.Group, "group", "g", "", "resource group (required)")
	apiCmd.Flags().StringVarP(&opts.Version, "version", "v", "", "resouce version (required)")
	apiCmd.Flags().StringVarP(&opts.Kind, "kind", "k", "", "resource kind (required)")
	apiCmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "attempt to create resource even if it already exists")
	apiCmd.Flags().BoolVarP(&opts.UseNamespace, "namespaced", "n", true, "resource is namespaced")
	apiCmd.Flags().BoolVarP(&opts.UseStatus, "status", "s", true, "resource is has status")

	ctrlCmd.Flags().StringVarP(&opts.Kind, "kind", "k", "", "resource kind (required)")
	ctrlCmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "attempt to create resource even if it already exists")

	crdCmd.PersistentFlags().StringVarP(&opts.Repo, "repo", "", "", "name to use for go module")

	initCmd.MarkFlagRequired("domain")

	apiCmd.MarkFlagRequired("group")
	apiCmd.MarkFlagRequired("version")
	apiCmd.MarkFlagRequired("kind")
	ctrlCmd.MarkFlagRequired("kind")

	generateCmd.AddCommand(apiCmd)
	generateCmd.AddCommand(ctrlCmd)

	crdCmd.AddCommand(initCmd)
	crdCmd.AddCommand(generateCmd)

	rootCmd.AddCommand(crdCmd)
}
