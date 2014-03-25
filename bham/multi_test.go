package bham

import (
	"bytes"
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
