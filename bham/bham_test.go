package bham

import (
	"bytes"
	"github.com/acsellers/assert"
	"testing"
	"text/template"
)

func TestParse(t *testing.T) {
	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(map[string]interface{}{})
		tree, err := Parse("test.bham", "%html\n\t%head\n\t\t%title wat")
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		b := new(bytes.Buffer)
		t.Execute(b, nil)
		test.AreEqual("<html><head><title>  wat </title> </head></html>", b.String())
	})
}

func TestParse2(t *testing.T) {
	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(map[string]interface{}{})
		tree, err := Parse("test.bham", "%html\n\t%head\n\t\t%title\n\t\t\twat")
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		b := new(bytes.Buffer)
		t.Execute(b, nil)
		test.AreEqual("<html><head><title>wat </title></head></html>", b.String())
	})
}

func TestParseIf(t *testing.T) {
	assert.Within(t, func(test *assert.Test) {
		tmpl := `%html
  %head
    = if .ShowWat
      %title wat
  %body
    #content
      moo
    `
		t := template.New("test").Funcs(map[string]interface{}{})
		tree, err := Parse("test.bham", tmpl)
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		b := new(bytes.Buffer)
		t.Execute(b, map[string]interface{}{"ShowWat": true})
		test.AreEqual("<html><head><title>  wat </title> </head><body><div id=\"content\">moo </div></body></html>", b.String())

		b.Reset()
		t.Execute(b, map[string]interface{}{"ShowWat": false})
		test.AreEqual("<html><head></head><body><div id=\"content\">moo </div></body></html>", b.String())
	})
}

func TestParseIfElse(t *testing.T) {
	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(map[string]interface{}{})
		tree, err := Parse("test.bham", "%html\n\t%head\n\t\t= if .ShowWat\n\t\t\t%title wat\n\t\t= else\n\t\t\t%title taw")
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		b := new(bytes.Buffer)
		t.Execute(b, map[string]interface{}{"ShowWat": true})
		test.AreEqual("<html><head><title>  wat </title> </head></html>", b.String())

		b.Reset()
		t.Execute(b, map[string]interface{}{"ShowWat": false})
		test.AreEqual("<html><head><title>  taw </title> </head></html>", b.String())
	})
}

func TestParseRange(t *testing.T) {
	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(map[string]interface{}{})
		tree, err := Parse("test.bham", "%html\n\t%body\n\t\t= range .Wats\n\t\t\t%p wat")
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		b := new(bytes.Buffer)
		t.Execute(b, map[string]interface{}{"Wats": []int{1, 2}})
		test.AreEqual("<html><body><p>  wat </p> <p>  wat </p> </body></html>", b.String())

		b.Reset()
		t.Execute(b, map[string]interface{}{"Wats": []int{}})
		test.AreEqual("<html><body></body></html>", b.String())
	})
}

func TestParseRangeElse(t *testing.T) {
	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(map[string]interface{}{})
		tree, err := Parse("test.bham", "%html\n\t%body\n\t\t= range .Wats\n\t\t\t%p wat\n\t\t= else\n\t\t\t%p no wat")
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		b := new(bytes.Buffer)
		t.Execute(b, map[string]interface{}{"Wats": []int{1, 2}})
		test.AreEqual("<html><body><p>  wat </p> <p>  wat </p> </body></html>", b.String())

		b.Reset()
		t.Execute(b, map[string]interface{}{"Wats": []int{}})
		test.AreEqual("<html><body><p>  no wat </p> </body></html>", b.String())
	})
}

func TestTextPassthrough(t *testing.T) {
	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(map[string]interface{}{})
		tree, err := Parse("test.bham", `<!DOCTYPE html>
%html
  %body Test Line
    Test other line`)
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		b := new(bytes.Buffer)
		t.Execute(b, map[string]interface{}{"Wats": []int{1, 2}})
		test.AreEqual("<!DOCTYPE html> <html><body> Test Line Test other line </body></html>", b.String())
	})
}

func TestOutput(t *testing.T) {
	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(map[string]interface{}{})
		tree, err := Parse("test.bham", "<!DOCTYPE html>\n%html\n\t%body\n\t\t= .Name")
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		b := new(bytes.Buffer)
		t.Execute(b, map[string]interface{}{"Name": "Andrew"})
		test.AreEqual("<!DOCTYPE html> <html><body>Andrew</body></html>", b.String())
	})
}

