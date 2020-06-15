package gotemplate

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"github.com/shoobyban/mxj"
)

type testParserStruct struct {
	Input  string
	Format string
	Reg    map[string]ParserFunc
	Result interface{}
}

func TestParseStruct(t *testing.T) {
	tests := map[string]testParserStruct{
		"xml": {
			Input:  `<?xml version="1.0"?><a><b>B</b><c>C</c></a>`,
			Format: "xml",
			Result: mxj.Map{"a": map[string]interface{}{"c": "C", "b": "B"}},
		},
		"json": {
			Input:  `{"a":["b","c"]}`,
			Format: "json",
			Result: mxj.Map{"a": []interface{}{"b", "c"}},
		},
		"csv": {
			Input:  "A,B\nC,D\n",
			Format: "csv",
			Result: []map[string]string{{"A": "C", "B": "D"}},
		},
		"underscore": {
			Input:  "A_B\nC_D\n",
			Format: "_",
			Result: [][]string{{"A", "B"}, {"C", "D"}},
			Reg: map[string]ParserFunc{
				"_": func(content io.Reader) (interface{}, error) {
					r := csv.NewReader(content)
					r.Comma = '_'
					r.Comment = '#'
					return r.ReadAll()
				},
			},
		},
		"edi": {
			Input:  "ordn_1\nsmtg_2\norln_a_3\norln_b_4\norln_c_5\n",
			Format: "bfk",
			Result: map[string]interface{}{
				"ordn": "1",
				"smtg": "2",
				"orln": [][]string{
					{"a", "3"},
					{"b", "4"},
					{"c", "5"}},
			},
			Reg: map[string]ParserFunc{
				"bfk": func(content io.Reader) (interface{}, error) {
					ret := map[string]interface{}{}
					b, err := ioutil.ReadAll(content)
					if err != nil {
						return nil, err
					}
					lines := strings.Split(string(b), "\n")
					for _, line := range lines {
						items := strings.Split(strings.Trim(line, "\r"), "_")
						if items[0] == "" {
							continue
						}
						if v, ok := ret[items[0]]; ok {
							if reflect.ValueOf(v).Kind() == reflect.Slice {
								if reflect.ValueOf(v).Index(0).Kind() == reflect.Slice {
									ret[items[0]] = append(ret[items[0]].([][]string), items[1:])
								} else {
									ret[items[0]] = append([][]string{}, ret[items[0]].([]string), items[1:])
								}
							} else if reflect.ValueOf(v).Kind() == reflect.String {
								ret[items[0]] = append([]interface{}{}, ret[items[0]], items[1:])
							} else if v == nil {
								ret[items[0]] = append([]interface{}{}, ret[items[0]], items[1:])
							} else {
								return nil, fmt.Errorf("Unhandled %#v", v)
							}
						} else {
							if len(items) > 2 {
								ret[items[0]] = items[1:]
							} else if len(items) == 2 {
								ret[items[0]] = items[1]
							} else {
								ret[items[0]] = nil
							}
						}
					}
					return ret, nil
				},
			},
		},
		"_items": {
			Input:  "##fn_2018042711432473\r\ntype_order_ack\r\nordn_20023\r\norln_115_73_1\r\n$$$$\r\n",
			Format: "bfk",
			Result: map[string]interface{}{"ordn": "20023", "orln": [][]string{{"115", "73", "1"}}, "$$$$": interface{}(nil), "##fn": "2018042711432473", "type": []string{"order", "ack"}},
			Reg: map[string]ParserFunc{
				"bfk": func(content io.Reader) (interface{}, error) {
					ret := map[string]interface{}{}
					b, err := ioutil.ReadAll(content)
					if err != nil {
						return nil, err
					}
					lines := strings.Split(string(b), "\n")
					forceslice := []string{"orln"}
					for _, line := range lines {
						items := strings.Split(strings.Trim(line, "\r"), "_")
						if items[0] == "" {
							continue
						}
						if len(forceslice) > 0 {
							for _, a := range forceslice {
								if a == items[0] {
									if _, ok := ret[items[0]]; !ok {
										ret[items[0]] = [][]string{}
									}
								}
							}
						}
						if v, ok := ret[items[0]]; ok {
							if reflect.ValueOf(v).Kind() == reflect.Slice {
								if len(v.([][]string)) == 0 {
									ret[items[0]] = append(v.([][]string), items[1:])
								} else if reflect.ValueOf(v).Index(0).Kind() == reflect.Slice {
									ret[items[0]] = append(ret[items[0]].([][]string), items[1:])
								} else {
									ret[items[0]] = append([][]string{}, ret[items[0]].([]string), items[1:])
								}
							} else if reflect.ValueOf(v).Kind() == reflect.String {
								ret[items[0]] = append([]interface{}{}, ret[items[0]], items[1:])
							} else if v == nil {
								ret[items[0]] = append([]interface{}{}, ret[items[0]], items[1:])
							} else {
								return nil, fmt.Errorf("Unhandled %#v", v)
							}
						} else {
							if len(items) > 2 {
								ret[items[0]] = items[1:]
							} else if len(items) == 2 {
								ret[items[0]] = items[1]
							} else {
								ret[items[0]] = nil
							}
						}
					}
					return ret, nil
				},
			},
		},
	}
	for name, test := range tests {
		l := NewParser()
		if test.Reg != nil {
			for name, parser := range test.Reg {
				l.RegisterParser(name, parser)
			}
		}
		res, err := l.ParseStruct(strings.NewReader(test.Input), test.Format)
		if err != nil {
			panic(err)
		}
		if !reflect.DeepEqual(res, test.Result) {
			t.Errorf("%s: %#v != %#v", name, res, test.Result)
		}
	}
}
