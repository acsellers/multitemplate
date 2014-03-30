package multitemplate

import (
	"bytes"
	"html/template"
	"testing"

	. "github.com/acsellers/assert"
)

func TestYieldingBlocks(t *testing.T) {
	Within(t, func(test *Test) {
		t, e := New("main").Parse("main", `{{ yield "test_block" }}`, "default")
		test.NoError(e)
		c := NewContext(nil)
		c.Blocks["test_block"] = template.HTML("test block content")
		c.Main = "main"
		b := bytes.Buffer{}
		test.NoError(t.ExecuteContext(&b, c))
		test.AreEqual("test block content", b.String())
	})
}

func TestYieldingInherits(t *testing.T) {
	Within(t, func(test *Test) {
		t, e := New("main").Parse("main", `{{ yield "test_block" }}`, "default")
		test.NoError(e)
		t, e = t.Parse("child", `{{ extends "main" }}{{ block "test_block" }}test block content{{ end_block }}`, "default")
		test.NoError(e)
		c := NewContext(nil)
		c.Main = "child"
		b := bytes.Buffer{}
		test.NoError(t.ExecuteContext(&b, c))
		test.AreEqual("test block content", b.String())
	})
}
