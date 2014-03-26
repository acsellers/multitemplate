package multitemplate

import (
	"bytes"
	"html/template"
)

type Context struct {
	// Templates set for yields
	Yields map[string]string
	// Data (strings), set for yields
	Content map[string]template.HTML
	// Base RenderArgs for the template
	Dot interface{}

	// internal, for simpler functions
	tmpl *Template
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

type Fallback string
