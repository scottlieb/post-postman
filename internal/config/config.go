package config

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"path"
	"post-postman/internal/app"
)

const (
	Global      = "global"
	cfgFileName = "config.yaml"
)

var (
	root string
)

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}
	root = path.Join(home, ".post-postman")
}

func SetDefaults(flags *pflag.FlagSet) {
	setDefault(flags, "method", "X", http.MethodGet)
	setDefault(flags, "proto", "", "http")
	setDefault(flags, "host", "", "localhost")
	setDefault(flags, "path", "p", "")
	setDefault(flags, "content-type", "", "")
	setDefault(flags, "body", "d", "")

	viper.SetEnvPrefix("PPM")
	viper.AutomaticEnv()
}

func setDefault(flags *pflag.FlagSet, name, shorthand, value string) {
	description := ""
	if value != "" {
		description = fmt.Sprintf("default is '%s'", value)
	}

	if shorthand != "" {
		flags.StringP(name, shorthand, "", description)
	} else {
		flags.String(name, "", description)
	}

	err := viper.BindPFlag(name, flags.Lookup(name))
	if err != nil {
		panic(err)
	}

	if value != "" {
		viper.SetDefault(name, value)
	}
}

func InitConfig() error {
	err := createConfig()
	if err != nil {
		return errors.Wrap(err, "create global config")
	}

	err = readInConfig()
	if err != nil {
		return errors.Wrap(err, "init global config")
	}
	return nil
}

func CreateRequestConfig(parts ...string) error {
	for i, part := range parts {

		root = path.Join(root, part)

		// If we are at the last part, create it if it doesn't exist.
		if i == len(parts)-1 {
			err := createConfig()
			if err != nil {
				return errors.Wrap(err, "create request config")
			}
		}

		err := readInConfig()
		if err != nil {
			return errors.Wrap(err, "init request config")
		}
	}
	return nil
}

func InitRequestConfig(parts ...string) error {
	for _, part := range parts {
		root = path.Join(root, part)
		err := readInConfig()
		if err != nil {
			return errors.Wrap(err, "init request config")
		}
	}
	return nil
}

func createConfig() error {
	err := os.MkdirAll(root, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "create root dir")
	}
	cfgFilePath := path.Join(root, cfgFileName)
	fh, err := os.Create(cfgFilePath)
	if err != nil {
		return errors.Wrap(err, "create new config file")
	}
	err = fh.Close()
	if err != nil {
		return errors.Wrap(err, "close config file")
	}
	return nil
}

func readInConfig() error {
	cfgFilePath := path.Join(root, cfgFileName)
	viper.SetConfigFile(cfgFilePath)
	err := viper.MergeInConfig()
	if err != nil {
		return errors.Wrap(err, "read in config")
	}
	return nil
}

func Request() (app.Request, error) {
	res := app.Request{}
	err := viper.Unmarshal(&res)
	if err != nil {
		return app.Request{}, err
	}
	return res, nil
}

func Remove() error {
	return os.Remove(viper.ConfigFileUsed())
}

func Persist() error {
	return viper.WriteConfig()
}
