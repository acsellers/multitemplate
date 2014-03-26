package multitemplate

import (
	"bytes"
	"html/template"
	"io"
)

func NewContext(t *Template, data interface{}) *Context {
	c := &Context{tmpl: t}
	c.Yields = make(map[string]string)
	c.Content = make(map[string]template.HTML)
	c.Dot = data
	c.Output = NewPouchWriter()
	return c
}

type Context struct {
	// Templates set for yields
	Yields map[string]string
	// Data (strings), set for yields
	Content map[string]template.HTML
	// Base RenderArgs for the template
	Dot interface{}

	// Name of the parent template
	parent string
	// internal, for exec
	tmpl *Template
	// blocks need this
	Output *PouchWriter
}

func (c *Context) exec(name string, dot interface{}) (template.HTML, error) {
	b := bytes.Buffer{}
	e := c.tmpl.ExecuteTemplate(&b, name, dot)
	return template.HTML(b.String()), e
}

func (c *Context) execWithFallback(name string, f Fallback, dot interface{}) (template.HTML, error) {
	if c.Yields[name] != "" {
		return c.exec(c.Yields[name], dot)
	}
	if c.Content[name] != "" {
		return c.Content[name], nil
	}
	return c.exec(string(f), dot)
}

func (c *Context) Close(w io.Writer) error {
	if c.parent != "" {
		temp := c.parent
		c.parent = ""
		for temp != "" {
			c.Output.Reset()
			e := c.tmpl.ExecuteTemplate(w, temp, c.Dot)
			if e != nil {
				return e
			}
			temp = c.parent
		}
	}
	_, e := c.Output.root.WriteTo(w)
	return e
}

type Fallback string
