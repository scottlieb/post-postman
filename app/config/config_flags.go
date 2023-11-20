package config

import (
	"github.com/spf13/pflag"
	"reflect"
	"strings"
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

type Flags struct {
	changed func(string) bool
	curlFlags
	requestFlags
}

func InitFlags(flags *pflag.FlagSet) (*Flags, error) {
	res := Flags{
		changed: func(s string) bool {
			return flags.Changed(s)
		},
	}

	err := initFlags(&res.curlFlags, flags)
	if err != nil {
		return nil, err
	}

	err = initFlags(&res.requestFlags, flags)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (f *Flags) CurlConfig() Map {
	res := Map{}
	res.Merge(fromFlags(f.curlFlags, f.changed))
	return res
}

func (f *Flags) RequestConfig() Map {
	res := Map{}
	res.Merge(fromFlags(f.requestFlags, f.changed))
	return res
}

// initFlags takes a generic struct and a FlagSet and inits the application
// flags bused on the struct tags.
func initFlags(flagStruct interface{}, flags *pflag.FlagSet) error {
	v := reflect.ValueOf(flagStruct).Elem()

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

func fromFlags(cfgStruct interface{}, changed func(string) bool) Map {
	res := Map{}
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

// header implements the pflag.Value interface
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
