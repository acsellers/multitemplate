package mustache

import (
	"bytes"
	"encoding/json"
	"html/template"
	"testing"

	"github.com/acsellers/assert"
)

/*
Comment tags represent content that should never appear in the resulting
output.

The tag's content may contain any substring (including newlines) EXCEPT the
closing delimiter.

Comment tags SHOULD be treated as standalone when appropriate.

*/

func TestCOMMENTS0(t *testing.T) {
	// Inline

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `12345{{! Comment Block! }}67890`, template.FuncMap{})
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`1234567890`, b.String())
	})
}

func TestCOMMENTS1(t *testing.T) {
	// Multiline

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `12345{{!
  This is a
  multi-line comment...
}}67890
`, template.FuncMap{})
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`1234567890
`, b.String())
	})
}

func TestCOMMENTS2(t *testing.T) {
	// Standalone

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `Begin.
{{! Comment Block! }}
End.
`, template.FuncMap{})
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

func TestCOMMENTS3(t *testing.T) {
	// Indented Standalone

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `Begin.
  {{! Indented Comment Block! }}
End.
`, template.FuncMap{})
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

func TestCOMMENTS4(t *testing.T) {
	// Standalone Line Endings

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `|
{{! Standalone Comment }}
|`, template.FuncMap{})
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

func TestCOMMENTS5(t *testing.T) {
	// Standalone Without Previous Line

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `  {{! I'm Still Standalone }}
!`, template.FuncMap{})
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`!`, b.String())
	})
}

func TestCOMMENTS6(t *testing.T) {
	// Standalone Without Newline

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `!
  {{! I'm Still Standalone }}`, template.FuncMap{})
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`!
`, b.String())
	})
}

func TestCOMMENTS7(t *testing.T) {
	// Multiline Standalone

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `Begin.
{{!
Something's going on here...
}}
End.
`, template.FuncMap{})
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

func TestCOMMENTS8(t *testing.T) {
	// Indented Multiline Standalone

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `Begin.
  {{!
    Something's going on here...
  }}
End.
`, template.FuncMap{})
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

func TestCOMMENTS9(t *testing.T) {
	// Indented Inline

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `  12 {{! 34 }}
`, template.FuncMap{})
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`  12 
`, b.String())
	})
}

func TestCOMMENTS10(t *testing.T) {
	// Surrounding Whitespace

	assert.Within(t, func(test *assert.Test) {
		t := template.New("test").Funcs(testFuncs)
		trees, err := Parse("test.mustache", `12345 {{! Comment Block! }} 67890`, template.FuncMap{})
		test.IsNil(err)
		for name, tree := range trees {
			t, err = t.AddParseTree(name, tree)
			test.IsNil(err)
		}

		data := make(map[string]interface{})
		test.IsNil(json.Unmarshal([]byte(`{}`), &data))
		b := new(bytes.Buffer)
		test.IsNil(t.ExecuteTemplate(b, "test", data))
		test.AreEqual(`12345  67890`, b.String())
	})
}
