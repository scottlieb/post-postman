package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"path"
	"post-postman/app"
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

var cfg app.Runtime

func init() {
	cmd.AddCommand(execCmd)
	cmd.AddCommand(createCmd)
	cmd.AddCommand(editCmd)
	cmd.AddCommand(removeCmd)
	cmd.AddCommand(describeCmd)
	cmd.AddCommand(resetCmd)

	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	appCfg, err := app.NewRuntime(path.Join(home, ".post-postman"), cmd.PersistentFlags())
	cobra.CheckErr(err)
	cfg = *appCfg
}
