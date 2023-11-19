package config

import (
	"github.com/spf13/pflag"
	"reflect"
	"strings"
)

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

// InitFlags takes a generic struct and a FlagSet and inits the application
// flags bused on the struct tags.
func InitFlags(flagStruct interface{}, flags *pflag.FlagSet) error {
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

func FromFlags(cfgStruct interface{}, changed func(string) bool) Map {
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

func Read(in []byte) Map {
	res := Map{}
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
