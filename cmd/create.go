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

func create(_ *cobra.Command, args []string) {
	// TODO: read in flag values before write
	if len(args) == 0 {
		println("create requires a collection name")
		return
	}

	partsArg := args[0]
	if strings.Contains(partsArg, "/:;\\") {
		println("invalid request name")
		return
	}

	parts := strings.Split(partsArg, ".")

	// All but the last part
	err := cfg.Navigate(parts[:len(parts)-1]...)
	checkErr(err)

	err = cfg.Create(parts[len(parts)-1])
	checkErr(err)

	err = cfg.ReadIn()
	checkErr(err)

	err = cfg.WriteOut()
	checkErr(err)
}
