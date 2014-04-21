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
		Name:        "Yielding Extend",
		Description: "A block should be able to be yielded in a parent template.",
		Expected:    "test block content",
		Main:        "child",
		Templates: map[string]string{
			"main":  `{{ yield "test_block" }}`,
			"child": `{{ extend "main" }}{{ block "test_block" }}test block content{{ end_block }}`,
		},
	},
	tableTest{
		Name:        "Overriding blocks",
		Description: "A block in an earlier template should override the block in a later template",
		Expected:    "correct",
		Main:        "child",
		Templates: map[string]string{
			"main":  `{{ block "test_block" }}incorrect{{ end_block }}`,
			"child": `{{ extend "main" }}{{ block "test_block" }}correct{{ end_block }}`,
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
			"child":  `{{ extend "parent" }}{{ block "block1" }}child{{ end_block }}`,
			"parent": `{{ block "block1" }}parent{{ end_block }}{{ block "block2" }}parent{{ end_block }}`,
			"layout": `layout {{ block "block1" }}layout{{ end_block }} {{ block "block2" }}layout{{ end_block }} layout`,
		},
	},

	tableTest{
		Name:        "Nested blocks in Parent templates",
		Description: "Blocks must not nest, and inner blocks should execute as if it was the last template",
		Expected:    "layout parent child layout",
		Main:        "child",
		Layout:      "layout",
		Templates: map[string]string{
			"child":  `{{ extend "parent" }}{{ block "child_block" }}child{{ end_block }}`,
			"parent": `{{ block "test_block" }}parent {{ block "child_block" }}parent{{ end_block }}{{ end_block }}`,
			"layout": `layout {{ block "test_block" }}layout{{ end_block }} layout`,
		},
	},

	tableTest{
		Name:        "Extended templates in both Main and Layout templates",
		Description: "Both Main and Layout templates should be able to be extended",
		Expected:    "one two three four",
		Main:        "main_child",
		Layout:      "layout_child",
		Templates: map[string]string{
			"layout_child":  `{{ extend "layout_parent" }}{{ block "child_layout" }}two{{ end_block }}`,
			"layout_parent": `one {{ yield "child_layout" }} {{ yield "parent_main" }} {{ yield "child_main" }}`,
			"main_child":    `{{ extend "main_parent" }}{{ block "child_main" }}four{{ end_block }}`,
			"main_parent":   `{{ block "parent_main" }}three{{ end_block }}`,
		},
	},

	tableTest{
		Name:        "Default content in layout template",
		Description: "If a a block was not set in the main template, we should see the default content",
		Expected:    "layout default main",
		Main:        "main",
		Layout:      "layout",
		Templates: map[string]string{
			"layout": `layout {{ block "nonexistent" }}default{{ end_block }} {{ block "content" }}content{{ end_block }}`,
			"main":   `{{ block "content" }}main{{ end_block }}`,
		},
	},

	tableTest{
		Name:        "Extended Main Templates with Layouts",
		Description: "Main template must be able to be yielded in a layout",
		Expected:    "lay ext child end out",
		Main:        "main_child",
		Layout:      "layout",
		Templates: map[string]string{
			"main_child":  `{{ extend "main_parent" }}{{ block "content" }}child{{ end_block }}`,
			"main_parent": `ext {{ block "content" }}error{{ end_block }} end`,
			"layout":      `lay {{ yield }} out`,
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
				c.Blocks[k] = RenderedBlock{template.HTML(b), User}
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
