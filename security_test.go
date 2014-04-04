package multitemplate

import (
	"bytes"
	"html/template"
	"testing"

	. "github.com/acsellers/assert"
)

func TestSentinel(t *testing.T) {
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
