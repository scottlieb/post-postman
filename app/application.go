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

func (r *Runtime) NavigateDir(parts ...string) error {
	do := func(cfg *Runtime) error {
		return nil
	}

	return r.doAndNavigateDir(do, parts...)
}

func (r *Runtime) NavigateDirAndReadIn(parts ...string) error {
	do := func(cfg *Runtime) error {
		return cfg.ReadIn()
	}

	return r.doAndNavigateDir(do, parts...)
}

func (r *Runtime) doAndNavigateDir(do func(config *Runtime) error, parts ...string) error {
	err := do(r)
	if err != nil {
		return err
	}

	if len(parts) == 0 {
		return nil
	}

	r.pwd = path.Join(r.pwd, parts[0])
	_, err = os.Stat(r.pwd)
	if err != nil {
		return CollectionErr{
			msg:        "no such collection",
			collection: parts[0],
		}
	}

	err = r.doAndNavigateDir(do, parts[1:]...)
	if err != nil {
		switch e := err.(type) {
		case CollectionErr:
			e.collection = fmt.Sprintf("%s.%s", parts[0], e.collection)
			return e
		default:
			return err
		}
	}

	return nil
}

func (r *Runtime) CreateDir(name string) error {
	r.pwd = path.Join(r.pwd, name)
	err := os.Mkdir(r.pwd, os.ModePerm)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return CollectionErr{
				msg:        "collection already exists",
				collection: name,
			}
		}
		return FatalErr{err}
	}
	return nil
}

func (r *Runtime) ReadIn() error {
	// TODO request config
	bytes, err := readFile(path.Join(r.pwd, curlConfig))
	if err != nil {
		return err
	}
	r.curlConfig.Merge(config.Read(bytes))
	r.curlConfig.Merge(r.flags.CurlConfig())

	bytes, err = readFile(path.Join(r.pwd, requestConfig))
	if err != nil {
		return err
	}
	r.requestConfig.Merge(config.Read(bytes))
	r.requestConfig.Merge(r.flags.RequestConfig())

	return nil
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

func (r *Runtime) Cmd() []string {
	flags := make([]string, 0)
	for k, v := range r.curlConfig {
		switch vv := v.(type) {
		case string:
			flags = append(flags, "--"+k)
			flags = append(flags, vv)
		case bool:
			flags = append(flags, "--"+k)
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

	rawUrl, ok := r.requestConfig.GetString("url")
	if ok {
		userUrl, err := url.Parse(rawUrl)
		if err != nil {
			println("WARNING: Could not parse given URL")
		}
		reqUrl = userUrl
	}

	scheme, ok := r.requestConfig.GetString("scheme")
	if ok {
		reqUrl.Scheme = scheme
	}

	host, ok := r.requestConfig.GetString("host")
	if ok {
		reqUrl.Host = host
	}

	reqPath, ok := r.requestConfig.GetString("path")
	if ok {
		reqUrl.Path = reqPath
	}

	return append(flags, reqUrl.String())
}

func (r *Runtime) ForceRemove() error {
	err := os.RemoveAll(r.pwd)
	if err != nil {
		return FatalErr{err}
	}
	return nil
}

func (r *Runtime) Remove() error {
	children, err := r.getChildren()
	if err != nil {
		return err
	}

	if len(children) != 0 {
		return CollectionErr{
			msg:        "cannot remove collection with children",
			collection: strings.Join(children, ","),
		}
	}

	err = os.RemoveAll(r.pwd)
	if err != nil {
		return FatalErr{err}
	}
	return nil
}

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
