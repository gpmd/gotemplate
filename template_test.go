package gotemplate

import (
	"testing"
)

type testTemplateStruct struct {
	Template string
	Values   map[string]interface{}
	Result   string
}

func TestTemplate(t *testing.T) {
	tests := map[string]testTemplateStruct{
		"isset": {
			Template: `{{if isset .A "z"}}OK{{end}}`,
			Values:   map[string]interface{}{"A": map[string]interface{}{"z": 1, "p": "a & b"}},
			Result:   "OK",
		},
		"not_isset": {
			Template: `{{if isset .A "c"}}OK{{end}}`,
			Values:   map[string]interface{}{"A": map[string]interface{}{"z": 1, "p": "a & b"}},
			Result:   "",
		},
		"xml_array": {
			Template: `{{xml_array .A "products" "product"}}`,
			Values:   map[string]interface{}{"A": []interface{}{map[string]interface{}{"z": 1, "p": "a & b"}, map[string]interface{}{"z": 2, "p": "b"}}},
			Result:   "<?xml version=\"1.0\"?>\n<products>\n  <product>\n    <p>a &amp; b</p>\n    <z>1</z>\n  </product>\n  <product>\n    <p>b</p>\n    <z>2</z>\n  </product>\n</products>",
		},
		"json_escape": {
			Template: `{{json_escape .A}}`,
			Values: map[string]interface{}{"A": `dog "fish"
 cat`},
			Result: `dog \"fish\"\n cat`,
		},
		"md5": {
			Template: `{{md5 .A}}`,
			Values:   map[string]interface{}{"A": []interface{}{}, "B": 1},
			Result:   `456a37d61262ccf952ee9768cbe32d94`,
		},
		"urlencode": {
			Template: `{{urlencode "Some & % / - Query"}}`,
			Result:   `Some+%26+%25+%2F+-+Query`,
		},
		"urldecode": {
			Template: `{{urldecode "Some+%26+%25+%2F+-+Query"}}`,
			Result:   `Some & % / - Query`,
		},
		"url_path": {
			Template: `{{url_path "Some Nice - URL"}}`,
			Result:   "some-nice-url",
		},
		"intf": {
			Template: "{{if .A}}.{{end}}",
			Values:   map[string]interface{}{"A": []interface{}{}, "B": 1},
			Result:   "",
		},
		"seq0123": {
			Template: `{{range seq 0 3}}{{.}} {{ end }}`,
			Result:   `0 1 2 3 `,
		},
		"seq123": {
			Template: `{{range seq 3}}{{.}} {{ end }}`,
			Result:   `1 2 3 `,
		},
		"replace": {
			Template: `{{replace "aBcd" "B" "b"}}`,
			Result:   "abcd",
		},
		"map": {
			Template: `{{ $m := createMap }}{{ $m := setItem $m "a" "b" }}{{ $m := setItem $m "c" "d" }}{{ range $i,$item := $m }} {{ $i }}:{{ $item }}{{ end }}`,
			Result:   " a:b c:d",
		},
		"slice": {
			Template: `{{ $slice := mkSlice "a" "b" "c" }}{{ range $slice }}{{.}}{{ end }}`,
			Result:   "abc",
		},
		"unique": {
			Template: `{{ $slice := mkSlice "a" "a" "b" "b" "c" }}{{ range (unique $slice) }}{{.}}{{ end }}`,
			Result:   "abc",
		},
		"reReplaceAll": {
			Template: `{{reReplaceAll "\"" "\\\"" .A }}`,
			Values:   map[string]interface{}{"A": `ab"cd"ef`},
			Result:   `ab\"cd\"ef`,
		},
		"reReplaceAll2": {
			Template: `{{reReplaceAll "\"" "&quot;" .A }}`,
			Values:   map[string]interface{}{"A": `ab"cd"ef`},
			Result:   `ab&quot;cd&quot;ef`,
		},

		"explode": {
			Template: `{{explode "1|2|3" "|"}}`,
			Result:   "[1 2 3]",
		},
		"in_array": {
			Template: `{{in_array "1" (explode "1|2|3" "|")}}`,
			Result:   "true",
		},
		"mapto": {
			Template: `{{range $key, $item := .test1}} {{ mapto $item "1:OK|2:Not OK|3:Maybe" "|:" }}{{end}}`,
			Values: map[string]interface{}{
				"test1": []string{"1", "2", "3"},
			},
			Result: " OK Not OK Maybe",
		},
		"xmldecode": {
			Template: `{{(xml_decode .val).analysis_code_15}}`,
			Values: map[string]interface{}{
				"val": `<?xml version=\"1.0\"?><analysis_code_15>Carneval &quot;Cool&quot; Point</analysis_code_15>`,
			},
			Result: `Carneval "Cool" Point`,
		},
		"xmlencode": {
			Template: `{{xml_encode .}}`,
			Values: map[string]interface{}{
				"analysis_code_15": "Carneval \"Cool\" Point",
			},
			Result: `<analysis_code_15>Carneval &quot;Cool&quot; Point</analysis_code_15>`,
		},
		"jsondecode": {
			Template: `{{(json_decode .val).analysis_code_15}}`,
			Values: map[string]interface{}{
				"val": `{"analysis_code_15":"Carneval \"Cool\" Point"}`,
			},
			Result: `Carneval "Cool" Point`,
		},
		"jsonencode": {
			Template: `{{json_encode .}}`,
			Values: map[string]interface{}{
				"analysis_code_15": "Carneval \"Cool\" Point",
			},
			Result: `{"analysis_code_15":"Carneval \"Cool\" Point"}`,
		},
		"divdec": {
			Template: `{{$l := len .a}}Len: {{$l}}{{ $b := (div $l .b) }}{{ $a := (div $l .c) }} A+B={{ add $a $b | decimal "1,2" }}`,
			Values: map[string]interface{}{
				"a": []interface{}{1, 2, 3},
				"b": 183.33,
				"c": 149.61,
			},
			Result: `Len: 3 A+B=110.98`,
		},
		"decimal": {
			Template: `All: '{{ "0.812545" | decimal "6,6" }}, {{ "0.1" | decimal "0,1" }}, {{ decimal "0,6" .test1 }}, {{ decimal "0,6" .test2 }}, {{ decimal "0,0" .test2 }}'`,
			Values: map[string]interface{}{
				"test1": nil,
				"test2": 10010.2342342342,
			},
			Result: "All: '0.812545, 0.1, 0.0, 10010.234234, 10010'",
		},
		"dec": {
			Template: `Dec: {{ "3" | int | sub 1 }}`,
			Result:   "Dec: 2",
		},
		"fixlen": {
			Template: `Fix: '{{ "A" | fixlen 5 }}'`,
			Result:   "Fix: 'A    '",
		},
		"fixlen2": {
			Template: `Fix: '{{ "ABCDEFG" | fixlen 5 }}'`,
			Result:   "Fix: 'ABCDE'",
		},
		"item": {
			Template: `C: '{{ item "1234-22" "-" 0 }}' D: '{{ item "1234-22" "-" 1 }}'`,
			Result:   "C: '1234' D: '22'",
		},
		"mapto2": {
			Template: `{{mapto "a" "a:True|b:False" "|:"}}`,
			Result:   "True",
		},
		"mapto3": {
			Template: `{{mapto "asdf" "a:A|b:B|*:C" "|:"}}`,
			Result:   "C",
		},
		"int": {
			Template: `'{{int "0123"}}'`,
			Result:   `'123'`,
		},
		"limit": {
			Template: `{{ limit "1234567890" 3 }}|{{ limit 1234 3 }}|{{ limit "12" 3 }}`,
			Result:   `123|1234|12`,
		},
		"empty": {
			Template: `{{ "1234567890" | empty }}|{{ "" | empty }}|{{ 3 | empty }}|{{.test1|empty}}|{{ $c := .test2|empty }}{{ concat "x" $c "y" }}|{{ $d := .test3|empty }}{{ concat "x" $d "y" }}`,
			Values: map[string]interface{}{
				"test1": []string{"1", "2", "3"},
				"test2": map[string]interface{}{},
				"test3": []interface{}{[]interface{}{}},
			},
			Result: `1234567890||3|[1 2 3]|xy|xy`,
		},
		"filter": {
			Template: `{{ $c := filter .countries "data.[iso=GB]" }}{{ $c.name }}`,
			Values: map[string]interface{}{
				"countries": map[string]interface{}{
					"data": []interface{}{
						map[string]interface{}{
							"iso":  "GB",
							"name": "Great Britain",
						},
						map[string]interface{}{
							"iso":  "US",
							"name": "United States",
						},
					},
				},
			},
			Result: `Great Britain`,
		},
		"deepfilter": {
			Template: `{{ $c := filter .countries "data.[iso=GB]" }}{{ $r := filter $c "states.[name=Surrey]" }}{{ $r.id }}`,
			Values: map[string]interface{}{
				"countries": map[string]interface{}{
					"data": []interface{}{
						map[string]interface{}{
							"iso":  "GB",
							"name": "Great Britain",
							"states": []interface{}{
								map[string]interface{}{
									"id":   1,
									"name": "Surrey",
								},
							},
						},
						map[string]interface{}{
							"iso":  "US",
							"name": "United States",
							"states": []interface{}{
								map[string]interface{}{
									"id":   2,
									"name": "Alabama",
								},
							},
						},
					},
				},
			},
			Result: `1`,
		},
	}
	for name, test := range tests {
		res, err := Template(test.Template, test.Values)
		if err != nil {
			panic(err)
		}
		if res != test.Result {
			t.Errorf("%s: %#v != %#v", name, res, test.Result)
		}
	}
}
