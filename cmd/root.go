/*
Copyright © 2023 Amitai Gottlieb
*/
package cmd

import (
	"os"
	"path"
	"post-postman/cmd/create"
	"post-postman/cmd/edit"
	"post-postman/cmd/remove"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Cmd represents the base command when called without any subcommands
var Cmd = &cobra.Command{
	Use:   "ppm",
	Short: "A CLI REST client with support for collections.",
	Long:  "Welcome to post-postman! A CLI REST client with support for collections!",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the Cmd.
func Execute() {
	err := Cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	Cmd.AddCommand(create.Cmd)
	Cmd.AddCommand(remove.Cmd)
	Cmd.AddCommand(edit.Cmd)

	cobra.OnInitialize(initConfig)

	Cmd.PersistentFlags().String("config", "", "config file (default is $HOME/.post-postman/config)")
	viper.SetDefault("cfgFile", "")
	err := viper.BindPFlag("cfgFile", Cmd.PersistentFlags().Lookup("config"))
	cobra.CheckErr(err)

	Cmd.PersistentFlags().StringP("collection", "c", "", "collection to use (default is 'global')")
	viper.SetDefault("collection", "global")
	err = viper.BindPFlag("collection", Cmd.PersistentFlags().Lookup("collection"))
	cobra.CheckErr(err)

	// TODO: Make this better.
	cobra.OnFinalize(func() { _ = viper.SafeWriteConfig() })
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	defaultRoot := path.Join(home, ".post-postman")
	viper.SetDefault("root", defaultRoot)

	if cfgFile := viper.GetString("cfgFile"); cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".post-postman" (without extension).
		viper.AddConfigPath(defaultRoot)
		viper.SetConfigType("json")
		viper.SetConfigName("config")
	}

	// Read in environment variables that match.
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	_ = viper.ReadInConfig()
}