func TestFindAttrs(t *testing.T) {
	assert.Within(t, func(test *assert.Test) {
		a, b := findAttrs("((((((")
		test.AreEqual(a, "")
		test.AreEqual(b, "((((((")

		a, b = findAttrs("()")
		test.AreEqual(a, "")
		test.AreEqual(b, "")

		a, b = findAttrs("(ng-app)hiip")
		test.AreEqual(a, "ng-app")
		test.AreEqual(b, "hiip")

	})
}

func TestAttribute(t *testing.T) {
	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(map[string]interface{}{})
		tree, err := Parse("test.bham", "<!DOCTYPE html>\n%html(ng-app)\n\t%body(ng-controller=\"PageController\")")
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		b := new(bytes.Buffer)
		t.Execute(b, map[string]interface{}{"Name": "Andrew"})
		test.AreEqual("<!DOCTYPE html> <html ng-app><body ng-controller=\"PageController\"></body></html>", b.String())
	})
}

func TestClass(t *testing.T) {
	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(map[string]interface{}{})
		tree, err := Parse("test.bham", "%html\n\t%body\n\t\t%div.see.me(class=\"soon\")")
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		b := new(bytes.Buffer)
		t.Execute(b, nil)
		test.AreEqual("<html><body><div class=\"see me soon\"></div></body></html>", b.String())
	})
}

func TestId(t *testing.T) {
	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(map[string]interface{}{})
		tree, err := Parse("test.bham", "%html\n\t%body\n\t\t%div#see(id=\"me\")")
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		b := new(bytes.Buffer)
		t.Execute(b, map[string]interface{}{"Name": "Andrew"})

		test.AreEqual("<html><body><div id=\"see_me\"></div></body></html>", b.String())
	})
}

func TestBareId(t *testing.T) {
	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(map[string]interface{}{})
		tree, err := Parse("test.bham", "%html\n\t%body\n\t\t#see(id=\"me\")")
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		b := new(bytes.Buffer)
		t.Execute(b, map[string]interface{}{"Name": "Andrew"})

		test.AreEqual("<html><body><div id=\"see_me\"></div></body></html>", b.String())
	})
}

func TestWith(t *testing.T) {
	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(map[string]interface{}{})
		tree, err := Parse("test.bham", "%html\n\t= with $name := \"Killer\"\n\t\t= $name")
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		b := new(bytes.Buffer)
		test.IsNil(t.Execute(b, nil))

		test.AreEqual("<html>Killer</html>", b.String())
	})
}

func FuncTestVar(t *testing.T) {
	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(map[string]interface{}{})
		tree, err := Parse("test.bham", "%html\n\t= $name := \"Killer\"\n\t\t= $name")
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		b := new(bytes.Buffer)
		t.Execute(b, nil)

		test.AreEqual("<html>Killer</html>", b.String())
	})
}

func TestEmbedded(t *testing.T) {
	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(map[string]interface{}{})
		tree, err := Parse("test.bham", "%html {{ .Name }}\n\t%head\n\t\t%title(class=\"{{ .Class }}\") wat")
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		b := new(bytes.Buffer)
		t.Execute(b, map[string]interface{}{"Name": "Andrew Sellers", "Class": "big"})
		test.AreEqual("<html> Andrew Sellers<head><title class=\"big\">  wat </title> </head></html>", b.String())
	})
}

func TestFilter(t *testing.T) {
	tmplContent := `%html
  :javascript
    var name = "Andrew";
`

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(map[string]interface{}{})
		tree, err := Parse("test.bham", tmplContent)
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		b := new(bytes.Buffer)
		t.Execute(b, map[string]interface{}{"Name": "Andrew", "Class": "big"})
		test.AreEqual("<html><script type=\"text/javascript\">var name = \"Andrew\";</script></html>", b.String())

	})
}

/*
func TestFilter2(t *testing.T) {
	tmplContent := `%html
  :javascript
    var name = "{{ .Name | js }}";
`

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(map[string]interface{}{})
		tree, err := Parse("test.bham", tmplContent)
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		b := new(bytes.Buffer)
		t.Execute(b, map[string]interface{}{"Name": "Andrew", "Class": "big"})
		test.AreEqual("<html><script type=\"text/javascript\"> var name = \"Andrew\";</script></html>", b.String())

	})
}
*/
