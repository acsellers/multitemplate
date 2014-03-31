package multitemplate

import (
	"html/template"
	textTmpl "text/template"
	"text/template/parse"
)

// Delimeters for the standard Go template parser
var GoLeftDelim, GoRightDelim string

type defaultParser struct {
	left, right string
}

func (ms *defaultParser) ParseTemplate(name, src string, funcs template.FuncMap) (map[string]*parse.Tree, error) {
	var t *textTmpl.Template
	var e error
	tf := textTmpl.FuncMap(funcs)
	if GoRightDelim != "" || GoLeftDelim != "" {
		t, e = textTmpl.New(name).Funcs(tf).Delims(GoLeftDelim, GoRightDelim).Parse(src)
	} else {
		t, e = textTmpl.New(name).Funcs(tf).Parse(src)
	}
	if e != nil {
		return nil, e
	}

	ret := make(map[string]*parse.Tree)
	for _, t := range t.Templates() {
		ret[t.Name()] = t.Tree
	}
	return ret, nil
}

func (ms *defaultParser) String() string {
	return "html/template: Standard Library Template"
}

func init() {
	Parsers["tmpl"] = &defaultParser{}
}
