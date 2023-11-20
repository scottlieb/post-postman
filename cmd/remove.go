package cmd

import (
	"github.com/spf13/cobra"
	"strings"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:     "remove <name>",
	Aliases: []string{"rm"},
	Args:    cobra.ExactArgs(1),
	Short:   "TODO",
	Long:    "TODO",
}

func init() {
	removeCmd.Run = remove
	removeCmd.Flags().BoolVarP(&force, "force", "f", false, "force remove")
}

var force bool

func remove(_ *cobra.Command, args []string) {
	partsArg := args[0]
	if strings.Contains(partsArg, "/:;\\") {
		println("invalid request name")
		return
	}

	parts := strings.Split(partsArg, ".")

	err := cfg.Navigate(parts...)
	checkErr(err)

	if force {
		err = cfg.ForceRemove()
	} else {
		err = cfg.Remove()
	}
	checkErr(err)
}
