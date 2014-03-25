package multitemplate

import (
	"text/template"
	"text/template/parse"

	"github.com/acsellers/multitemplate"
)

type defaultParser struct{}

func (ms *multiStruct) ParseTemplate(name, src string, funcs template.FuncMap) (map[string]*parse.Tree, error) {
	t, e := template.New(name).Funcs(funcs).Parse(src)
	if e != nil {
		return nil, e
	}
	ret := make(map[string]*parse.Tree)
	for _, t := range t.Templates() {
		ret[t.Name()] = t.Tree
	}
	return ret, nil
}

func (ms *multiStruct) String() string {
	return "html/template: Standard Library Template"
}

func init() {
	ms := multiStruct{}
	multitemplate.Parsers["default"] = &ms
	multitemplate.Parsers["tmpl"] = &ms
}
