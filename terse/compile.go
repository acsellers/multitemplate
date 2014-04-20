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
	var hold []parse.Node
	var prefix string
	watchType := ErrorToken
	if len(tt.roots) > 0 {
		for _, rn := range tt.roots {
			cf, ft := rn.FollowupToken()
			switch {
			case cf:
				if len(hold) > 0 {
					nodes = append(nodes, hold...)
				}
				watchType = ft
				hold = rn.Compile(prefix)
			case watchType == rn.Type:
				en := rn.Compile(prefix)
				switch watchType {
				case ElseRangeToken:
					if ln, ok := hold[0].(*parse.RangeNode); ok {
						ln.ElseList = &parse.ListNode{
							NodeType: parse.NodeList,
							Nodes:    en,
						}
					}
				case ElseIfToken:
					if ln, ok := hold[0].(*parse.IfNode); ok {
						ln.ElseList = &parse.ListNode{
							NodeType: parse.NodeList,
							Nodes:    en,
						}
					}
				}
				nodes = append(nodes, hold...)
				hold = []parse.Node{}
				watchType = ErrorToken
			case watchType != ErrorToken:
				watchType = ErrorToken
				nodes = append(nodes, hold...)
				hold = []parse.Node{}
			default:
				nodes = append(nodes, rn.Compile(prefix)...)
			}
			prefix = "\n"
		}
		if len(hold) > 0 {
			nodes = append(nodes, hold...)
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
