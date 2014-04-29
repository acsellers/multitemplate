package bham

import (
	"bytes"
	"html/template"
	"testing"

	. "github.com/acsellers/assert"
	"github.com/acsellers/multitemplate"
)

func TestTemplate(tst *testing.T) {
	Within(tst, func(test *Test) {
		for p, s := range simpleTest.Cases {
			t, e := multitemplate.New("simple").Parse("simple", s, p)
			test.IsNil(e)
			b := &bytes.Buffer{}
			t.Execute(b, nil)
			test.AreEqual(b.String(), simpleTest.Expected)
		}
	})
}

func TestVariable(t *testing.T) {
	tmpl, e := multitemplate.New("vars").Funcs(template.FuncMap{"wat": func(string, string) string {
		return "now"
	}}).Parse("vars", `= $f := "123"
= wat $f "nana"`, "bham")
	if e != nil {
		t.Log("Parse Error: ", e)
		t.Fail()
	}
	b := &bytes.Buffer{}
	tmpl.Execute(b, nil)
	if b.String() != "now" {
		t.Errorf("Expected: '123', found: '%s'", b.String())
	}
}

type templateTest struct {
	Expected string
	Cases    map[string]string
}

var simpleTest = templateTest{
	Expected: "<b>Test</b>",
	Cases: map[string]string{
		"stdlib": "<b>Test</b>",
	},
}
