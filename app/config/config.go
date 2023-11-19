package config

import (
	"fmt"
	"github.com/spf13/pflag"
	"reflect"
	"strings"
)

type RequestFlags struct {
	URL    string `short:"u" desc:"Request URL"`
	Scheme string `short:"s" desc:"HTTP scheme to use"`
	Host   string `short:"y" desc:"Request host"`
	Path   string `short:"p" desc:"Request path"`
}

type CurlFlags struct {
	Data    string   `short:"d" desc:"HTTP POST data"`
	Request string   `short:"X" desc:"Specify request method to use" default:"GET"`
	Verbose bool     `short:"v" desc:"Make the operation more talkative"`
	Header  []string `short:"H" desc:"Pass custom header(s) to server"`
}

type configMap map[string]interface{}

type header struct {
	val *[]string
}

func (h header) String() string {
	return ""
}

func (h header) Set(s string) error {
	*(h.val) = append(*(h.val), s)
	return nil
}

func (h header) Type() string {
	return "key:value"
}

func InitFlags(cfgStruct interface{}, flags *pflag.FlagSet) error {
	v := reflect.ValueOf(cfgStruct).Elem()

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		name := strings.ToLower(field.Name)
		tag := field.Tag
		value := v.Field(i).Addr().Interface()
		shorthand := tag.Get("short")
		description := tag.Get("desc")
		defaultVal := tag.Get("default")

		switch vv := value.(type) {
		case *string:
			flags.StringVarP(vv, name, shorthand, defaultVal, description)
		case *bool:
			flags.BoolVarP(vv, name, shorthand, false, description)
		case *[]string:
			flags.VarP(header{vv}, name, shorthand, description)
		}
	}

	return nil
}

func flagsToConfig(cfgStruct interface{}, changed func(string) bool) configMap {
	res := configMap{}
	v := reflect.ValueOf(cfgStruct)

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		name := strings.ToLower(field.Name)
		if !changed(name) {
			continue
		}

		value := v.Field(i).Interface()
		res[name] = value
	}

	return res
}

func readConfig(in []byte) configMap {
	res := configMap{}
	fields := strings.Fields(string(in))
	for _, field := range fields {
		parts := strings.Split(field, "=")
		if len(parts) == 1 {
			res[field] = true
		}
		if len(parts) == 2 && parts[0] == "header" {
			hdrs, ok := res["header"]
			if ok {
				res["header"] = append(hdrs.([]string), parts[1])
				continue
			}
			res["header"] = []string{parts[1]}
			continue
		}
		if len(parts) == 2 {
			res[parts[0]] = parts[1]
		}
		if len(parts) > 2 {
			println("Warning: bad config field:", field)
		}
	}
	return res
}

func (c configMap) Merge(other configMap) {
	for k, v := range other {
		if k == "header" {
			hdrs, ok := c[k]
			if ok {
				c[k] = append(hdrs.([]string), v.([]string)...)
				continue
			}
			c[k] = v.([]string)
		}

		c[k] = v
	}
}

func (c configMap) String() string {
	res := ""
	for k, v := range c {
		switch vv := v.(type) {
		case string:
			res += fmt.Sprintf("%s=%s\n", k, vv)
		case bool:
			res += fmt.Sprintf("%s\n", k)
		case []string:
			for _, s := range vv {
				res += fmt.Sprintf("%s=%s\n", k, s)
			}
		}
	}
	return res
}
