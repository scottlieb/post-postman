package config

import (
	"fmt"
	"strings"
)

type Map map[string]interface{}

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

func (c Map) GetString(key string) (string, bool) {
	res, ok := c[key]
	if !ok {
		return "", false
	}
	resStr, ok := res.(string)
	return resStr, ok
}

func (c Map) Merge(other Map) {
	for k, v := range other {
		if k == "header" {
			hdrs, ok := c[k]
			if ok {
				c[k] = append(hdrs.([]string), v.([]string)...)
				continue
			}
			c[k] = v.([]string)
			continue
		}
		c[k] = v
	}
}

func (c Map) String() string {
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
