package terse

import (
	"html/template"
	"text/template/parse"
)

func compile(name string, funcs template.FuncMap, tt tokenTree) (map[string]*parse.Tree, error) {
	if tt.err != nil {
		return map[string]*parse.Tree{}, tt.err
	}

	var nodes []parse.Node
	var prefix string
	if len(tt.roots) > 0 {
		for _, rn := range tt.roots {
			nodes = append(nodes, rn.Compile(prefix)...)
			prefix = "\n"
		}
	}

	pt := &parse.Tree{
		Name: name,
		Root: &parse.ListNode{
			NodeType: parse.NodeList,
			Pos:      parse.Pos(0),
			Nodes:    nodes,
		},
	}

	return map[string]*parse.Tree{name: pt}, nil
}
