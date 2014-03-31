package multitemplate

import (
	"bytes"
	"testing"

	. "github.com/acsellers/assert"
)

func TestTemplate(tst *testing.T) {
	Within(tst, func(test *Test) {
		for p, s := range simpleTest.Cases {
			t, e := New("simple").Parse("simple", s, p)
			test.IsNil(e)
			b := &bytes.Buffer{}
			t.Execute(b, nil)
			test.AreEqual(b.String(), simpleTest.Expected)
		}
	})
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

func TestBlock(tst *testing.T) {
	Within(tst, func(test *Test) {
		code := `{{ define "view.html" }}
{{ extends "base.html" }}

{{ block "header" }}<styles>{{ end_block }}
{{ end }}
{{ define "base.html" }}<before>
{{ block "header"}}<links>{{ end_block }}
<after>{{ end }}

`
		t, e := New("view.html").Parse("view.html", code, "stdlib")
		test.IsNil(e)
		b := &bytes.Buffer{}
		test.NoError(t.ExecuteTemplate(b, "view.html", nil))
		test.AreEqual(b.String(), "\n<before>\n<styles>\n<after>")

	})
}
