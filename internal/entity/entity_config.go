package entity

import (
	"fmt"
	"github.com/spf13/viper"
	"path"
	"strings"
)

const (
	Global          = "global"
	CfgFileName     = "config"
	CfgFileNameJSON = CfgFileName + ".json"
)

type CollectionConfig struct {
	Auth        Auth     `json:"auth,omitempty"`
	Proto       string   `json:"proto,omitempty"`
	Host        string   `json:"host,omitempty"`
	GlobHeaders []Header `json:"globHeaders,omitempty"`
}

type ReqOpt func(RequestConfig) RequestConfig

type RequestConfig struct {
	CollectionConfig `mapstructure:",squash"`
	Method           string   `json:"method,omitempty"`
	Path             string   `json:"path,omitempty"`
	Headers          []Header `json:"headers,omitempty"`
	ContentType      string   `json:"contentType,omitempty"`
	Body             string   `json:"body,omitempty"`
}

func (c RequestConfig) toUrl() string {
	// TODO path.Join is not good enough. Need to write a better function.
	noProto := path.Join(c.Host, c.Path)
	return fmt.Sprintf("%s://%s", c.Proto, noProto)
}

func (c RequestConfig) toCurl() []string {
	return strings.Split(fmt.Sprintf("-v -X%s %s", c.Method, c.toUrl()), " ")
}

const (
	AuthBearerToken = "bearer-token"
	AuthBasic       = "basic"
)

type Auth struct {
	Type        string `json:"type,omitempty"`
	BearerToken string `json:"bearerToken,omitempty"`
	User        string `json:"user,omitempty"`
	Password    string `json:"password,omitempty"`
}

type Header struct {
	Key string
	Val string
}

func newDefaultConfig() *viper.Viper {
	cfg := viper.New()
	cfg.SetDefault("method", "GET")
	cfg.SetDefault("proto", "http")
	return cfg
}
