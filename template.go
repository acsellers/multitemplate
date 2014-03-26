// multitemplate allows for multiple template parsers that emit
// text/template/parse trees. The interface is intended to be a
// reminiscent of text|html/template. Notable departures are the
// absence of a Delims method and Parse taking a name, src, and
// parser.
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

var Parsers = make(map[string]Parser)

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
	return &Template{Tmpl: template.New(name).Funcs(template.FuncMap{})}
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
	return tmpl.Funcs(GenerateFuncs(t)), nil
}

func (t *Template) Execute(w io.Writer, data interface{}) error {
	if t.ctx == nil {
		t.ctx = &Context{}
	}
	return t.Tmpl.Execute(w, data)
}

func (t *Template) ExecuteTemplate(w io.Writer, name string, data interface{}) error {
	if t.ctx == nil {
		t.ctx = &Context{}
	}

	return t.Tmpl.ExecuteTemplate(w, name, data)
}

func (t *Template) Funcs(fm template.FuncMap) *Template {
	if t.funcs == nil {
		t.funcs = fm
	} else {
		for k, v := range fm {
			t.funcs[k] = v
		}
	}
	return &Template{t.Tmpl.Funcs(fm), t.Base, nil, t.funcs}
}

func (t *Template) Lookup(name string) *Template {
	return &Template{t.Tmpl.Lookup(name), t.Base, nil, t.funcs}
}

func (t *Template) Name() string {
	return t.Tmpl.Name()
}

func (t *Template) Parse(name, src, parser string) (*Template, error) {
	p, ok := Parsers[parser]
	if !ok {
		p = &defaultParser{}
	}

	trees, err := p.ParseTemplate(name, src, t.funcs)
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
	base, exts := Extensions(name)
	if len(exts) > 0 {
		for _, ext := range exts {
			if _, ok := Parsers[ext]; ok {
			} else {
				base = base + "." + ext
			}
		}
	}
	return
}

func Extensions(filename string) (string, []string) {
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
