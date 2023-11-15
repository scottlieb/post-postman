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

// CollectionCmd represents the collection command
var CollectionCmd = &cobra.Command{
	Use:     "collection <name>",
	Aliases: []string{"col"},
	Short:   "Create a ppm collection",
	Long:    "Create a ppm collection",
}

func init() {
	CollectionCmd.Run = createCollection
}

func createCollection(_ *cobra.Command, args []string) {
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

	// TODO: Probably need better checks here:
	name := args[0]
	if strings.ContainsAny(name, "./:;\\") {
		println("invalid request name")
		return
	}

	err := entity.NewCollection(root, name)
	cobra.CheckErr(err)
}
