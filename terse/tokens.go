package terse

import "text/template/parse"

type tokenTree struct {
	roots []*token
	err   error
}

type token struct {
	Opening  []*token
	Closing  []*token
	Children []*token
	Type     tokenType
	Content  string
	Pos      int
}

var errorToken = &token{Type: ErrorToken}

type tokenType int

const (
	ErrorToken tokenType = iota
	TextToken
	HTMLToken
	ExecToken
	CommentToken
	IfToken
	ElseIfToken
	RangeToken
	ElseRangeToken
	TagToken
	TagOpenToken
	TagCloseToken
)

func (t *token) Compile(prefix string) []parse.Node {
	if t.Type == ErrorToken {
		panic("Error token's should not make it to compilation")
	}

	if t.Content != "" {
		switch t.Type {
		case TextToken:
			return []parse.Node{
				&parse.TextNode{
					NodeType: parse.NodeText,
					Text:     []byte(prefix + t.Content),
				},
			}
		case HTMLToken:
			return []parse.Node{
				&parse.TextNode{
					NodeType: parse.NodeText,
					Text:     []byte(prefix + t.Content),
				},
			}
		case IfToken:
			bn := parse.BranchNode{
				NodeType: parse.NodeIf,
				Pos:      parse.Pos(t.Pos),
			}

			an, _ := actionNode(t.Content, &resources{})
			bn.Pipe = an.Pipe
			bn.List = &parse.ListNode{
				NodeType: parse.NodeList,
				Pos:      parse.Pos(t.Children[0].Pos),
				Nodes:    t.ChildCompile(prefix),
			}

			return []parse.Node{
				&parse.IfNode{bn},
			}
		}
	} else {
		switch t.Type {
		case TextToken:
			ps := []parse.Node{}
			ps = append(ps, t.OpeningCompile(prefix)...)
			return append(ps, t.ChildCompile(prefix+"  ")...)
		case HTMLToken:
			ps := []parse.Node{}
			ps = append(ps, t.OpeningCompile(prefix)...)
			ps = append(ps, t.ChildCompile(prefix+"  ")...)
			return append(ps, t.ClosingCompile(prefix)...)
		case ElseIfToken:
			return t.ChildCompile(prefix)
		}
	}
	return []parse.Node{}
}

func (t *token) FollowupToken() (bool, tokenType) {
	switch t.Type {
	case IfToken:
		return true, ElseIfToken
	case RangeToken:
		return true, ElseRangeToken
	}
	return false, ErrorToken
}
func (t *token) OpeningCompile(prefix string) []parse.Node {
	ns := []parse.Node{}
	for _, ot := range t.Opening {
		ns = append(ns, ot.Compile(prefix)...)
	}
	return ns
}

func (t *token) ClosingCompile(prefix string) []parse.Node {
	ns := []parse.Node{}
	for _, ot := range t.Opening {
		ns = append(ns, ot.Compile(prefix)...)
	}
	return ns
}

func (t *token) ChildCompile(prefix string) []parse.Node {
	if len(prefix) > 0 && prefix[0] != '\n' {
		prefix = "\n" + prefix
	}
	if len(t.Children) == 0 {
		return []parse.Node{}
	}
	nodes := []parse.Node{}
	for _, child := range t.Children {
		nodes = append(nodes, child.Compile(prefix)...)
	}
	return nodes
}
