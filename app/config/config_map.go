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

		if len(parts) > 2 {
			println("Warning: bad config field:", field)
			continue
		}

		key, value := parts[0], parts[1]

		if key == "header" || key == "path" {
			prev, ok := res[key]
			if ok {
				res[key] = append(prev.([]string), value)
			} else {
				res[key] = []string{value}
			}
			continue
		}

		res[key] = value
	}

	return res
}

func (c Map) GetString(key string) string {
	res, ok := c[key]
	if !ok {
		return ""
	}
	resStr, _ := res.(string)
	return resStr
}

func (c Map) GetStringSlice(key string) []string {
	res, ok := c[key]
	if !ok {
		return nil
	}
	resSlice, _ := res.([]string)
	return resSlice
}

func (c Map) Merge(other Map) {
	for key, otherVal := range other {
		val, ok := c[key]
		if !ok {
			c[key] = otherVal
			continue
		}

		switch vv := val.(type) {
		case string: // Just override:
			c[key] = otherVal
		case bool: // Toggle:
			c[key] = vv != otherVal.(bool)
		case []string: // Append:
			c[key] = append(vv, otherVal.([]string)...)
		}
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
