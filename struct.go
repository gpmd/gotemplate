package gotemplate

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/shoobyban/mxj"
	"github.com/spf13/cast"
)

func jsonEncode(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	return string(b), err
}

func xmlEncode(v interface{}) (string, error) {
	mv := mxj.Map(v.(map[string]interface{}))
	mxj.XMLEscapeChars(true)
	b, err := mv.Xml()
	return string(b), err
}

func xmlArray(v interface{}, roottag, itemtag string) (string, error) {
	mv := mxj.Map(map[string]interface{}{itemtag: v.([]interface{})})
	mxj.XMLEscapeChars(true)
	b, err := mv.XmlIndent("", "  ", roottag)
	return "<?xml version=\"1.0\"?>\n" + string(b), err
}

func decode(s, format string) (interface{}, error) {
	p := NewParser()
	res, err := p.ParseStruct(strings.NewReader(s), format)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse json '%s': %v", s, err)
	}
	return res, nil
}

func jsonDecode(s string) (interface{}, error) {
	return decode(s, "json")
}

func xmlDecode(s string) (interface{}, error) {
	return decode(s, "xml")
}

// jsonEscape escapes a variable (mostly string) for using inside a JSON as string
func jsonEscape(i interface{}) string {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	s := string(b)
	return s[1 : len(s)-1]
}

func asJSON(s interface{}) string {
	jsonBytes, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(jsonBytes)
}

func pathValue(keys []string, s interface{}, f string) (v interface{}) {
	var key string
	var nextkeys []string
	if len(keys) == 0 {
		if f == "" {
			return s
		}
		key = ""
		nextkeys = keys
	} else {
		key = keys[0]
		nextkeys = keys[1:]
	}
	filter := ""
	var (
		i  int64
		ok bool
	)
	var err error

	if key != "" && key[:1] == "[" && key[len(key)-1:] == "]" {
		key, filter = "", key[1:len(key)-1]
	}

	switch s.(type) {
	case map[string]interface{}:
		if key == "" {
			m := map[string]interface{}{}
			found := true
			if f != "" {
				found = false
				fparts := strings.Split(f, "=")
				for k, item := range s.(map[string]interface{}) {
					if k == fparts[0] && item == fparts[1] {
						found = true
					}
				}
			}
			if found {
				for k, item := range s.(map[string]interface{}) {
					m[k] = pathValue(nextkeys, item, filter)
				}
			}
			if len(m) > 0 {
				v = m
			}
		} else if v, ok = s.(map[string]interface{})[key]; !ok {
			err = fmt.Errorf("Key not present. [Key:%s]", key)
		}
	case []interface{}:
		array := s.([]interface{})
		a := []interface{}{}
		if f != "" {
			return a
		}
		if key == "" {
			for _, item := range array {
				pv := pathValue(nextkeys, item, filter)
				if pv != nil {
					a = append(a, pv)
				}
			}
			if len(a) == 1 {
				v = a[0]
			} else if len(a) > 0 {
				v = a
			}
		} else if i, err = strconv.ParseInt(key, 10, 64); err == nil {
			if int(i) < len(array) {
				v = array[i]
			} else {
				err = fmt.Errorf("Index out of bounds. [Index:%d] [Array:%v]", i, array)
			}
		}
	}
	return pathValue(nextkeys, v, "")
}

func explode(s, sep string) []string {
	return strings.Split(s, sep)
}

func inArray(needle interface{}, haystack interface{}) bool {
	switch reflect.TypeOf(haystack).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(haystack)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(needle, s.Index(i).Interface()) == true {
				return true
			}
		}
	}

	return false
}

func isSet(a interface{}, key interface{}) bool {
	av := reflect.ValueOf(a)
	kv := reflect.ValueOf(key)

	switch av.Kind() {
	case reflect.Array, reflect.Chan, reflect.Slice:
		k, err := cast.ToIntE(key)
		if err != nil {
			return false
		}
		if av.Len() > k {
			return true
		}
	case reflect.Map:
		if kv.Type() == av.Type().Key() {
			return av.MapIndex(kv).IsValid()
		}
	default:
		return false
	}

	return false
}

func mapto(item, mapvals, separators string) string {
	maps := strings.Split(mapvals, separators[:1])
	mapping := map[string]string{}
	for _, v := range maps {
		vv := strings.Split(v, separators[1:])
		if len(vv) < 2 {
			return ""
		}
		mapping[vv[0]] = vv[1]
	}
	if ret, ok := mapping[item]; ok {
		return ret
	}
	if ret, ok := mapping["*"]; ok {
		return ret
	}
	return item
}

func createMap() map[string]interface{} {
	return map[string]interface{}{}
}

func setItem(m map[string]interface{}, a string, b interface{}) map[string]interface{} {
	m[a] = b
	return m
}

func mkSlice(args ...interface{}) []interface{} {
	return args
}

func decimalFmt(format string, number interface{}) string {
	num := emptyStr(number)
	f, _ := strconv.ParseFloat(num, 64)
	i := strings.Split(format, ",")
	s := fmt.Sprintf(fmt.Sprintf("%%%s.%sf", i[0], i[1]), f)
	if i[1] != "0" {
		s = strings.TrimRight(s, "0")
		if s[len(s)-1:] == "." {
			s += "0"
		}
	}
	return s
}

func last(x int, a interface{}) bool {
	return x == reflect.ValueOf(a).Len()-1
}

func filterPath(s interface{}, p string) interface{} {
	return pathValue(strings.Split(p, "."), s, "")
}

func toAbs(float float64) float64 {
	if float < 0 {
		float = float * -1
	}
	return float
}
