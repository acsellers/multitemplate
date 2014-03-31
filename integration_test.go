package multitemplate

import (
	"bytes"
	"html/template"
	"testing"

	. "github.com/acsellers/assert"
)

type tableTest struct {
	Name, Description string
	Expected          string
	Main, Layout      string
	Templates         map[string]string
	Blocks            map[string]string
	Yields            map[string]string
	RenderArgs        map[string]string
}

var tableTests = []tableTest{
	tableTest{
		Name:        "Yielding Blocks",
		Description: "A block of content set in the context must be able to be yielded.",
		Expected:    "test block content",
		Main:        "main",
		Templates: map[string]string{
			"main": `{{ yield "test_block" }}`,
		},
		Blocks: map[string]string{
			"test_block": "test block content",
		},
	},
	tableTest{
		Name:        "Yielding Extends",
		Description: "A block should be able to be yielded in a parent template.",
		Expected:    "test block content",
		Main:        "child",
		Templates: map[string]string{
			"main":  `{{ yield "test_block" }}`,
			"child": `{{ extends "main" }}{{ block "test_block" }}test block content{{ end_block }}`,
		},
	},
	tableTest{
		Name:        "Overriding blocks",
		Description: "A block in an earlier template should override the block in a later template",
		Expected:    "correct",
		Main:        "child",
		Templates: map[string]string{
			"main":  `{{ block "test_block" }}incorrect{{ end_block }}`,
			"child": `{{ extends "main" }}{{ block "test_block" }}correct{{ end_block }}`,
		},
	},
	tableTest{
		Name:        "Blocks between mains and layouts",
		Description: "A block defined in a main template should be available in the Layout template",
		Expected:    "layout correct end",
		Main:        "main",
		Layout:      "layout",
		Templates: map[string]string{
			"layout": `layout {{ yield "test_block" }} end`,
			"main":   `{{ block "test_block" }}correct{{ end_block }}`,
		},
	},
	tableTest{
		Name:        "Layouts with Blocks",
		Description: "Blocks defined in main templates should be visible in layout templates",
		Expected:    "layout main layout",
		Main:        "main",
		Layout:      "layout",
		Templates: map[string]string{
			"main":   `{{ block "testing" }}main{{ end_block }}`,
			"layout": `layout {{ block "testing" }}layout{{ end_block }} layout`,
		},
	},

	tableTest{
		Name:        "Extended Templates with blocks and layouts",
		Description: "Blocks defined in templates that are extended for the main template should show up in the layout template",
		Expected:    "layout child parent layout",
		Main:        "child",
		Layout:      "layout",
		Templates: map[string]string{
			"child":  `{{ extends "parent" }}{{ block "block1" }}child{{ end_block }}`,
			"parent": `{{ block "block1" }}parent{{ end_block }}{{ block "block2" }}parent{{ end_block }}`,
			"layout": `layout {{ block "block1" }}layout{{ end_block }} {{ block "block2" }}layout{{ end_block }} layout`,
		},
	},
}

func TestTables(t *testing.T) {
	Within(t, func(test *Test) {
		for _, tt := range tableTests {
			test.Section(tt.Name)
			t := New("test_templates")
			for name, tmpl := range tt.Templates {
				var e error
				t, e = t.Parse(name, tmpl, "default")
				test.NoError(e)
			}
			c := NewContext(tt.RenderArgs)
			c.Main = tt.Main
			c.Layout = tt.Layout
			for k, b := range tt.Blocks {
				c.Blocks[k] = template.HTML(b)
			}
			for k, t := range tt.Yields {
				c.Yields[k] = t
			}
			b := bytes.Buffer{}
			test.NoError(t.ExecuteContext(&b, c))
			test.AreEqual(tt.Expected, b.String())
		}
	})
}
