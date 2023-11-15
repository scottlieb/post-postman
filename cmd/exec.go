/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"post-postman/internal/entity"
	"strings"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec <request-name>",
	Short: "Execute a request.",
	Long:  "Execute a request.",
}

func init() {
	Cmd.AddCommand(execCmd)
	execCmd.Run = exec
}

func exec(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		err := execCmd.Help()
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

	println("executing request:", name, "using collection:", collection)

	err := entity.ExecuteRequest(root, collection, name)
	cobra.CheckErr(err)
}
