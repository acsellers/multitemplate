package mustache

import (
	ht "html/template"
	"strings"
	"text/template/parse"

	"github.com/acsellers/multitemplate"
)

var (
	LeftDelim        = "{{"
	RightDelim       = "}}"
	LeftEscapeDelim  = "{{{"
	RightEscapeDelim = "}}}"
)

type multiStruct struct{}

func (ms *multiStruct) ParseTemplate(name, src string, funcs ht.FuncMap) (map[string]*parse.Tree, error) {
	return Parse(name, src, funcs)
}
func (ms *multiStruct) String() string {
	return "mustache: Logic-less templates"
}

func init() {
	ms := multiStruct{}
	multitemplate.Parsers["mustache"] = &ms
}

func Parse(templateName, templateContent string, funcs ht.FuncMap) (map[string]*parse.Tree, error) {
	i := strings.Index(templateName, ".mustache")
	name := templateName[:i] + templateName[i+len(".mustache"):]

	proto := &protoTree{
		source:     templateContent,
		localRight: RightDelim,
		localLeft:  LeftDelim,
		funcs:      funcs,
		tree: &parse.Tree{
			Name:      name,
			ParseName: templateName,
			Root: &parse.ListNode{
				NodeType: parse.NodeList,
				Nodes: []parse.Node{
					&parse.ActionNode{
						NodeType: parse.NodeAction,
						Pipe: &parse.PipeNode{
							NodeType: parse.NodePipe,
							Decl: []*parse.VariableNode{
								&parse.VariableNode{
									NodeType: parse.NodeVariable,
									Ident:    []string{"$mustacheCurrent"},
								},
							},
							Cmds: []*parse.CommandNode{
								&parse.CommandNode{
									NodeType: parse.NodeCommand,
									Args:     []parse.Node{&parse.DotNode{}},
								},
							},
						},
					},
				},
			},
		},
	}
	proto.list = proto.tree.Root
	proto.parse()

	return proto.templates(), proto.err
}
