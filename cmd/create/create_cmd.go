/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package create

import (
	"github.com/spf13/cobra"
)

// Cmd represents the create command
var Cmd = &cobra.Command{
	Use:   "create <entity> <name>",
	Short: "Create a new pmm entity.",
	Long: `Create a new collection or request.
Example Usages:
pmm create collection my-api
pmm create request my-request --path '/my/api/request'`,
}

func init() {
	Cmd.AddCommand(CollectionCmd)
	Cmd.AddCommand(RequestCmd)
}
