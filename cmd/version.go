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
	"runtime"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

// ___ ___                ___________           .__
// /   |   \   ____ ___  __\__    ___/___   ____ |  |__
// /    ~    \_/ __ \\  \/  / |    |_/ __ \_/ ___\|  |  \
// \    Y    /\  ___/ >    <  |    |\  ___/\  \___|   Y  \
// \___|_  /  \___  >__/\_ \ |____| \___  >\___  >___|  /
// \/       \/      \/            \/     \/     \/

const (
	version = "v0.1.6"
	Logo    = `  ___ ___                ___________           .__
 /   |   \   ____ ___  __\__    ___/___   ____ |  |__
/    ~    \_/ __ \\  \/  / |    |_/ __ \_/ ___\|  |  \
\    Y    /\  ___/ >    <  |    |\  ___/\  \___|   Y  \
 \___|_  /  \___  >__/\_ \ |____| \___  >\___  >___|  /
       \/       \/      \/            \/     \/     \/  `
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "hexctl version",
	Long:  `show hexctl version.`,
	Run: func(cmd *cobra.Command, args []string) {
		color.Printf("%s\n%s\n%s/%s\n", Logo, version, runtime.GOOS, runtime.GOARCH)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
