package terse

import (
	"html/template"
	"text/template/parse"
)

func compile(name string, funcs template.FuncMap, tt tokenTree) (map[string]*parse.Tree, error) {
	if tt.err != nil {
		return map[string]*parse.Tree{}, tt.err
	}

	pt := &parse.Tree{
		Name: name,
		Root: &parse.ListNode{
			NodeType: parse.NodeList,
			Pos:      parse.Pos(0),
		},
	}
	for _, rn := range tt.roots {
		pt.Root.Nodes = append(pt.Root.Nodes, rn.Compile())
	}
	return map[string]*parse.Tree{name: pt}, nil
}
