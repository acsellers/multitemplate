package multitemplate

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template/parse"
)

// To Register a parser, set it in this map for each extension that would
// correspond to it.
var Parsers = make(map[string]Parser)

// The interface you must have to implement a Parser
type Parser interface {
	ParseTemplate(name, src string, funcs template.FuncMap) (map[string]*parse.Tree, error)
	String() string
}

type Template struct {
	Tmpl  *template.Template
	Base  string
	ctx   *Context
	funcs template.FuncMap
}

func Must(t *Template, err error) *Template {
	if err != nil {
		panic(err)
	}
	return t
}

func New(name string) *Template {
	return &Template{Tmpl: template.New(name).Funcs(template.FuncMap{}), Base: name}
}

func ParseFiles(filenames ...string) (*Template, error) {
	return (&Template{Tmpl: template.New("root")}).ParseFiles(filenames...)
}

func ParseGlob(pattern string) (*Template, error) {
	return (&Template{Tmpl: template.New("root")}).ParseGlob(pattern)
}

func (t *Template) AddParseTree(name string, tree *parse.Tree) (*Template, error) {
	var e error
	t.Tmpl, e = t.Tmpl.AddParseTree(name, tree)
	return t, e
}

func (t *Template) Clone() (*Template, error) {
	tmpl, err := t.Tmpl.Clone()
	return &Template{tmpl, t.Base, nil, t.funcs}, err
}

func (t *Template) Context(ctx *Context) (*Template, error) {
	tmpl, err := t.Clone()
	if err != nil {
		return nil, err
	}
	tmpl.ctx = ctx
	ctx.tmpl = tmpl
	return tmpl.Funcs(generateFuncs(tmpl)), nil
}

func (t *Template) Execute(w io.Writer, data interface{}) error {
	var tt *Template
	if t.ctx == nil {
		tt, _ = t.Context(NewContext(data))
	}

	e := tt.Tmpl.Execute(tt.ctx.output, data)
	if e == nil {
		return tt.ctx.close(w)
	}
	return e
}

func (t *Template) ExecuteContext(w io.Writer, ctx *Context) error {
	tt, e := t.Context(ctx)
	if e != nil {
		return e
	}

	main := ctx.Main
	if ctx.Layout != "" {
		ctx.mainContent, e = tt.ctx.exec(ctx.Main, ctx.Dot)
		if e != nil {
			return e
		}
		main = ctx.Layout
	}
	return tt.Tmpl.ExecuteTemplate(w, main, ctx.Dot)
}

func (t *Template) ExecuteTemplate(w io.Writer, name string, data interface{}) error {
	tt := t
	if t.ctx == nil {
		tt, _ = t.Context(NewContext(data))
	}

	e := tt.Tmpl.ExecuteTemplate(tt.ctx.output, name, data)
	if e == nil {
		return tt.ctx.close(w)
	}
	return e
}

func (t *Template) Funcs(fm template.FuncMap) *Template {
	if t.funcs == nil {
		t.funcs = fm
	} else {
		for k, v := range fm {
			t.funcs[k] = v
		}
	}
	t.Tmpl.Funcs(fm)
	return t
}

func (t *Template) Lookup(name string) *Template {
	tmpl := t.Tmpl.Lookup(name)
	if tmpl != nil {
		return &Template{tmpl, t.Base, nil, t.funcs}
	}
	return nil
}

func (t *Template) Name() string {
	return t.Tmpl.Name()
}

func (t *Template) Parse(name, src, parser string) (*Template, error) {
	p, ok := Parsers[parser]
	if !ok {
		p = &defaultParser{}
	}

	t2, _ := t.Clone()
	trees, err := p.ParseTemplate(name, src, t2.Funcs(generateFuncs(t)).funcs)
	if err != nil {
		return nil, err
	}
	for n, tree := range trees {
		t, err = t.AddParseTree(n, tree)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

func (t *Template) ParseFiles(filenames ...string) (*Template, error) {
	for _, f := range filenames {
		n, p := t.stripBase(f)
		b, e := ioutil.ReadFile(f)
		if e != nil {
			return t, e
		}
		t, e = t.Parse(n, string(b), p)
		if e != nil {
			return t, e
		}
	}
	return t, nil
}

func (t *Template) ParseGlob(pattern string) (*Template, error) {
	filenames, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	if len(filenames) == 0 {
		return nil, fmt.Errorf("multitemplate: no files match pattern: %#q", pattern)
	}
	return t.ParseFiles(filenames...)
}

func (t *Template) stripBase(filename string) (name, parser string) {
	if strings.HasPrefix(filename, t.Base) {
		filename = filename[len(t.Base):]
	}
	if filename[0] == '/' || filename[0] == '\\' {
		name = filename[1:]
	}
	base, exts := extensions(name)
	if len(exts) > 0 {
		for _, ext := range exts {
			if _, ok := Parsers[ext]; ok {
				parser = ext
			} else {
				base = base + "." + ext
			}
		}
	}
	name = base
	return
}

func extensions(filename string) (string, []string) {
	name := filepath.Base(filename)
	dirs := filename[:len(filename)-len(name)]
	ext := filepath.Ext(name)
	exts := []string{}
	for ext != "" {
		exts = append(exts, ext[1:])
		name = name[:len(name)-len(ext)]
		ext = filepath.Ext(name)
	}
	return dirs + name, exts
}

func (t *Template) Templates() []*Template {
	tmpls := t.Tmpl.Templates()
	ret := make([]*Template, len(tmpls))
	for i, tmpl := range tmpls {
		ret[i] = &Template{tmpl, t.Base, nil, t.funcs}
	}
	return ret
}
