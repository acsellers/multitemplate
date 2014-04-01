package mustache

import (
	"bytes"
	"html/template"
	"testing"
	. "github.com/acsellers/assert"
)

/*
Lambdas are a special-cased data type for use in interpolations and
sections.

When used as the data value for an Interpolation tag, the lambda MUST be
treatable as an arity 0 function, and invoked as such.  The returned value
MUST be rendered against the default delimiters, then interpolated in place
of the lambda.

When used as the data value for a Section tag, the lambda MUST be treatable
as an arity 1 function, and invoked as such (passing a String containing the
unprocessed section contents).  The returned value MUST be rendered against
the current delimiters, then interpolated in place of the section.

*/

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

/*
// I'd like to interject and say this feature is about as WTF as it gets
func TestLAMBDAS1(t *testing.T) {
	// Interpolation - Expansion
  // Seriously WTF is this?

	Within(t, func(test *Test) {
		extraFuncs := template.FuncMap{
			"lambda": func() string {
				return "planet"
			},
		}
		t := template.New("test").Funcs(testFuncs).Funcs(extraFuncs)
		trees, err := Parse("test.mustache", `Hello, {{lambda}}!`, template.FuncMap{})
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}
		t.Funcs(template.FuncMap{})

		data := map[string]interface{}{"planet": "world"}
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`Hello, world!`, b.String())
	})
}

func TestLAMBDAS2(t *testing.T) {
	// Interpolation - Alternate Delimiters

	Within(t, func(test *Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `{{= | | =}}
Hello, (|&lambda|)!`, template.FuncMap{})
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}
		t.Funcs(template.FuncMap{})

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"lambda":{"__tag__":"code","clojure":"(fn [] \"|planet| =\u003e {{planet}}\")","js":"function() { return \"|planet| =\u003e {{planet}}\" }","perl":"sub { \"|planet| =\u003e {{planet}}\" }","php":"return \"|planet| =\u003e {{planet}}\";","python":"lambda: \"|planet| =\u003e {{planet}}\"","ruby":"proc { \"|planet| =\u003e {{planet}}\" }"},"planet":"world"}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`Hello, (|planet| => world)!`, b.String())
	})
}

func TestLAMBDAS3(t *testing.T) {
	// Interpolation - Multiple Calls

	Within(t, func(test *Test) {
		var index int
		extraFuncs := template.FuncMap{
			"lambda": func() string {
				index++
				return fmt.Sprint(index)
			},
		}
		t := template.New("test").Funcs(testFuncs).Funcs(extraFuncs)
		trees, err := Parse("test.mustache", `{{lambda}} == {{{lambda}}} == {{lambda}}`, extraFuncs)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}
		t.Funcs(template.FuncMap{})

		data := make(map[string]interface{})
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`1 == 2 == 3`, b.String())
	})
}

func TestLAMBDAS4(t *testing.T) {
	// Escaping

	Within(t, func(test *Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `<{{lambda}}{{{lambda}}}`, template.FuncMap{})
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}
		t.Funcs(template.FuncMap{})

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"lambda":{"__tag__":"code","clojure":"(fn [] \"\u003e\")","js":"function() { return \"\u003e\" }","perl":"sub { \"\u003e\" }","php":"return \"\u003e\";","python":"lambda: \"\u003e\"","ruby":"proc { \"\u003e\" }"}}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`<&gt;>`, b.String())
	})
}

*/

func TestLAMBDAS5(t *testing.T) {
	// Section

	Within(t, func(test *Test) {
		extraFuncs := template.FuncMap{
			"lambda": func() map[string]interface{} {
				return map[string]interface{}{"x": "yes"}
			},
		}
		t := template.New("test").Funcs(testFuncs).Funcs(extraFuncs)
		trees, err := Parse("test.mustache", `<{{#lambda}}{{x}}{{/lambda}}>`, template.FuncMap{})
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}
		t.Funcs(template.FuncMap{})

		data := make(map[string]interface{})
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`<yes>`, b.String())
	})
}

/*
func TestLAMBDAS6(t *testing.T) {
	// Section - Expansion

	Within(t, func(test *Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `<{{#lambda}}-{{/lambda}}>`, template.FuncMap{})
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}
		t.Funcs(template.FuncMap{})

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"lambda":{"__tag__":"code","clojure":"(fn [text] (str text \"{{planet}}\" text))","js":"function(txt) { return txt + \"{{planet}}\" + txt }","perl":"sub { $_[0] . \"{{planet}}\" . $_[0] }","php":"return $text . \"{{planet}}\" . $text;","python":"lambda text: \"%s{{planet}}%s\" % (text, text)","ruby":"proc { |text| \"#{text}{{planet}}#{text}\" }"},"planet":"Earth"}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`<-Earth->`, b.String())
	})
}

func TestLAMBDAS7(t *testing.T) {
	// Section - Alternate Delimiters

	Within(t, func(test *Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `{{= | | =}}<|#lambda|-|/lambda|>`, template.FuncMap{})
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}
		t.Funcs(template.FuncMap{})

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"lambda":{"__tag__":"code","clojure":"(fn [text] (str text \"{{planet}} =\u003e |planet|\" text))","js":"function(txt) { return txt + \"{{planet}} =\u003e |planet|\" + txt }","perl":"sub { $_[0] . \"{{planet}} =\u003e |planet|\" . $_[0] }","php":"return $text . \"{{planet}} =\u003e |planet|\" . $text;","python":"lambda text: \"%s{{planet}} =\u003e |planet|%s\" % (text, text)","ruby":"proc { |text| \"#{text}{{planet}} =\u003e |planet|#{text}\" }"},"planet":"Earth"}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`<-{{planet}} => Earth->`, b.String())
	})
}

func TestLAMBDAS8(t *testing.T) {
	// Section - Multiple Calls

	Within(t, func(test *Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `{{#lambda}}FILE{{/lambda}} != {{#lambda}}LINE{{/lambda}}`, template.FuncMap{})
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}
		t.Funcs(template.FuncMap{})

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"lambda":{"__tag__":"code","clojure":"(fn [text] (str \"__\" text \"__\"))","js":"function(txt) { return \"__\" + txt + \"__\" }","perl":"sub { \"__\" . $_[0] . \"__\" }","php":"return \"__\" . $text . \"__\";","python":"lambda text: \"__%s__\" % (text)","ruby":"proc { |text| \"__#{text}__\" }"}}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`__FILE__ != __LINE__`, b.String())
	})
}

func TestLAMBDAS9(t *testing.T) {
	// Inverted Section

	Within(t, func(test *Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `<{{^lambda}}{{static}}{{/lambda}}>`, template.FuncMap{})
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}
		t.Funcs(template.FuncMap{})

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"lambda":{"__tag__":"code","clojure":"(fn [text] false)","js":"function(txt) { return false }","perl":"sub { 0 }","php":"return false;","python":"lambda text: 0","ruby":"proc { |text| false }"},"static":"static"}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`<>`, b.String())
	})
}
*/
