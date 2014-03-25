package mustache

import (
	"bytes"
	"encoding/json"
	"github.com/acsellers/assert"
	"html/template"
	"testing"
)

/*
Inverted Section tags and End Section tags are used in combination to wrap a
section of the template.

These tags' content MUST be a non-whitespace character sequence NOT
containing the current closing delimiter; each Inverted Section tag MUST be
followed by an End Section tag with the same content assert.Within the same
section.

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

This section MUST NOT be rendered unless the data list is empty.

Inverted Section and End Section tags SHOULD be treated as standalone when
appropriate.

*/

func TestINVERTED0(t *testing.T) {
	// Falsey

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `"{{^boolean}}This should be rendered.{{/boolean}}"`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"boolean":false}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`"This should be rendered."`, b.String())
	})
}

func TestINVERTED1(t *testing.T) {
	// Truthy

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `"{{^boolean}}This should not be rendered.{{/boolean}}"`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"boolean":true}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`""`, b.String())
	})
}

func TestINVERTED2(t *testing.T) {
	// Context

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `"{{^context}}Hi {{name}}.{{/context}}"`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"context":{"name":"Joe"}}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`""`, b.String())
	})
}

func TestINVERTED3(t *testing.T) {
	// List

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `"{{^list}}{{n}}{{/list}}"`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"list":[{"n":1},{"n":2},{"n":3}]}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`""`, b.String())
	})
}

func TestINVERTED4(t *testing.T) {
	// Empty List

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `"{{^list}}Yay lists!{{/list}}"`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"list":[]}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`"Yay lists!"`, b.String())
	})
}

func TestINVERTED5(t *testing.T) {
	// Doubled

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `{{^bool}}
* first
{{/bool}}
* {{two}}
{{^bool}}
* third
{{/bool}}
`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"bool":false,"two":"second"}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`* first
* second
* third
`, b.String())
	})
}

func TestINVERTED6(t *testing.T) {
	// Nested (Falsey)

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `| A {{^bool}}B {{^bool}}C{{/bool}} D{{/bool}} E |`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"bool":false}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`| A B C D E |`, b.String())
	})
}

func TestINVERTED7(t *testing.T) {
	// Nested (Truthy)

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `| A {{^bool}}B {{^bool}}C{{/bool}} D{{/bool}} E |`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"bool":true}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`| A  E |`, b.String())
	})
}

func TestINVERTED8(t *testing.T) {
	// Context Misses

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `[{{^missing}}Cannot find key 'missing'!{{/missing}}]`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`[Cannot find key 'missing'!]`, b.String())
	})
}

func TestINVERTED9(t *testing.T) {
	// Dotted Names - Truthy

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `"{{^a.b.c}}Not Here{{/a.b.c}}" == ""`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"a":{"b":{"c":true}}}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`"" == ""`, b.String())
	})
}

func TestINVERTED10(t *testing.T) {
	// Dotted Names - Falsey

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `"{{^a.b.c}}Not Here{{/a.b.c}}" == "Not Here"`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"a":{"b":{"c":false}}}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`"Not Here" == "Not Here"`, b.String())
	})
}

func TestINVERTED11(t *testing.T) {
	// Dotted Names - Broken Chains

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `"{{^a.b.c}}Not Here{{/a.b.c}}" == "Not Here"`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"a":{}}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`"Not Here" == "Not Here"`, b.String())
	})
}

func TestINVERTED12(t *testing.T) {
	// Surrounding Whitespace

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", ` | {{^boolean}}	|	{{/boolean}} | 
`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"boolean":false}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(` | 	|	 | 
`, b.String())
	})
}

func TestINVERTED13(t *testing.T) {
	// Internal Whitespace

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", ` | {{^boolean}} {{! Important Whitespace }}
 {{/boolean}} | 
`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"boolean":false}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(` |  
  | 
`, b.String())
	})
}

func TestINVERTED14(t *testing.T) {
	// Indented Inline Sections

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", ` {{^boolean}}NO{{/boolean}}
 {{^boolean}}WAY{{/boolean}}
`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"boolean":false}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(` NO
 WAY
`, b.String())
	})
}

func TestINVERTED15(t *testing.T) {
	// Standalone Lines

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `| This Is
{{^boolean}}
|
{{/boolean}}
| A Line
`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"boolean":false}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`| This Is
|
| A Line
`, b.String())
	})
}

func TestINVERTED16(t *testing.T) {
	// Standalone Indented Lines

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `| This Is
  {{^boolean}}
|
  {{/boolean}}
| A Line
`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"boolean":false}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`| This Is
|
| A Line
`, b.String())
	})
}

func TestINVERTED17(t *testing.T) {
	// Standalone Line Endings

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `|
{{^boolean}}
{{/boolean}}
|`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"boolean":false}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`|
|`, b.String())
	})
}

func TestINVERTED18(t *testing.T) {
	// Standalone Without Previous Line

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `  {{^boolean}}
^{{/boolean}}
/`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"boolean":false}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`^
/`, b.String())
	})
}

func TestINVERTED19(t *testing.T) {
	// Standalone Without Newline

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `^{{^boolean}}
/
  {{/boolean}}`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"boolean":false}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`^
/
`, b.String())
	})
}

func TestINVERTED20(t *testing.T) {
	// Padding

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `|{{^ boolean }}={{/ boolean }}|`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"boolean":false}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`|=|`, b.String())
	})
}
