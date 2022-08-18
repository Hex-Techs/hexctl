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
	"github.com/Hex-Techs/hexctl/pkg/kubeconfig"
	"github.com/spf13/cobra"
)

var (
	kubeconfigPath string
	src            string
	dst            string
	cascade        bool
	byteFormat     bool
)

// kcCmd represents the kc command
var kcCmd = &cobra.Command{
	Use:              "kubeconfig",
	Aliases:          []string{"kc"},
	TraverseChildren: true,
	Short:            "manage your kubeconfig and context",
	Long: `kc helps you manage kubeconfig files and contexts.

- show current context:

  hexctl kc show

- get a context kubeconfig:

  hexctl kc get

you must have kubectl command already.`,
}

var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "merge tow kubeconfig file in ~/.kube/config",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := validate.ValidateArgs(args, -1)
		cobra.CheckErr(err)
		mgr, err := kubeconfig.NewKCMgr(kubeconfigPath)
		cobra.CheckErr(err)
		cobra.CheckErr(mgr.MergeContext(src))
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete a context from kubeconfig",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := validate.ValidateArgs(args, -1)
		cobra.CheckErr(err)
		mgr, err := kubeconfig.NewKCMgr(kubeconfigPath)
		cobra.CheckErr(err)
		cobra.CheckErr(mgr.DeleteContext())
	},
}

var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "switch your kube context",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := validate.ValidateArgs(args, -1)
		cobra.CheckErr(err)
		mgr, err := kubeconfig.NewKCMgr(kubeconfigPath)
		cobra.CheckErr(err)
		cobra.CheckErr(mgr.SwitchContext(cascade))
	},
}

var switchNsCmd = &cobra.Command{
	Use:     "namespace",
	Aliases: []string{"ns"},
	Short:   "switch your current kube context default namespace",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := validate.ValidateArgs(args, -1)
		cobra.CheckErr(err)
		mgr, err := kubeconfig.NewKCMgr(kubeconfigPath)
		cobra.CheckErr(err)
		cobra.CheckErr(mgr.SwitchNamespace())
	},
}

var rename = &cobra.Command{
	Use:   "rename",
	Short: "rename a context from kubeconfig",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := validate.ValidateArgs(args, -1)
		cobra.CheckErr(err)
		mgr, err := kubeconfig.NewKCMgr(kubeconfigPath)
		cobra.CheckErr(err)
		cobra.CheckErr(mgr.RenameContext())
	},
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get a context kubeconfig in kubeconfig",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := validate.ValidateArgs(args, -1)
		cobra.CheckErr(err)
		mgr, err := kubeconfig.NewKCMgr(kubeconfigPath)
		cobra.CheckErr(err)
		cobra.CheckErr(mgr.GetContext(dst, byteFormat))
	},
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "list all kube context",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := validate.ValidateArgs(args, -1)
		cobra.CheckErr(err)
		mgr, err := kubeconfig.NewKCMgr(kubeconfigPath)
		cobra.CheckErr(err)
		mgr.ListContext()
	},
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "show your current kube context",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := validate.ValidateArgs(args, -1)
		cobra.CheckErr(err)
		mgr, err := kubeconfig.NewKCMgr(kubeconfigPath)
		cobra.CheckErr(err)
		mgr.ShowCurrentContext()
	},
}

func init() {
	mergeCmd.Flags().StringVarP(&src, "src", "s", "", "Specify the kubeconfig file to merge (required)")

	switchCmd.Flags().BoolVarP(&cascade, "cascade", "c", false, "Whether to switch namespace")

	kcCmd.PersistentFlags().StringVarP(&kubeconfigPath, "kubeconfig", "", "", "Specify the kubeconfig file to modify, default ~/.kube/config")

	getCmd.Flags().StringVarP(&dst, "dst", "d", "", "Stores the contents in the specified file")
	getCmd.Flags().BoolVarP(&byteFormat, "byte", "b", false, "Whether to output the contents in byte format")

	mergeCmd.MarkFlagRequired("src")

	kcCmd.AddCommand(listCmd, switchCmd, showCmd, switchNsCmd, mergeCmd, deleteCmd, getCmd, rename)
	rootCmd.AddCommand(kcCmd)
}
