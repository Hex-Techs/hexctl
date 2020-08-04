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
	"github.com/Fize/n/pkg/kc"
	"github.com/spf13/cobra"
)

// kcCmd represents the kc command
var kcCmd = &cobra.Command{
	Use:   "kc",
	Short: "manage your kubeconfig and context",
	Long: `kc helps you manage kubeconfig files and contexts,
it will switch context or show current context.`,
}

var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "switch your kube context",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 {
			cmd.Help()
			return
		}
		kc.SwitchConfig()
	},
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "show your current kube context",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 {
			cmd.Help()
			return
		}
		kc.Show()
	},
}

func init() {
	kcCmd.AddCommand(switchCmd)
	kcCmd.AddCommand(showCmd)
	rootCmd.AddCommand(kcCmd)
}
