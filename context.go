package multitemplate

import (
	"bytes"
	"html/template"
	"io"
)

func NewContext(data interface{}) *Context {
	c := &Context{}
	c.Yields = make(map[string]string)
	c.Blocks = make(map[string]template.HTML)
	c.Dot = data
	c.output = newPouchWriter()
	return c
}

// A Context allows you to setup more specialized template executions,
// like those involving layouts
type Context struct {
	// Main template to be rendered, not layout
	Main        string
	mainContent template.HTML
	// Layout for rendering
	Layout          string
	executingLayout bool

	// Templates set for yields
	Yields map[string]string
	// Blocks (pre-rendered HTML)
	Blocks map[string]template.HTML
	// Base RenderArgs for the template
	Dot interface{}

	// Name of the parent template
	parent string
	// internal, for exec
	tmpl *Template
	// blocks need this
	output *pouchWriter
}

func (c *Context) openableScope() bool {
	return !c.output.nesting() && (c.parent != "" || (c.Layout != "" && !c.executingLayout))
}

func (c *Context) exec(name string, dot interface{}) (template.HTML, error) {
	b := bytes.Buffer{}
	e := c.tmpl.ExecuteTemplate(&b, name, dot)
	return template.HTML(b.String()), e
}

func (c *Context) execWithFallback(name string, f fallback, dot interface{}) (template.HTML, error) {
	if c.Yields[name] != "" {
		return c.exec(c.Yields[name], dot)
	}
	if c.Blocks[name] != "" {
		return c.Blocks[name], nil
	}
	return c.exec(string(f), dot)
}

func (c *Context) Close(w io.Writer) error {
	if c.parent != "" {
		temp := c.parent
		c.parent = ""
		for temp != "" {
			c.output.Reset()
			e := c.tmpl.Tmpl.ExecuteTemplate(c.output, temp, c.Dot)
			if e != nil {
				return e
			}
			temp = c.parent
		}
	}
	_, e := io.WriteString(w, c.output.root.String())
	return e
}

type fallback string
