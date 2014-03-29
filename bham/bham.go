package bham

import (
	"html/template"
	"strings"
	"text/template/parse"

	"github.com/acsellers/multitemplate"
)

var Doctypes = map[string]string{
	"":             `<!DOCTYPE html>`,
	"Transitional": `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">`,
	"Strict":       `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">`,
	"Frameset":     `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Frameset//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-frameset.dtd">`,
	"5":            `<!DOCTYPE html>`,
	"1.1":          `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN" "http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd">`,
	"Basic":        `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML Basic 1.1//EN" "http://www.w3.org/TR/xhtml-basic/xhtml-basic11.dtd">`,
	"Mobile":       `<!DOCTYPE html PUBLIC "-//WAPFORUM//DTD XHTML Mobile 1.2//EN" "http://www.openmobilealliance.org/tech/DTD/xhtml-mobile12.dtd">`,
	"RDFa":         `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML+RDFa 1.0//EN" "http://www.w3.org/MarkUp/DTD/xhtml-rdfa-1.dtd">`,
}

type multiStruct struct{}

func (ms *multiStruct) ParseTemplate(name, src string, funcs template.FuncMap) (map[string]*parse.Tree, error) {
	return Parse(name, src)
}
func (ms *multiStruct) String() string {
	return "bham: Blocky Hypertext Abstraction Markup"
}

func init() {
	ms := multiStruct{}
	multitemplate.Parsers["bham"] = &ms
}

// parse will return a parse tree containing a single
func Parse(name, text string) (map[string]*parse.Tree, error) {
	pt := &protoTree{source: text, name: name}
	pt.lex()
	pt.analyze()
	pt.compile()

	return map[string]*parse.Tree{name: pt.outputTree}, pt.err
}

type protoTree struct {
	name       string
	source     string
	lineList   []templateLine
	nodes      []protoNode
	currNodes  []protoNode
	tokenList  []token
	outputTree *parse.Tree
	err        error
}

type protoNode struct {
	level      int
	identifier int
	content    string
	filter     FilterHandler
	list       []protoNode
	elseList   []protoNode
}

func (pn protoNode) needsRuntimeData() bool {
	return strings.Contains(pn.content, LeftDelim) &&
		strings.Contains(pn.content, RightDelim)
}
