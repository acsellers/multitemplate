package mustache

import (
	"bytes"
	"encoding/json"
	"github.com/acsellers/assert"
	"html/template"
	"testing"
)

/*
Partial tags are used to expand an external template into the current
template.

The tag's content MUST be a non-whitespace character sequence NOT containing
the current closing delimiter.

This tag's content names the partial to inject.  Set Delimiter tags MUST NOT
affect the parsing of a partial.  The partial MUST be rendered against the
context stack local to the tag.  If the named partial cannot be found, the
empty string SHOULD be used instead, as in interpolations.

Partial tags SHOULD be treated as standalone when appropriate.  If this tag
is used standalone, any whitespace preceding the tag should treated as
indentation, and prepended to each line of the partial before rendering.

*/

func TestPARTIALS0(t *testing.T) {
	// Basic Behavior

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `"{{>text}}"`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		trees, err = Parse("text.mustache", `from partial`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`"from partial"`, b.String())
	})
}

func TestPARTIALS2(t *testing.T) {
	// Context

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `"{{>partial}}"`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		trees, err = Parse("partial.mustache", `*{{text}}*`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"text":"content"}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`"*content*"`, b.String())
	})
}

func TestPARTIALS3(t *testing.T) {
	// Recursion

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `{{>node}}`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		trees, err = Parse("node.mustache", `{{content}}<{{#nodes}}{{>node}}{{/nodes}}>`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"content":"X","nodes":[{"content":"Y","nodes":[]}]}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`X&lt;Y&lt;>>`, b.String())
	})
}

func TestPARTIALS4(t *testing.T) {
	// Surrounding Whitespace

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `| {{>partial}} |`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		trees, err = Parse("partial.mustache", `	|	`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`| 	|	 |`, b.String())
	})
}

func TestPARTIALS5(t *testing.T) {
	// Inline Indentation

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `  {{data}}  {{> partial}}
`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		trees, err = Parse("partial.mustache", `>
>`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"data":"|"}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`  |  >
>
`, b.String())
	})
}

func TestPARTIALS6(t *testing.T) {
	// Standalone Line Endings

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `|
{{>partial}}
|`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		trees, err = Parse("partial.mustache", `>`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`|
>|`, b.String())
	})
}

func TestPARTIALS7(t *testing.T) {
	// Standalone Without Previous Line

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `  {{>partial}}
>`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		trees, err = Parse("partial.mustache", `>
>`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`  >
>>`, b.String())
	})
}

func TestPARTIALS8(t *testing.T) {
	// Standalone Without Newline

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `>
  {{>partial}}`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		trees, err = Parse("partial.mustache", `>
>`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`>
  >
>`, b.String())
	})
}

func TestPARTIALS9(t *testing.T) {
	// Standalone Indentation

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `\
 {{>partial}}
/
`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		trees, err = Parse("partial.mustache", `|
{{{content}}}
|
`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"content":"\u003c\n-\u003e"}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`\
 |
<
->
|
/
`, b.String())
	})
}

func TestPARTIALS10(t *testing.T) {
	// Padding Whitespace

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `|{{> partial }}|`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		trees, err = Parse("partial.mustache", `[]`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"boolean":true}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`|[]|`, b.String())
	})
}
