package cmd

import (
	"github.com/spf13/cobra"
	"post-postman/internal/config"
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
}

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

	parts := strings.Split(partsArg, ".")

	err := config.InitRequestConfig(parts...)
	cobra.CheckErr(err)

	err = config.Remove()
	cobra.CheckErr(err)
}
