/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package remove

import (
	"github.com/spf13/cobra"
)

// Cmd represents the delete command
var Cmd = &cobra.Command{
	Use:     "delete <entity> <name>",
	Aliases: []string{"del", "remove", "rm"},
	Short:   "TODO",
	Long:    "TODO",
}

func init() {
	Cmd.AddCommand(CollectionCmd)
	Cmd.AddCommand(RequestCmd)
}
