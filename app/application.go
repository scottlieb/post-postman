package app

import (
	"fmt"
	"github.com/spf13/pflag"
	"io"
	"net/url"
	"os"
	"path"
	"post-postman/app/config"
)

type requestFlags struct {
	URL    string `short:"u" desc:"Request URL"`
	Scheme string `short:"s" desc:"HTTP scheme to use"`
	Host   string `short:"y" desc:"Request host"`
	Path   string `short:"p" desc:"Request path"`
}

type curlFlags struct {
	Data    string   `short:"d" desc:"HTTP POST data"`
	Request string   `short:"X" desc:"Specify request method to use" default:"GET"`
	Verbose bool     `short:"v" desc:"Make the operation more talkative"`
	Header  []string `short:"H" desc:"Pass custom header(s) to server"`
}

type applicationFlags struct {
	changed func(string) bool
	curlFlags
	requestFlags
}

func (fc *applicationFlags) curlConfig() config.Map {
	res := config.Map{}
	res.Merge(config.FromFlags(fc.curlFlags, fc.changed))
	return res
}

func (fc *applicationFlags) requestConfig() config.Map {
	res := config.Map{}
	res.Merge(config.FromFlags(fc.requestFlags, fc.changed))
	return res
}

func initApplicationFlags(flags *pflag.FlagSet) (*applicationFlags, error) {
	res := applicationFlags{
		changed: func(s string) bool {
			return flags.Changed(s)
		},
	}

	err := config.InitFlags(&res.curlFlags, flags)
	if err != nil {
		return nil, err
	}

	err = config.InitFlags(&res.requestFlags, flags)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

const (
	curlConfig    = "curl.cfg"
	requestConfig = "req.cfg"
)

type Runtime struct {
	pwd           string
	flags         *applicationFlags
	curlConfig    config.Map
	requestConfig config.Map
}

func NewRuntime(root string, cmdFlags *pflag.FlagSet) (*Runtime, error) {
	flags, err := initApplicationFlags(cmdFlags)
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

func (cfg *Runtime) NavigateDir(parts ...string) error {
	do := func(cfg *Runtime) error {
		return nil
	}

	return cfg.doAndNavigateDir(do, parts...)
}

func (cfg *Runtime) NavigateDirAndReadIn(parts ...string) error {
	do := func(cfg *Runtime) error {
		return cfg.ReadIn()
	}

	return cfg.doAndNavigateDir(do, parts...)
}

func (cfg *Runtime) doAndNavigateDir(do func(config *Runtime) error, parts ...string) error {
	err := do(cfg)
	if err != nil {
		return err
	}

	if len(parts) == 0 {
		return nil
	}

	cfg.pwd = path.Join(cfg.pwd, parts[0])
	_, err = os.Stat(cfg.pwd)
	if err != nil {
		return CollectionErr{
			msg:        "no such collection",
			collection: parts[0],
		}
	}

	err = cfg.doAndNavigateDir(do, parts[1:]...)
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

func (cfg *Runtime) CreateDir(name string) error {
	cfg.pwd = path.Join(cfg.pwd, name)
	err := os.Mkdir(cfg.pwd, os.ModePerm)
	if err != nil {
		return FatalErr{err}
	}
	return nil
}

func (cfg *Runtime) ReadIn() error {
	// TODO request config
	bytes, err := readFile(path.Join(cfg.pwd, curlConfig))
	if err != nil {
		return err
	}
	cfg.curlConfig.Merge(config.Read(bytes))
	cfg.curlConfig.Merge(cfg.flags.curlConfig())

	bytes, err = readFile(path.Join(cfg.pwd, requestConfig))
	if err != nil {
		return err
	}
	cfg.requestConfig.Merge(config.Read(bytes))
	cfg.requestConfig.Merge(cfg.flags.requestConfig())

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

func (cfg *Runtime) WriteOut() error {
	err := writeConfig(cfg.curlConfig, path.Join(cfg.pwd, curlConfig))
	if err != nil {
		return err
	}

	err = writeConfig(cfg.requestConfig, path.Join(cfg.pwd, requestConfig))
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

func (cfg *Runtime) Cmd() []string {
	flags := make([]string, 0)
	for k, v := range cfg.curlConfig {
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

	rawUrl, ok := cfg.requestConfig.GetString("url")
	if ok {
		userUrl, err := url.Parse(rawUrl)
		if err != nil {
			println("WARNING: Could not parse given URL")
		}
		reqUrl = userUrl
	}

	scheme, ok := cfg.requestConfig.GetString("scheme")
	if ok {
		reqUrl.Scheme = scheme
	}

	host, ok := cfg.requestConfig.GetString("host")
	if ok {
		reqUrl.Host = host
	}

	reqPath, ok := cfg.requestConfig.GetString("path")
	if ok {
		reqUrl.Path = reqPath
	}

	return append(flags, reqUrl.String())
}
