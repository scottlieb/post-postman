/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package edit

import (
	"github.com/spf13/viper"
	"post-postman/internal/entity"
	"strings"

	"github.com/spf13/cobra"
)

// CollectionCmd represents the collection command
var CollectionCmd = &cobra.Command{
	Use:     "collection",
	Aliases: []string{"col"},
	Short:   "TODO",
	Long:    "TODO",
}

func init() {
	CollectionCmd.Run = removeCollection
}

func removeCollection(_ *cobra.Command, args []string) {
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

	name := args[0]
	if strings.Contains(name, "/:;\\") {
		println("invalid request name")
		return
	}

	println("deleting", name)
	err := entity.RemoveCollection(root, name)
	cobra.CheckErr(err)
}
