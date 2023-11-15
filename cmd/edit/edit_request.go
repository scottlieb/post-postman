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

// RequestCmd represents the collection command
var RequestCmd = &cobra.Command{
	Use:     "request",
	Aliases: []string{"req"},
	Short:   "TODO",
	Long:    "TODO",
}

var (
	method      string
	host        string
	path        string
	contentType string
	body        string
)

func init() {
	RequestCmd.Run = editRequest

	RequestCmd.Flags().StringVarP(&method, "method", "X", "", "HTTP method")
	RequestCmd.Flags().StringVar(&host, "host", "", "Host of the request")
	RequestCmd.Flags().StringVarP(&path, "path", "p", "", "Path of the request")
}

func editRequest(_ *cobra.Command, args []string) {
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

	name := args[0]
	if strings.Contains(name, "/:;\\") {
		println("invalid request name")
		return
	}

	err := entity.EditRequest(
		root,
		collection,
		name,
		func(cfg entity.RequestConfig) entity.RequestConfig {
			if method != "" {
				cfg.Method = method
			}
			return cfg
		},
		func(cfg entity.RequestConfig) entity.RequestConfig {
			if host != "" {
				cfg.Host = host
			}
			return cfg
		},
		func(cfg entity.RequestConfig) entity.RequestConfig {
			if path != "" {
				cfg.Path = path
			}
			return cfg
		},
	)
	cobra.CheckErr(err)
}
