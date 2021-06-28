package cmd

import (
	"github.com/Hex-Techs/hexctl/pkg/crd"
	"github.com/spf13/cobra"
)

type options struct {
	GVK        *crd.GVK
	WorkOption *crd.WorkOption
	Owner      string
	Repo       string
}

var opts = &options{
	GVK: &crd.GVK{
		Force: false,
	},
	WorkOption: &crd.WorkOption{
		Options: "",
	},
}

var crdCmd = &cobra.Command{
	Use:              "crd",
	TraverseChildren: true,
	Short:            "crd is a command scaffold for kubernetes extensions",
	Long: `Provides libraries and tools to create new projects, APIs and controllers.
Includes tools for packaging artifacts into an installer container.

Typical project lifecycle:

- initialize a project:

  hexctl crd init --domain example.com --owner "The Hex-Techs authors"

- create one or more a new resource APIs and add your code to them:

  hexctl crd create api --group <group> --version <version> --kind <Kind>

After the scaffold is written, api will run make on the project.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init kubernetes extensions project",
	Example: `	# Init project
	- hexctl crd init --domain example.com --owner "The Hex-Techs authors"`,
	Run: func(cmd *cobra.Command, args []string) {
		crd.Init(opts.Owner, opts.Repo, opts.GVK)
	},
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Scaffold a Kubernetes extensions api",
	Example: `	# Create a Foos API with Group: foo, Version: v1alpha1 and Kind: Foo
	- hexctl crd generate api --group foo --version v1alpha1 --kind Foo`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Scaffold a Kubernetes API",
	Run: func(cmd *cobra.Command, args []string) {
		crd.CreateAPI(opts.GVK)
	},
}

var codeCmd = &cobra.Command{
	Use:   "code",
	Short: "Scaffold a Kubernetes code",
	Long: `Generate necessary scripts and tool files for code generator.

Details: https://github.com/kubernetes/code-generator

You can use these tools to generate informer, client, lister, openapi and so on.`,
	Run: func(cmd *cobra.Command, args []string) {
		crd.Generate(opts.GVK, opts.WorkOption)
	},
}

func init() {
	initCmd.Flags().StringVarP(&opts.GVK.Domain, "domain", "d", "", "domain for groups (default \"my.domain\")")

	apiCmd.Flags().StringVarP(&opts.GVK.Group, "group", "g", "", "resource group (required)")
	apiCmd.Flags().StringVarP(&opts.GVK.Version, "version", "v", "", "resouce version (required)")
	apiCmd.Flags().StringVarP(&opts.GVK.Kind, "kind", "k", "", "resource kind (required)")
	apiCmd.Flags().BoolVarP(&opts.GVK.Force, "force", "f", false, "attempt to create resource even if it already exists")
	apiCmd.Flags().BoolVarP(&opts.GVK.UseNamespace, "namespaced", "n", true, "resource is namespaced")

	codeCmd.Flags().StringVarP(&opts.GVK.Group, "group", "g", "", "resource group (required)")
	codeCmd.Flags().StringVarP(&opts.GVK.Version, "version", "v", "", "resouce version (required)")
	codeCmd.Flags().StringVarP(&opts.WorkOption.Generated, "generated", "", "", "output package (required)")
	codeCmd.Flags().StringVarP(&opts.WorkOption.API, "apis", "", "", "apis package (required)")

	crdCmd.PersistentFlags().StringVarP(&opts.Owner, "owner", "", "", "specify the crd owner")
	crdCmd.PersistentFlags().StringVarP(&opts.Repo, "repo", "", "", "name to use for go module")

	apiCmd.MarkFlagRequired("group")
	apiCmd.MarkFlagRequired("version")
	apiCmd.MarkFlagRequired("kind")

	codeCmd.MarkFlagRequired("group")
	codeCmd.MarkFlagRequired("version")
	codeCmd.MarkFlagRequired("generated")
	codeCmd.MarkFlagRequired("apis")

	generateCmd.AddCommand(apiCmd)
	generateCmd.AddCommand(codeCmd)

	crdCmd.AddCommand(initCmd)
	crdCmd.AddCommand(generateCmd)

	rootCmd.AddCommand(crdCmd)
}
