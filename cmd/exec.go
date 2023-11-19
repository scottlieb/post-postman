package cmd

import (
	"github.com/spf13/cobra"
	execute "os/exec"
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
		err := cfg.NavigateDirAndReadIn()
		checkErr(err)

		curlCmd := execute.Command("curl", cfg.Cmd()...)
		println(curlCmd.String())
		out, err := curlCmd.CombinedOutput()
		println(string(out))
		checkErr(err)

		return
	}

	partsArg := args[0]
	if strings.Contains(partsArg, "/:;\\") {
		println("invalid request name")
		return
	}

	parts := strings.Split(partsArg, ".")
	err := cfg.NavigateDirAndReadIn(parts...)
	checkErr(err)

	curlCmd := execute.Command("curl", cfg.Cmd()...)
	println(curlCmd.String())
	out, err := curlCmd.CombinedOutput()
	println(string(out))
	checkErr(err)
}
