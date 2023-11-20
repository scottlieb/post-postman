package app

import "fmt"

type FatalErr struct {
	error
}

type RequestErr struct {
	msg     string
	request string
}

func (err RequestErr) Error() string {
	return fmt.Sprintf("%s: %s", err.msg, err.request)
}
