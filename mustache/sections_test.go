package mustache

import (
	"bytes"
	"encoding/json"
	"html/template"
	"testing"

	"github.com/acsellers/assert"
)

/*
Section tags and End Section tags are used in combination to wrap a section
of the template for iteration

These tags' content MUST be a non-whitespace character sequence NOT
containing the current closing delimiter; each Section tag MUST be followed
by an End Section tag with the same content assert.Within the same section.

This tag's content names the data to replace the tag.  Name resolution is as
follows:
  1) Split the name on periods; the first part is the name to resolve, any
  remaining parts should be retained.
  2) Walk the context stack from top to bottom, finding the first context
  that is a) a hash containing the name as a key OR b) an object responding
  to a method with the given name.
  3) If the context is a hash, the data is the value associated with the
  name.
  4) If the context is an object and the method with the given name has an
  arity of 1, the method SHOULD be called with a String containing the
  unprocessed contents of the sections; the data is the value returned.
  5) Otherwise, the data is the value returned by calling the method with
  the given name.
  6) If any name parts were retained in step 1, each should be resolved
  against a context stack containing only the result from the former
  resolution.  If any part fails resolution, the result should be considered
  falsey, and should interpolate as the empty string.
If the data is not of a list type, it is coerced into a list as follows: if
the data is truthy (e.g. `!!data == true`), use a single-element list
containing the data, otherwise use an empty list.

For each element in the data list, the element MUST be pushed onto the
context stack, the section MUST be rendered, and the element MUST be popped
off the context stack.

Section and End Section tags SHOULD be treated as standalone when
appropriate.

*/

func TestSECTIONS0(t *testing.T) {
	// Truthy

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `"{{#boolean}}This should be rendered.{{/boolean}}"`, template.FuncMap{})
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"boolean":true}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`"This should be rendered."`, b.String())
	})
}

func TestSECTIONS1(t *testing.T) {
	// Falsey

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `"{{#boolean}}This should not be rendered.{{/boolean}}"`, template.FuncMap{})
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"boolean":false}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`""`, b.String())
	})
}

func TestSECTIONS2(t *testing.T) {
	// Context

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `"{{#context}}Hi {{name}}.{{/context}}"`, template.FuncMap{})
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"context":{"name":"Joe"}}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`"Hi Joe."`, b.String())
	})
}

func TestSECTIONS3(t *testing.T) {
	// Deeply Nested Contexts

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `{{#a}}
{{one}}
{{#b}}
{{one}}{{two}}{{one}}
{{#c}}
{{one}}{{two}}{{three}}{{two}}{{one}}
{{#d}}
{{one}}{{two}}{{three}}{{four}}{{three}}{{two}}{{one}}
{{#e}}
{{one}}{{two}}{{three}}{{four}}{{five}}{{four}}{{three}}{{two}}{{one}}
{{/e}}
{{one}}{{two}}{{three}}{{four}}{{three}}{{two}}{{one}}
{{/d}}
{{one}}{{two}}{{three}}{{two}}{{one}}
{{/c}}
{{one}}{{two}}{{one}}
{{/b}}
{{one}}
{{/a}}
`, template.FuncMap{})
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"a":{"one":1},"b":{"two":2},"c":{"three":3},"d":{"four":4},"e":{"five":5}}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`1
121
12321
1234321
123454321
1234321
12321
121
1
`, b.String())
	})
}

func TestSECTIONS4(t *testing.T) {
	// List

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `"{{#list}}{{item}}{{/list}}"`, template.FuncMap{})
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"list":[{"item":1},{"item":2},{"item":3}]}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`"123"`, b.String())
	})
}

func TestSECTIONS5(t *testing.T) {
	// Empty List

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `"{{#list}}Yay lists!{{/list}}"`, template.FuncMap{})
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"list":[]}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`""`, b.String())
	})
}

func TestSECTIONS6(t *testing.T) {
	// Doubled

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `{{#bool}}
* first
{{/bool}}
* {{two}}
{{#bool}}
* third
{{/bool}}
`, template.FuncMap{})
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"bool":true,"two":"second"}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`* first
* second
* third
`, b.String())
	})
}

func TestSECTIONS7(t *testing.T) {
	// Nested (Truthy)

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `| A {{#bool}}B {{#bool}}C{{/bool}} D{{/bool}} E |`, template.FuncMap{})
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"bool":true}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`| A B C D E |`, b.String())
	})
}

