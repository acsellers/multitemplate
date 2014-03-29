package bham

import (
	"bytes"
	"html/template"
	"testing"

	"github.com/acsellers/assert"
)

const (
	large = `<!DOCTYPE html>

%html
  = if .flash.success
    %p.fSuccess {{.flash.success}}
  = print ""
  Created with the
  %a(href="http://github.com/robfig/revel")
    Revel Framework
`
	result = `<!DOCTYPE html> <html>Created with the <a href="http://github.com/robfig/revel">Revel Framework </a></html>`
)

func TestLarge(t *testing.T) {
	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(map[string]interface{}{})
		tree, err := Parse("test.bham", large, template.FuncMap{})
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test.bham"])
		test.IsNil(err)

		b := new(bytes.Buffer)
		test.IsNil(t.Execute(b, nil))
		test.AreEqual(result, b.String())
	})
}
