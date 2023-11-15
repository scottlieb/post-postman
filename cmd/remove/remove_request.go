/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package remove

import (
	"github.com/spf13/viper"
	"post-postman/internal/entity"
	"strings"

	"github.com/spf13/cobra"
)

// RequestCmd represents the collection command
var RequestCmd = &cobra.Command{
	Use:     "request",
	Aliases: []string{"req"},
	Short:   "TODO",
	Long:    "TODO",
}

func init() {
	RequestCmd.Run = removeRequest
}

func removeRequest(_ *cobra.Command, args []string) {
	if len(args) == 0 {
		err := RequestCmd.Help()
		cobra.CheckErr(err)
		return
	}

	root := viper.GetString("root")
	if root == "" {
		println("no project root configured")
		return
	}

	collection := viper.GetString("collection")

	name := args[0]
	if strings.Contains(name, "/:;\\") {
		println("invalid request name")
		return
	}

	println("deleting", name)
	err := entity.RemoveRequest(root, collection, name)
	cobra.CheckErr(err)
}
