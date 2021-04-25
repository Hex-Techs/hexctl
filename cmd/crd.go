package cmd

import (
	"github.com/Hex-Techs/hexctl/cmd/app/crd"
)

func init() {
	rootCmd.AddCommand(crd.NewCrdCommand())
}
