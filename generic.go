package gotemplate

import (
	"crypto/md5"
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/cast"
)

func md5hash(data interface{}) string {
	s := fmt.Sprintf("%#v", data)
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func replace(input, from, to string) string {
	return strings.Replace(input, from, to, -1)
}

func urlencode(s string) string {
	return url.QueryEscape(s)
}

func urldecode(s string) (string, error) {
	r, err := url.QueryUnescape(s)
	if err != nil {
		return "", err
	}
	return r, nil
}

// From Hugo with love see https://gohugo.io/functions/seq/
// seq LAST
// seq FIRST LAST
// seq FIRST INCREMENT LAST
// It's named and used in the model of GNU's seq.
// Examples:
// anything wrong: empty list(!)
// 3: 1, 2, 3
// 1 2 4: 1, 3
// -3: -1, -2, -3
// 1 4: 1, 2, 3, 4
// 1 -2: 1, 0, -1, -2
func seq(args ...interface{}) []int {
	if len(args) < 1 || len(args) > 3 {
		return []int{}
	}

	intArgs := cast.ToIntSlice(args)
	if len(intArgs) < 1 || len(intArgs) > 3 {
		return []int{}
	}

	var inc = 1
	var last int
	var first = intArgs[0]

	if len(intArgs) == 1 {
		last = first
		if last == 0 {
			return []int{}
		} else if last > 0 {
			first = 1
		} else {
			first = -1
			inc = -1
		}
	} else if len(intArgs) == 2 {
		last = intArgs[1]
		if last < first {
			inc = -1
		}
	} else {
		inc = intArgs[1]
		last = intArgs[2]
		if inc == 0 {
			return []int{}
		}
		if first < last && inc < 0 {
			return []int{}
		}
		if first > last && inc > 0 {
			return []int{}
		}
	}

	// sanity check
	if last < -100000 {
		return []int{}
	}
	size := ((last - first) / inc) + 1

	// sanity check
	if size <= 0 || size > 2000 {
		return []int{}
	}

	seq := make([]int, size)
	val := first
	for i := 0; ; i++ {
		seq[i] = val
		val += inc
		if (inc < 0 && val < last) || (inc > 0 && val > last) {
			break
		}
	}

	return seq
}

func unique(e []interface{}) []interface{} {
	r := []interface{}{}

	for _, s := range e {
		if !inArray(s, r[:]) {
			r = append(r, s)
		}
	}
	return r
}

func item(s, sep string, num int) string {
	i := strings.Split(s, sep)
	if len(i) <= num {
		return ""
	}
	return i[num]
}

func escape(str string) string {
	return strings.Replace(str, "\"", "\\\"", -1)
}
