package cmd

import (
	"github.com/spf13/cobra"
	"strings"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:     "remove <entity> <name>",
	Aliases: []string{"rm"},
	Short:   "TODO",
	Long:    "TODO",
}

func init() {
	removeCmd.Run = remove

	removeCmd.Flags().BoolVarP(&force, "force", "f", false, "force remove")
}

var force bool

func remove(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		err := cmd.Help()
		cobra.CheckErr(err)
		return
	}

	partsArg := args[0]
	if strings.Contains(partsArg, "/:;\\") {
		println("invalid request name")
		return
	}
}
