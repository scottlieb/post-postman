/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package create

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"post-postman/internal/entity"
	"strings"
)

// RequestCmd represents the collection command
var RequestCmd = &cobra.Command{
	Use:     "request <name>",
	Aliases: []string{"req"},
	Short:   "Create a ppm request",
	Long:    "Create a ppm request",
}

func init() {
	RequestCmd.Run = createRequest
}

func createRequest(_ *cobra.Command, args []string) {
	if len(args) == 0 {
		err := CollectionCmd.Help()
		cobra.CheckErr(err)
		return
	}

	root := viper.GetString("root")
	if root == "" {
		println("no project root configured")
		return
	}

	collection := viper.GetString("collection")

	// TODO: Probably need better checks here:
	name := args[0]
	if strings.ContainsAny(name, "./:;\\") {
		println("invalid request name")
		return
	}

	err := entity.NewRequest(root, collection, name)
	cobra.CheckErr(err)
}