func TestSECTIONS8(t *testing.T) {
	// Nested (Falsey)

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `| A {{#bool}}B {{#bool}}C{{/bool}} D{{/bool}} E |`, template.FuncMap{})
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"bool":false}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`| A  E |`, b.String())
	})
}

func TestSECTIONS9(t *testing.T) {
	// Context Misses

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `[{{#missing}}Found key 'missing'!{{/missing}}]`, template.FuncMap{})
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`[]`, b.String())
	})
}

func TestSECTIONS10(t *testing.T) {
	// Implicit Iterator - String

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `"{{#list}}({{.}}){{/list}}"`, template.FuncMap{})
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"list":["a","b","c","d","e"]}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`"(a)(b)(c)(d)(e)"`, b.String())
	})
}

func TestSECTIONS11(t *testing.T) {
	// Implicit Iterator - Integer

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `"{{#list}}({{.}}){{/list}}"`, template.FuncMap{})
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"list":[1,2,3,4,5]}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`"(1)(2)(3)(4)(5)"`, b.String())
	})
}

func TestSECTIONS12(t *testing.T) {
	// Implicit Iterator - Decimal

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `"{{#list}}({{.}}){{/list}}"`, template.FuncMap{})
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"list":[1.1,2.2,3.3,4.4,5.5]}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`"(1.1)(2.2)(3.3)(4.4)(5.5)"`, b.String())
	})
}

func TestSECTIONS13(t *testing.T) {
	// Dotted Names - Truthy

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `"{{#a.b.c}}Here{{/a.b.c}}" == "Here"`, template.FuncMap{})
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"a":{"b":{"c":true}}}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`"Here" == "Here"`, b.String())
	})
}

func TestSECTIONS14(t *testing.T) {
	// Dotted Names - Falsey

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `"{{#a.b.c}}Here{{/a.b.c}}" == ""`, template.FuncMap{})
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"a":{"b":{"c":false}}}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`"" == ""`, b.String())
	})
}

func TestSECTIONS15(t *testing.T) {
	// Dotted Names - Broken Chains

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `"{{#a.b.c}}Here{{/a.b.c}}" == ""`, template.FuncMap{})
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"a":{}}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`"" == ""`, b.String())
	})
}

func TestSECTIONS16(t *testing.T) {
	// Surrounding Whitespace

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", ` | {{#boolean}}	|	{{/boolean}} | 
`, template.FuncMap{})
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"boolean":true}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(` | 	|	 | 
`, b.String())
	})
}

func TestSECTIONS17(t *testing.T) {
	// Internal Whitespace

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", ` | {{#boolean}} {{! Important Whitespace }}
 {{/boolean}} | 
`, template.FuncMap{})
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"boolean":true}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(` |  
  | 
`, b.String())
	})
}

func TestSECTIONS18(t *testing.T) {
	// Indented Inline Sections

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", ` {{#boolean}}YES{{/boolean}}
 {{#boolean}}GOOD{{/boolean}}
`, template.FuncMap{})
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"boolean":true}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(` YES
 GOOD
`, b.String())
	})
}

func TestSECTIONS19(t *testing.T) {
	// Standalone Lines

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `| This Is
{{#boolean}}
|
{{/boolean}}
| A Line
`, template.FuncMap{})
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"boolean":true}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`| This Is
|
| A Line
`, b.String())
	})
}

func TestSECTIONS20(t *testing.T) {
	// Indented Standalone Lines

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `| This Is
  {{#boolean}}
|
  {{/boolean}}
| A Line
`, template.FuncMap{})
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"boolean":true}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`| This Is
|
| A Line
`, b.String())
	})
}

func TestSECTIONS21(t *testing.T) {
	// Standalone Line Endings

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `|
{{#boolean}}
{{/boolean}}
|`, template.FuncMap{})
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"boolean":true}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`|
|`, b.String())
	})
}

func TestSECTIONS22(t *testing.T) {
	// Standalone Without Previous Line

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `  {{#boolean}}
#{{/boolean}}
/`, template.FuncMap{})
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"boolean":true}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`#
/`, b.String())
	})
}

func TestSECTIONS23(t *testing.T) {
	// Standalone Without Newline

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `#{{#boolean}}
/
  {{/boolean}}`, template.FuncMap{})
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"boolean":true}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`#
/
`, b.String())
	})
}

func TestSECTIONS24(t *testing.T) {
	// Padding

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `|{{# boolean }}={{/ boolean }}|`, template.FuncMap{})
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"boolean":true}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`|=|`, b.String())
	})
}
