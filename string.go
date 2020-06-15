package gotemplate

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/kennygrant/sanitize"
)

// slugify
func urlPath(title string) string {
	reg, err := regexp.Compile("[^A-Za-z0-9]+")
	if err != nil {
		return ""
	}
	prettyurl := reg.ReplaceAllString(title, "-")
	prettyurl = strings.ToLower(strings.Trim(prettyurl, "-"))
	return prettyurl
}

func reReplaceAll(pattern, repl, text string) string {
	re := regexp.MustCompile(pattern)
	return re.ReplaceAllString(text, repl)
}

func empty(a interface{}) interface{} {
	k := reflect.ValueOf(a).Kind()
	if k == reflect.Int || k == reflect.Int16 || k == reflect.Int32 ||
		k == reflect.Int64 || k == reflect.Int8 || k == reflect.Bool ||
		k == reflect.Float32 || k == reflect.Float64 || k == reflect.Uint ||
		k == reflect.Uint16 || k == reflect.Uint32 || k == reflect.Uint64 ||
		k == reflect.Uint8 || k == reflect.Func {
		return a
	}
	v := reflect.ValueOf(a)
	if a == nil ||
		(k == reflect.Slice && v.Len() < 1) ||
		(k == reflect.Struct && v.NumField() < 1) ||
		(k == reflect.Map && v.Len() < 1) {
		return ""
	}
	if k == reflect.Struct {
		p := fmt.Sprintf("%#v", a)
		if p == "[]interface {}{}" || p == "map[string]interface {}{}" {
			return ""
		}
	}
	if k == reflect.Slice {
		for i := 0; i < v.Len(); i++ {
			if empty(v.Index(i)) != "" {
				return a
			}
		}
		return ""
	}
	if k == reflect.Map {
		for _, mk := range v.MapKeys() {
			if empty(v.MapIndex(mk)) != "" {
				return a
			}
		}
	}
	return a
}

func emptyStr(a interface{}) string {
	k := reflect.ValueOf(a).Kind()
	if k == reflect.Int || k == reflect.Int16 || k == reflect.Int32 ||
		k == reflect.Int64 || k == reflect.Int8 || k == reflect.Bool ||
		k == reflect.Float32 || k == reflect.Float64 || k == reflect.Uint ||
		k == reflect.Uint16 || k == reflect.Uint32 || k == reflect.Uint64 ||
		k == reflect.Uint8 || k == reflect.Func {
		return fmt.Sprintf("%v", a)
	}
	if a == nil || reflect.ValueOf(a).Len() < 1 {
		return ""
	}
	return fmt.Sprintf("%v", a)
}

func sanitise(str string) string {
	return sanitize.Name(strings.Replace(str, "/", " ", -1))
}

func limit(data interface{}, length int) interface{} {
	switch reflect.ValueOf(data).Kind() {
	case reflect.String:
		if len(data.(string)) > length {
			return data.(string)[:length]
		}
		return data

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf(fmt.Sprintf("%%-%dd", length), data)
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf(fmt.Sprintf("%%-%d.4f", length), data)
	}
	return data
}

func fixlen(length int, data interface{}) interface{} {
	switch reflect.ValueOf(data).Kind() {
	case reflect.String:
		return fmt.Sprintf(fmt.Sprintf("%%-%d.%ds", length, length), data)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf(fmt.Sprintf("%%-%d.%dd", length, length), data)
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf(fmt.Sprintf("%%-%d.4f", length), data)
	}
	return strings.Repeat(" ", length)
}

func fixlenright(length int, data interface{}) interface{} {
	switch reflect.ValueOf(data).Kind() {
	case reflect.String:
		return fmt.Sprintf(fmt.Sprintf("%%%d.%ds", length, length), data)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf(fmt.Sprintf("%%%d.%dd", length, length), data)
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf(fmt.Sprintf("%%%d.4f", length), data)
	}
	return strings.Repeat(" ", length)
}

func concat(ss ...string) string {
	return strings.Join(ss, "")
}

func toint(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func tofloat(s string) (float64, error) {
	if s == "" {
		s = "0.0"
	}
	return strconv.ParseFloat(s, 64)
}

func conditional(s1, s2 string) string {
	if s1 != "" {
		return s1
	}
	return s2
}

func notconditional(s1, s2 string) string {
	if s1 == "" {
		return s1
	}
	return s2
}
