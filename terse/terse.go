package terse

import (
	"html/template"
	"text/template/parse"

	"github.com/acsellers/multitemplate"
)

var (
	LeftDelim  = "{{"
	RightDelim = "}}"
)

func init() {
	ms := multiStruct{}
	multitemplate.Parsers["terse"] = &ms
}

type multiStruct struct{}

func (*multiStruct) ParseTemplate(name, src string, funcs template.FuncMap) (map[string]*parse.Tree, error) {
	return compile(name, funcs, tokenize(scan(src)))
}
func (*multiStruct) String() string {
	return "terse: HTML Templating gone concise"
}
