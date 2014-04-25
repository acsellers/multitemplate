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
	tt := tokenize(scan(src))
	if tt.err != nil {
		return map[string]*parse.Tree{}, tt.err
	}
	return compile(name, funcs, tt)
}
func (*multiStruct) String() string {
	return "terse: HTML Templating gone concise"
}
