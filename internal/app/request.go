package app

import (
	"fmt"
	"os/exec"
	"path"
	"strings"
)

const (
	AuthBearerToken = "bearer-token"
	AuthBasic       = "basic"
)

type Auth struct {
	Type        string
	BearerToken string
	User        string
	Password    string
}

type Header struct {
	Key string
	Val string
}

type Request struct {
	Auth        Auth
	Host        string
	PathPrefix  string
	Proto       string
	Method      string
	Path        string
	GlobHeaders []Header
	Headers     []Header
	ContentType string
	Body        string
}

func (r Request) Execute() error {
	cmd := exec.Command("curl", r.curl()...)
	fmt.Printf("Executing:\n%s\n", cmd)
	out, err := cmd.CombinedOutput()
	println(string(out))
	return err
}

func (r Request) url() string {
	// TODO path.Join is not good enough. Need to write a better function.
	noProto := path.Join(r.Host, r.Path)
	return fmt.Sprintf("%s://%s", r.Proto, noProto)
}

func (r Request) curl() []string {
	return strings.Split(fmt.Sprintf("-cfg -X%s %s", r.Method, r.url()), " ")
}
