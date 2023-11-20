package cmd

import (
	"github.com/spf13/cobra"
	"os"
	execute "os/exec"
	"strings"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec <request-name>",
	Short: "Execute a request.",
	Long:  "Execute a request.",
}

var dryRun bool

func init() {
	execCmd.Run = exec

	execCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print cURL command without executing")
}

func exec(_ *cobra.Command, args []string) {
	if len(args) == 0 {
		err := cfg.NavigateAndReadIn()
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
	err := cfg.NavigateAndReadIn(parts...)
	checkErr(err)

	curlCmd := execute.Command("curl", cfg.Cmd()...)

	if dryRun {
		println(curlCmd.String())
		return
	}

	curlCmd.Stdout = os.Stdout
	curlCmd.Stderr = os.Stderr
	err = curlCmd.Run()
	checkErr(err)
}
