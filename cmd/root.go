package cmd

import (
	"post-postman/internal/config"

	"github.com/spf13/cobra"
)

// cmd represents the base command when called without any subcommands
var cmd = &cobra.Command{
	Use:   "ppm",
	Short: "A CLI REST client with support for nested requests",
	Long: `Welcome to post-postman!
Post-postman is a CLI REST client with support for nested requests. It was built as
a free, open-source and private alternative to the popular tool 'postman'.`,
}

func Execute() {
	cobra.CheckErr(cmd.Execute())
}

var collection = config.Global

func init() {
	cmd.AddCommand(execCmd)
	cmd.AddCommand(createCmd)
	cmd.AddCommand(editCmd)
	cmd.AddCommand(removeCmd)

	config.SetDefaults(cmd.PersistentFlags())

	cmd.PersistentFlags().StringVarP(&collection, "collection", "c", "", "")

	cobra.OnInitialize(func() {
		err := config.InitConfig()
		cobra.CheckErr(err)
	})
}
