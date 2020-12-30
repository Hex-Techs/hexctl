package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var crdCmd = &cobra.Command{
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

var initCmd = &cobra.Command{
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

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Scaffold a Kubernetes API",
	Example: `	# Create a Foos API with Group: foo, Version: v1alpha1 and Kind: Foo 
	- hexctl crd create api --group foo --version v1alpha1 --kind Foo`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Scaffold a Kubernetes API",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update kubernetes API clients",
	Long: `After API changed, must need execute this command.
This command will update kubernetes API clients.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "build image for the project",
	Long:  ``,
	Example: `	#Build image
	- hexctl crd build examlpe.com/hex-techs/image:v0.0.0`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "push image to docker harbor",
	Long:  ``,
	Example: `	#Push image, it will push image that you made.
	- hexctl crd push`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Invokes a specific generator",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var k8sCmd = &cobra.Command{
	Use:   "k8s",
	Short: "generate k8s clients whit the API",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var openApiCmd = &cobra.Command{
	Use:   "openapi",
	Short: "generate k8s openAPI whit the API",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var crdRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run an custom project in a variety of environments",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// print k8s version, client-go version, crd version，go version, GOOS=darwin GOARCH=amd64
var crdVersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version of hexctl crd",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("hexctl crd version: %q, kubernetes version: %q, go version: %q, GOOS: %q, GOARCH: %q",
			crdVersion, kubernetesVersion, runtime.Version(), runtime.GOOS, runtime.GOARCH)
	},
}

func init() {
	crdCmd.AddCommand(initCmd)
	crdCmd.AddCommand(createCmd)
	crdCmd.AddCommand(updateCmd)
	crdCmd.AddCommand(buildCmd)
	crdCmd.AddCommand(pushCmd)
	crdCmd.AddCommand(generateCmd)
	crdCmd.AddCommand(crdRunCmd)
	crdCmd.AddCommand(crdVersionCmd)

	createCmd.AddCommand(apiCmd)
	generateCmd.AddCommand(k8sCmd)
	generateCmd.AddCommand(openApiCmd)

	rootCmd.AddCommand(crdCmd)
}
