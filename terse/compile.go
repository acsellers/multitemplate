package terse

import (
	"html/template"
	"text/template/parse"
)

func compile(name string, funcs template.FuncMap, tt tokenTree) (map[string]*parse.Tree, error) {
	if tt.err != nil {
		return map[string]*parse.Tree{}, tt.err
	}

	setResources(tt.roots, &resources{funcs: funcs})

	pt := &parse.Tree{
		Name: name,
		Root: &parse.ListNode{
			NodeType: parse.NodeList,
			Pos:      parse.Pos(0),
			Nodes:    compileTokens(tt.roots, ""),
		},
	}

	return map[string]*parse.Tree{name: pt}, nil
}

func setResources(t []*token, r *resources) {
	for _, c := range t {
		c.Rsc = r
		setResources(c.Opening, r)
		setResources(c.Closing, r)
		setResources(c.Children, r)
	}
}

func compileTokens(ts []*token, prefix string) []parse.Node {
	var nodes []parse.Node
	var hold []parse.Node
	watchType := ErrorToken
	if len(ts) > 0 {
		for _, rn := range ts {
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
				case ElseWithToken:
					if ln, ok := hold[0].(*parse.WithNode); ok {
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
			if !cf && len(prefix) == 0 {
				prefix = "\n"
			}
		}
		if len(hold) > 0 {
			nodes = append(nodes, hold...)
		}
	}
	return nodes
}
