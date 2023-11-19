package config

import (
	"fmt"
	"github.com/spf13/pflag"
	"io"
	"net/url"
	"os"
	"path"
)

type flagsCfg struct {
	changed func(string) bool
	curl    CurlFlags
	request RequestFlags
}

func (fc *flagsCfg) curlConfig() configMap {
	res := configMap{}
	res.Merge(flagsToConfig(fc.curl, fc.changed))
	return res
}

func (fc *flagsCfg) requestConfig() configMap {
	res := configMap{}
	res.Merge(flagsToConfig(fc.request, fc.changed))
	return res
}

func initApplicationFlags(flags *pflag.FlagSet) (*flagsCfg, error) {
	res := flagsCfg{
		changed: func(s string) bool {
			return flags.Changed(s)
		},
	}

	err := InitFlags(&res.curl, flags)
	if err != nil {
		return nil, err
	}

	err = InitFlags(&res.request, flags)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

const (
	curlConfig    = "curl.cfg"
	requestConfig = "req.cfg"
)

type ApplicationConfig struct {
	pwd        string
	flagsCfg   *flagsCfg
	curlCfg    configMap
	requestCfg configMap
}

func NewApplicationConfig(root string, flags *pflag.FlagSet) (*ApplicationConfig, error) {
	flagsCfg, err := initApplicationFlags(flags)
	if err != nil {
		return nil, err
	}

	return &ApplicationConfig{
		pwd:        root,
		flagsCfg:   flagsCfg,
		curlCfg:    configMap{},
		requestCfg: configMap{},
	}, nil
}

func (cfg *ApplicationConfig) Navigate(parts ...string) error {
	do := func(cfg *ApplicationConfig) error {
		return nil
	}

	return cfg.doAndNavigate(do, parts...)
}

func (cfg *ApplicationConfig) NavigateAndReadIn(parts ...string) error {
	do := func(cfg *ApplicationConfig) error {
		return cfg.ReadIn()
	}

	return cfg.doAndNavigate(do, parts...)
}

func (cfg *ApplicationConfig) doAndNavigate(do func(config *ApplicationConfig) error, parts ...string) error {
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

	err = cfg.doAndNavigate(do, parts[1:]...)
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

func (cfg *ApplicationConfig) Create(name string) error {
	cfg.pwd = path.Join(cfg.pwd, name)
	err := os.Mkdir(cfg.pwd, os.ModePerm)
	if err != nil {
		return FatalErr{err}
	}
	return nil
}

func (cfg *ApplicationConfig) ReadIn() error {
	// TODO request config
	bytes, err := readFile(path.Join(cfg.pwd, curlConfig))
	if err != nil {
		return err
	}
	cfg.curlCfg.Merge(readConfig(bytes))
	cfg.curlCfg.Merge(cfg.flagsCfg.curlConfig())

	bytes, err = readFile(path.Join(cfg.pwd, requestConfig))
	if err != nil {
		return err
	}
	cfg.requestCfg.Merge(readConfig(bytes))
	cfg.requestCfg.Merge(cfg.flagsCfg.requestConfig())

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

func (cfg *ApplicationConfig) WriteOut() error {
	err := writeConfig(cfg.curlCfg, path.Join(cfg.pwd, curlConfig))
	if err != nil {
		return err
	}

	err = writeConfig(cfg.requestCfg, path.Join(cfg.pwd, requestConfig))
	if err != nil {
		return err
	}

	return nil
}

func writeConfig(c configMap, fileName string) error {
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

func (cfg *ApplicationConfig) Cmd() []string {
	flags := make([]string, 0)
	for k, v := range cfg.curlCfg {
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

	rawUrl, ok := cfg.requestCfg["url"].(string)
	if ok {
		userUrl, err := url.Parse(rawUrl)
		if err != nil {
			println("WARNING: Could not parse given URL")
		}
		reqUrl = userUrl
	}

	scheme, ok := cfg.requestCfg["scheme"].(string)
	if ok {
		reqUrl.Scheme = scheme
	}

	host, ok := cfg.requestCfg["host"].(string)
	if ok {
		reqUrl.Host = host
	}

	reqPath, ok := cfg.requestCfg["path"].(string)
	if ok {
		reqUrl.Path = reqPath
	}

	return append(flags, reqUrl.String())
}
