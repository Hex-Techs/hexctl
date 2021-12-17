/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/Hex-Techs/hexctl/pkg/common/validate"
	"github.com/Hex-Techs/hexctl/pkg/kc"
	"github.com/spf13/cobra"
)

var kubeconfig string
var src string
var namespace bool

// kcCmd represents the kc command
var kcCmd = &cobra.Command{
	Use:              "kc",
	TraverseChildren: true,
	Short:            "manage your kubeconfig and context",
	Long: `kc helps you manage kubeconfig files and contexts.

- show current context:

  hexctl kc show

- get a context kubeconfig:

  hexctl kc get

you must have kubectl command already.`,
}

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "list all kube context",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := validate.ValidateArgs(args, -1)
		cobra.CheckErr(err)
		kc.Ls(kubeconfig)
	},
}

var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "switch your kube context",
	Run: func(cmd *cobra.Command, args []string) {
		cluster, err := validate.ValidateArgs(args, 0)
		if err != nil {
			if err.Error() != "args length is 0, but need 1" {
				cobra.CheckErr(err)
			}
			cluster = ""
		}
		kc.Switch(kubeconfig, cluster, namespace)
	},
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "show your current kube context",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := validate.ValidateArgs(args, -1)
		cobra.CheckErr(err)
		kc.Show(kubeconfig)
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete a context from kubeconfig",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := validate.ValidateArgs(args, -1)
		cobra.CheckErr(err)
		kc.Delete(kubeconfig)
	},
}

var nsCmd = &cobra.Command{
	Use:   "ns",
	Short: "switch your current kube context default namespace",
	Run: func(cmd *cobra.Command, args []string) {
		ns, err := validate.ValidateArgs(args, 0)
		if err != nil {
			if err.Error() != "args length is 0, but need 1" {
				cobra.CheckErr(err)
			}
			ns = ""
		}
		kc.Namespace(kubeconfig, ns)
	},
}

var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "merge tow kubeconfig file in ~/.kube/config",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := validate.ValidateArgs(args, -1)
		cobra.CheckErr(err)
		kc.Merge(src, kubeconfig)
	},
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get a context kubeconfig in kubeconfig",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := validate.ValidateArgs(args, -1)
		cobra.CheckErr(err)
		kc.GetContext(kubeconfig)
	},
}

func init() {
	mergeCmd.Flags().StringVarP(&src, "src", "s", "", "Specify the kubeconfig file to merge (required)")

	switchCmd.Flags().BoolVarP(&namespace, "namespace", "n", false, "Whether to switch namespace")

	kcCmd.PersistentFlags().StringVarP(&kubeconfig, "kubeconfig", "", "", "Specify the kubeconfig file to modify, default ~/.kube/config")

	mergeCmd.MarkFlagRequired("src")

	kcCmd.AddCommand(lsCmd)
	kcCmd.AddCommand(switchCmd)
	kcCmd.AddCommand(showCmd)
	kcCmd.AddCommand(deleteCmd)
	kcCmd.AddCommand(nsCmd)
	kcCmd.AddCommand(mergeCmd)
	kcCmd.AddCommand(getCmd)
	rootCmd.AddCommand(kcCmd)
}
