package cmd

import (
	"github.com/spf13/cobra"
	"post-postman/internal/config"
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

	err := config.CreateRequestConfig(parts...)
	cobra.CheckErr(err)

	err = config.Persist()
	cobra.CheckErr(err)
}
