package config

import "fmt"

type Map map[string]interface{}

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
