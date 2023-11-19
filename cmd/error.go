package cmd

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"post-postman/app"
)

func checkErr(err error) {
	if err == nil {
		return
	}

	if errors.Is(err, app.FatalErr{}) {
		_, _ = fmt.Fprintf(os.Stderr, "FATAL: %v", err)
		os.Exit(1)
	}

	fmt.Println(err.Error())
	os.Exit(0)
}
