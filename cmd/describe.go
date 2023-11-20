package cmd

import (
	"github.com/spf13/cobra"
	"strings"
)

// describeCmd represents the describe command
var describeCmd = &cobra.Command{
	Use:     "describe <name>",
	Aliases: []string{"desc"},
	Short:   "TODO",
	Long:    "TODO",
}

func init() {
	describeCmd.Run = describe
}

func describe(_ *cobra.Command, args []string) {
	if len(args) == 0 {
		err := cfg.ReadIn()
		checkErr(err)
		cfg.Describe()
		return
	}

	partsArg := args[0]
	if strings.Contains(partsArg, "/:;\\") {
		println("invalid request name")
		return
	}

	parts := strings.Split(partsArg, ".")
	err := cfg.Navigate(parts...)
	checkErr(err)

	err = cfg.ReadIn()
	checkErr(err)

	cfg.Describe()
}
