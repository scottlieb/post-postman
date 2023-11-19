package cmd

import (
	"github.com/spf13/cobra"
	"strings"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create <entity> <name>",
	Short: "TODO",
	Long:  "TODO",
}

func init() {
	createCmd.Run = create
}

func create(cmd *cobra.Command, args []string) {
	// TODO: read in flag values before write
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

	// All but the last part
	for i := 0; i < len(parts)-1; i++ {
		err := cfg.NavigateDir(parts...)
		checkErr(err)
	}

	err := cfg.CreateDir(parts[len(parts)-1])
	checkErr(err)

	err = cfg.WriteOut()
	checkErr(err)
}
