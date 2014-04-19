package terse

import (
	"bytes"
	"testing"

	"github.com/acsellers/multitemplate"
)

func TestParse(t *testing.T) {
	for _, test := range parseTests {
		tmpl := multitemplate.New("terse")
		var e error
		if len(test.Sources) == 0 {
			tmpl, e = tmpl.Parse("parse", test.Content, "terse")
			if e != nil {
				t.Logf("In test %s\n", test.Name)
				t.Error("Parse Error:", e)
				continue
			}
		} else {
			for tn, tc := range test.Sources {
				tmpl, e = tmpl.Parse(tn, tc, "terse")
				if e != nil {
					t.Logf("In test %s\n", test.Name)
					t.Error("Parse Error:", e)
					continue
				}
			}
		}
		b := &bytes.Buffer{}
		if test.Template == "" {
			e = tmpl.Execute(b, test.Data)
		} else {
			e = tmpl.ExecuteTemplate(b, test.Template, test.Data)
		}
		if e != nil {
			t.Logf("In test %s\n", test.Name)
			t.Error("Execute Error:", e)
		}
		if b.String() != test.Expected {
			t.Logf("In test %s\n", test.Name)
			t.Errorf("Result Error, Expected:`%s`\nReceived:`%s`", test.Expected, b.String())
		}
	}
}

var parseTests = []parseTest{
	parseTest{
		Name: "Blank Template",
	},
}

type parseTest struct {
	Name     string
	Content  string
	Sources  map[string]string
	Expected string
	Template string
	Data     interface{}
}
