package cmd

import (
	"github.com/spf13/cobra"
	"strings"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit <entity> <name>",
	Short: "TODO",
	Long:  "TODO",
}

func init() {
	editCmd.Run = edit
}

func edit(_ *cobra.Command, args []string) {
	if len(args) == 0 {
		err := cfg.ReadIn()
		checkErr(err)
		err = cfg.WriteOut()
		checkErr(err)
		return
	}

	partsArg := args[0]
	if strings.Contains(partsArg, "/:;\\") {
		println("invalid request name")
		return
	}

	parts := strings.Split(partsArg, ".")
	err := cfg.Navigate(parts...)
	checkErr(err)

	err = cfg.ReadIn()
	checkErr(err)

	err = cfg.WriteOut()
	checkErr(err)
}
