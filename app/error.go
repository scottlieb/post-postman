package app

import "fmt"

type FatalErr struct {
	error
}

type CollectionErr struct {
	msg        string
	collection string
}

func (c CollectionErr) Error() string {
	return fmt.Sprintf("%s: %s", c.msg, c.collection)
}
