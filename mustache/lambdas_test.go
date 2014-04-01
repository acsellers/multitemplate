package mustache

import (
	"bytes"
	"html/template"
	"testing"
	. "github.com/acsellers/assert"
)

// Lambdas are a special-cased data type for use in interpolations.

func TestLAMBDAS0(t *testing.T) {
	// Interpolation

	Within(t, func(test *Test) {
		extraFuncs := template.FuncMap{
			"lambda": func() string {
				return "world"
			},
		}
		t := template.New("test").Funcs(testFuncs).Funcs(extraFuncs)
		trees, err := Parse("test.mustache", `Hello, {{lambda}}!`, extraFuncs)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}
		data := make(map[string]interface{})
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`Hello, world!`, b.String())
	})
}
