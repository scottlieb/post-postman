package cmd

import (
	"github.com/spf13/cobra"
	"post-postman/internal/config"
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

func edit(cmd *cobra.Command, args []string) {
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

	err = config.Persist()
	cobra.CheckErr(err)
}
