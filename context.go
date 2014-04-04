package multitemplate

import (
	"bytes"
	"html/template"
	"io"
	"strings"
)

func NewContext(data interface{}) *Context {
	c := &Context{}
	c.Yields = make(map[string]string)
	c.Blocks = make(map[string]RenderedBlock)
	c.Dot = data
	c.output = newPouchWriter()
	return c
}

type RenderedBlock struct {
	Content template.HTML
	Type    Ruleset
}

type Ruleset string

const (
	User Ruleset = ""
	HTML Ruleset = "&lt;&#34;&#39;."
	CSS  Ruleset = "ZgotmplZ"
	JS   Ruleset = `"\u003c\"'."`
)

// A Context allows you to setup more specialized template executions,
// like those involving layouts
type Context struct {
	// Main template to be rendered, not layout
	Main        string
	mainContent RenderedBlock
	// Layout for rendering
	Layout          string
	executingLayout bool
	currentMode     string

	// Templates set for yields
	Yields map[string]string
	// Blocks (pre-rendered HTML)
	Blocks map[string]RenderedBlock
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
	canNest := !c.output.nesting()

	hasParent := c.parent != ""
	forthcomingLayout := c.Layout != "" && !c.executingLayout
	inactiveView := strings.TrimSpace(c.output.root.String()) == ""
	openableTemplate := hasParent || (forthcomingLayout && inactiveView)

	return canNest && openableTemplate
}

func (c *Context) exec(name string, dot interface{}) (RenderedBlock, error) {
	b := bytes.Buffer{}
	// Replace the output buffer so we don't have stale data hanging around
	// We need to have the rest of the context hanging around.
	temp := c.output
	c.output = newPouchWriter()

	e := c.tmpl.ExecuteTemplate(&b, name, dot)
	c.output = temp
	return RenderedBlock{template.HTML(b.String()), HTML}, e
}

func (c *Context) execWithFallback(name string, f fallback, dot interface{}) (RenderedBlock, error) {
	if c.Yields[name] != "" {
		return c.exec(c.Yields[name], dot)
	}
	if rb, ok := c.Blocks[name]; ok {
		return rb, nil
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
	if c.output.err != nil {
		return c.output.err
	}
	_, e := io.WriteString(w, c.output.root.String())
	return e
}

type fallback string
