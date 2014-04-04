package multitemplate

import (
	"bytes"
	"html/template"
	"testing"

	. "github.com/acsellers/assert"
)

func TestSentinelJS(t *testing.T) {
	Within(t, func(test *Test) {
		t := New("test_templates")
		var e error
		t, e = t.Parse("default", `<html><script type="text/javascript">{{ yield "break" }}</script></html>`, "default")
		test.NoError(e)

		c := NewContext(nil)
		c.Main = "default"
		c.Blocks["break"] = RenderedBlock{template.HTML(`NONONONO`), HTML}
		b := bytes.Buffer{}
		test.IsError(t.ExecuteContext(&b, c))
		b.Reset()

		c = NewContext(nil)
		c.Main = "default"
		c.Blocks["break"] = RenderedBlock{template.HTML(`NONONONO`), CSS}
		test.IsError(t.ExecuteContext(&b, c))

		c = NewContext(nil)
		c.Main = "default"
		c.Blocks["break"] = RenderedBlock{template.HTML(`YES`), JS}
		test.NoError(t.ExecuteContext(&b, c))
	})
}

func TestSentinelHTML(t *testing.T) {
	Within(t, func(test *Test) {
		t := New("test_templates")
		var e error
		t, e = t.Parse("default", `<html>{{ yield "break" }}</html>`, "default")
		test.NoError(e)

		c := NewContext(nil)
		c.Main = "default"
		c.Blocks["break"] = RenderedBlock{template.HTML(`NONONONO`), HTML}
		b := bytes.Buffer{}
		test.NoError(t.ExecuteContext(&b, c))
		b.Reset()

		c = NewContext(nil)
		c.Main = "default"
		c.Blocks["break"] = RenderedBlock{template.HTML(`NONONONO`), CSS}
		test.IsError(t.ExecuteContext(&b, c))

		c = NewContext(nil)
		c.Main = "default"
		c.Blocks["break"] = RenderedBlock{template.HTML(`YES`), JS}
		test.IsError(t.ExecuteContext(&b, c))
	})
}

func TestSentinelCSS(t *testing.T) {
	Within(t, func(test *Test) {
		t := New("test_templates")
		var e error
		t, e = t.Parse("default", `<html><style>{{ yield "break" }}</style></html>`, "default")
		test.NoError(e)

		c := NewContext(nil)
		c.Main = "default"
		c.Blocks["break"] = RenderedBlock{template.HTML(`NONONONO`), HTML}
		b := bytes.Buffer{}
		test.IsError(t.ExecuteContext(&b, c))
		b.Reset()

		c = NewContext(nil)
		c.Main = "default"
		c.Blocks["break"] = RenderedBlock{template.HTML(`NONONONO`), CSS}
		test.NoError(t.ExecuteContext(&b, c))

		c = NewContext(nil)
		c.Main = "default"
		c.Blocks["break"] = RenderedBlock{template.HTML(`YES`), JS}
		test.IsError(t.ExecuteContext(&b, c))
	})
}
