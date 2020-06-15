package gotemplate

import (
	"testing"
	"time"
)

type testDateStruct struct {
	Template string
	Values   map[string]interface{}
	Result   string
}

func TestDate(t *testing.T) {
	tests := map[string]testDateStruct{
		"ukdate": {
			Template: `Date: '{{ "2006-01-02 15:04:05" | date "ukshort" }}'`,
			Result:   "Date: '02/01/06'",
		},
		"timeformatminus": {
			Template: `{{timeformatminus "02/01/06 15:04:05" 5 }}`,
			Result:   time.Now().Add(time.Second * -5).Format("02/01/06 15:04:05"),
		},
		"timeformat": {
			Template: `{{timeformat "020106"}}`,
			Result:   time.Now().Format("020106"),
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
