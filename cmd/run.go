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
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/Fize/n/pkg/output"

	"github.com/Fize/n/pkg/run"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run a go project",
	Long: `run a go project. For example:

Gin or other web project, it will watch *.go file and when these file changed n will reload it,
you must have a main.go file in workdir and code in the directory named pkg.`,
	Run: func(cmd *cobra.Command, args []string) {
		stop := make(chan bool)
		pwd, _ := os.Getwd()
		pkgs, err := run.GetDirList(filepath.Join(pwd, "pkg"))
		cmds, err := run.GetDirList(filepath.Join(pwd, "cmd"))
		dirs := append(pkgs, cmds...)
		if err != nil {
			output.Fatalln(err)
		}
		go run.NewWatcher(dirs, stop)
		go run.Reload(command, stop)
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		for {
			select {
			case <-sigs:
				run.Kill()
				os.Exit(0)
			}
		}
	},
}

var command []string

func init() {
	runCmd.Flags().StringSliceVarP(&command, "cmd", "", []string{}, "app command")
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
