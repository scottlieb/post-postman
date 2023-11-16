package cmd

import (
	"github.com/spf13/cobra"
	"post-postman/internal/config"
	"strings"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec <request-name>",
	Short: "Execute a request.",
	Long:  "Execute a request.",
}

func init() {
	execCmd.Run = exec
}

func exec(_ *cobra.Command, args []string) {
	if len(args) == 0 {
		err := cmd.Help()
		cobra.CheckErr(err)
		return
	}

	partsArg := args[0]
	if strings.Contains(partsArg, "/:;\\") {
		println("invalid request name")
		return
	}

	parts := strings.Split(partsArg, ".")

	err := config.InitRequestConfig(parts...)
	cobra.CheckErr(err)

	req, err := config.Request()
	cobra.CheckErr(err)

	err = req.Execute()
	cobra.CheckErr(err)
}
