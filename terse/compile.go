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
	if len(tt.roots) > 0 {
		nodes = []parse.Node{tt.roots[0].Compile("")}
		for _, rn := range tt.roots[1:] {
			nodes = append(nodes, rn.Compile("\n"))
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
