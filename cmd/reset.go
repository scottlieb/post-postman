package cmd

import (
	"github.com/spf13/cobra"
	"strings"
)

// resetCmd represents the reset command
var resetCmd = &cobra.Command{
	Use:   "reset <name>",
	Short: "TODO",
	Long:  "TODO",
}

func init() {
	resetCmd.Run = reset
}

func reset(_ *cobra.Command, args []string) {
	if len(args) == 0 {
		err := cfg.WriteOut()
		checkErr(err)
		return
	}

	partsArg := args[0]
	if strings.Contains(partsArg, "/:;\\") {
		println("invalid request name")
		return
	}

	parts := strings.Split(partsArg, ".")
	err := cfg.NavigateDir(parts...)
	checkErr(err)

	err = cfg.WriteOut()
	checkErr(err)
}
