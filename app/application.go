package app

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"io"
	"net/url"
	"os"
	"path"
	"post-postman/app/config"
	"strings"
)

const (
	curlConfig    = "curl.cfg"
	requestConfig = "req.cfg"
)

type Runtime struct {
	pwd           string
	flags         *config.Flags
	curlConfig    config.Map
	requestConfig config.Map
}

func NewRuntime(root string, cmdFlags *pflag.FlagSet) (*Runtime, error) {
	flags, err := config.InitFlags(cmdFlags)
	if err != nil {
		return nil, err
	}

	return &Runtime{
		pwd:           root,
		flags:         flags,
		curlConfig:    config.Map{},
		requestConfig: config.Map{},
	}, nil
}

// Navigate navigates the runtime to the directory indicated by 'requests', starting from the application root.
// For example: Navigate("a", "b", "c") will change the working directory to '$APP_ROOT/a/b/c'
// It returns a RequestErr if the indicated path does not exist, or a FatalErr for any other error.
func (r *Runtime) Navigate(requests ...string) error {
	do := func(cfg *Runtime) error {
		return nil
	}

	return r.doAndNavigate(do, requests...)
}

// NavigateAndReadIn is like Navigate but it reads-in any configuration files as it navigates down the file tree.
// Finally, it reads in any configuration passed via command-line flags.
func (r *Runtime) NavigateAndReadIn(requests ...string) error {
	do := func(cfg *Runtime) error {
		return cfg.readInDir()
	}

	err := r.doAndNavigate(do, requests...)
	if err != nil {
		return err
	}

	r.readInFlags()
	return nil
}

// Create creates a new request-directory called "name" at the current PWD. It returns a RequestErr if the named request
// already exists, and a FatalErr for any other error.
func (r *Runtime) Create(name string) error {
	r.pwd = path.Join(r.pwd, name)
	err := os.Mkdir(r.pwd, os.ModePerm)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return RequestErr{
				msg:     "request already exists",
				request: name,
			}
		}
		return FatalErr{err}
	}
	return nil
}

// ReadIn reads in any configuration files in the current request-directory followed by any command-line flags.
func (r *Runtime) ReadIn() error {
	err := r.readInDir()
	if err != nil {
		return err
	}
	r.readInFlags()
	return nil
}

// WriteOut writes the current configuration to the current request-directory, overwriting any existing configuration
// files.
func (r *Runtime) WriteOut() error {
	err := writeConfig(r.curlConfig, path.Join(r.pwd, curlConfig))
	if err != nil {
		return err
	}

	err = writeConfig(r.requestConfig, path.Join(r.pwd, requestConfig))
	if err != nil {
		return err
	}

	return nil
}

// Remove all configuration files at the current request-directory, followed by the directory itself. Remove will fail
// if the current request has any children.
func (r *Runtime) Remove() error {
	children, err := r.getChildren()
	if err != nil {
		return err
	}

	if len(children) != 0 {
		return RequestErr{
			msg:     "cannot remove request with children",
			request: strings.Join(children, ","),
		}
	}

	err = os.RemoveAll(r.pwd)
	if err != nil {
		return FatalErr{err}
	}
	return nil
}

// ForceRemove is like Remove, but it removes a request along with all of its children.
func (r *Runtime) ForceRemove() error {
	err := os.RemoveAll(r.pwd)
	if err != nil {
		return FatalErr{err}
	}
	return nil
}

// Describe the request. Writes the description to stdOut.
func (r *Runtime) Describe() error {
	println("REQUEST CONFIG:")
	print(r.requestConfig.String())
	println()

	println("CURL CONFIG:")
	print(r.curlConfig.String())
	println()

	children, err := r.getChildren()
	if err != nil {
		return err
	}

	if len(children) == 0 {
		println("NO CHILDREN")
		return nil
	}

	println("CONTAINS:")
	println(strings.Join(children, ", "))
	return nil
}

// Cmd creates a list of arguments representing the request. The output of Cmd should be passed to exec.Cmd.
func (r *Runtime) Cmd() []string {
	flags := make([]string, 0)
	for k, v := range r.curlConfig {
		switch vv := v.(type) {
		case string:
			flags = append(flags, "--"+k)
			flags = append(flags, vv)
		case bool:
			if vv {
				flags = append(flags, "--"+k)
			}
		case []string:
			for _, s := range vv {
				flags = append(flags, "--"+k)
				flags = append(flags, s)
			}

		}
	}

	reqUrl, err := url.Parse("http://localhost")
	if err != nil {
		println("WARNING: Could not parse default URL")
	}

	rawUrl := r.requestConfig.GetString("url")
	if rawUrl != "" {
		userUrl, err := url.Parse(rawUrl)
		if err != nil {
			println("WARNING: Could not parse given URL")
		}
		reqUrl = userUrl
	}

	scheme := r.requestConfig.GetString("scheme")
	if scheme != "" {
		reqUrl.Scheme = scheme
	}

	host := r.requestConfig.GetString("host")
	if host != "" {
		reqUrl.Host = host
	}

	reqPath := r.requestConfig.GetStringSlice("path")
	if len(reqPath) > 0 {
		reqUrl.Path = path.Join(reqPath...)
	}

	return append(flags, reqUrl.String())
}

func (r *Runtime) doAndNavigate(do func(config *Runtime) error, requests ...string) error {
	err := do(r)
	if err != nil {
		return err
	}

	if len(requests) == 0 {
		return nil
	}

	r.pwd = path.Join(r.pwd, requests[0])
	_, err = os.Stat(r.pwd)
	if err != nil {
		return RequestErr{
			msg:     "no such request",
			request: requests[0],
		}
	}

	err = r.doAndNavigate(do, requests[1:]...)
	if err != nil {
		switch e := err.(type) {
		case RequestErr:
			e.request = fmt.Sprintf("%s.%s", requests[0], e.request)
			return e
		default:
			return err
		}
	}

	return nil
}

func (r *Runtime) readInDir() error {
	bytes, err := readFile(path.Join(r.pwd, curlConfig))
	if err != nil {
		return err
	}
	r.curlConfig.Merge(config.Read(bytes))

	bytes, err = readFile(path.Join(r.pwd, requestConfig))
	if err != nil {
		return err
	}
	r.requestConfig.Merge(config.Read(bytes))

	return nil
}

func (r *Runtime) readInFlags() {
	r.curlConfig.Merge(r.flags.CurlConfig())
	r.requestConfig.Merge(r.flags.RequestConfig())
}

func readFile(fileName string) ([]byte, error) {
	fh, err := os.OpenFile(fileName, os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, FatalErr{err}
	}
	defer func() { _ = fh.Close() }()

	bytes, err := io.ReadAll(fh)
	if err != nil {
		return nil, FatalErr{err}
	}

	return bytes, nil
}

func writeConfig(c config.Map, fileName string) error {
	fh, err := os.Create(fileName)
	if err != nil {
		return FatalErr{err}
	}
	defer func() { _ = fh.Close() }()

	_, err = fh.Write([]byte(c.String()))
	if err != nil {
		return FatalErr{err}
	}
	return nil
}

func (r *Runtime) getChildren() ([]string, error) {
	dirs, err := os.ReadDir(r.pwd)
	if err != nil {
		return nil, FatalErr{err}
	}

	children := make([]string, 0, len(dirs))
	for _, dir := range dirs {
		if dir.IsDir() {
			children = append(children, dir.Name())
		}
	}

	return children, nil
}
