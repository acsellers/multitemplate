package mustache

import (
	"bytes"
	"encoding/json"
	"github.com/acsellers/assert"
	"html/template"
	"testing"
)

/*
Set Delimiter tags are used to change the tag delimiters for all content
following the tag in the current compilation unit.

The tag's content MUST be any two non-whitespace sequences (separated by
whitespace) EXCEPT an equals sign ('=') followed by the current closing
delimiter.

Set Delimiter tags SHOULD be treated as standalone when appropriate.

*/

func TestDELIMITERS0(t *testing.T) {
	// Pair Behavior

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `{{=<% %>=}}(<%text%>)`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"text":"Hey!"}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`(Hey!)`, b.String())
	})
}

func TestDELIMITERS1(t *testing.T) {
	// Special Characters

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `({{=[ ]=}}[text])`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"text":"It worked!"}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`(It worked!)`, b.String())
	})
}

func TestDELIMITERS2(t *testing.T) {
	// Sections

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `[
{{#section}}
  {{data}}
  |data|
{{/section}}

{{= | | =}}
|#section|
  {{data}}
  |data|
|/section|
]
`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"data":"I got interpolated.","section":true}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`[
  I got interpolated.
  |data|

  {{data}}
  I got interpolated.
]
`, b.String())
	})
}

func TestDELIMITERS3(t *testing.T) {
	// Inverted Sections

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `[
{{^section}}
  {{data}}
  |data|
{{/section}}

{{= | | =}}
|^section|
  {{data}}
  |data|
|/section|
]
`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"data":"I got interpolated.","section":false}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`[
  I got interpolated.
  |data|

  {{data}}
  I got interpolated.
]
`, b.String())
	})
}

func TestDELIMITERS4(t *testing.T) {
	// Partial Inheritence

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `[ {{>include}} ]
{{= | | =}}
[ |>include| ]
`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		trees, err = Parse("include.mustache", `.{{value}}.`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"value":"yes"}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`[ .yes. ]
[ .yes. ]
`, b.String())
	})
}

func TestDELIMITERS5(t *testing.T) {
	// Post-Partial Behavior

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `[ {{>include}} ]
[ .{{value}}.  .|value|. ]
`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		trees, err = Parse("include.mustache", `.{{value}}. {{= | | =}} .|value|.`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{"value":"yes"}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`[ .yes.  .yes. ]
[ .yes.  .|value|. ]
`, b.String())
	})
}

func TestDELIMITERS6(t *testing.T) {
	// Surrounding Whitespace

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `| {{=@ @=}} |`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`|  |`, b.String())
	})
}

func TestDELIMITERS7(t *testing.T) {
	// Outlying Whitespace (Inline)

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", ` | {{=@ @=}}
`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(` | 
`, b.String())
	})
}

func TestDELIMITERS8(t *testing.T) {
	// Standalone Tag

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `Begin.
{{=@ @=}}
End.
`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`Begin.
End.
`, b.String())
	})
}

func TestDELIMITERS9(t *testing.T) {
	// Indented Standalone Tag

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `Begin.
  {{=@ @=}}
End.
`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`Begin.
End.
`, b.String())
	})
}

func TestDELIMITERS10(t *testing.T) {
	// Standalone Line Endings

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `|
{{= @ @ =}}
|`)
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
|`, b.String())
	})
}

func TestDELIMITERS11(t *testing.T) {
	// Standalone Without Previous Line

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `  {{=@ @=}}
=`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`=`, b.String())
	})
}

func TestDELIMITERS12(t *testing.T) {
	// Standalone Without Newline

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `=
  {{=@ @=}}`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`=
`, b.String())
	})
}

func TestDELIMITERS13(t *testing.T) {
	// Pair with Padding

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `|{{= @   @ =}}|`)
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`||`, b.String())
	})
}
