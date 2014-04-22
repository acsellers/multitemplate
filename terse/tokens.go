package terse

import (
	"reflect"
	"strings"
	"text/template/parse"
)

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
	Rsc      *resources
}

var errorToken = &token{Type: ErrorToken}

type tokenType int

const (
	ErrorToken tokenType = iota
	TextToken
	HTMLToken
	ExecToken
	BlockToken
	CommentToken
	IfToken
	ElseIfToken
	RangeToken
	ElseRangeToken
	TagToken
	FilterToken
	FilterContentToken
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

			an, _ := actionNode(t.Content, t.Rsc)
			bn.Pipe = an.Pipe
			bn.List = &parse.ListNode{
				NodeType: parse.NodeList,
				Pos:      parse.Pos(t.Children[0].Pos),
				Nodes:    t.ChildCompile(prefix),
			}

			return []parse.Node{
				&parse.IfNode{bn},
			}
		case RangeToken:
			bn := parse.BranchNode{
				NodeType: parse.NodeRange,
				Pos:      parse.Pos(t.Pos),
			}

			an, _ := actionNode(t.Content, t.Rsc)
			bn.Pipe = an.Pipe
			if len(prefix) == 0 || prefix[0] != '\n' {
				prefix = "\n" + prefix
			}
			bn.List = &parse.ListNode{
				NodeType: parse.NodeList,
				Pos:      parse.Pos(t.Children[0].Pos),
				Nodes:    t.ChildCompile(prefix),
			}

			return []parse.Node{
				&parse.RangeNode{bn},
			}
		case ExecToken:
			n, e := actionNode(t.Content, t.Rsc)
			if e != nil {
				return []parse.Node{}
			}
			na := []parse.Node{n}
			if len(t.Children) == 0 {
				return na
			}
			na = append(na, t.ChildCompile(prefix+"  ")...)

			c := n.Pipe.Cmds[len(n.Pipe.Cmds)-1]
			if in, ok := c.Args[0].(*parse.IdentifierNode); ok {
				if tf, ok := t.Rsc.funcs["end_"+in.Ident]; ok {
					rv := reflect.TypeOf(tf)
					if rv.Kind() == reflect.Func {
						na = append(na, &parse.TextNode{
							NodeType: parse.NodeText,
							Text:     []byte("\n"),
						})

						if rv.NumIn() == 0 {
							n, _ := actionNode("end_"+in.Ident, t.Rsc)
							na = append(na, n)
						} else {
							n, e := actionNode("end_"+strings.TrimSpace(t.Content), t.Rsc)
							if e == nil {
								na = append(na, n)
							}
						}
					}
				}
			}
			return na
		case FilterToken:
			filter := Filters[t.Content]
			content := ""
			for _, c := range t.Children {
				content += "\n" + c.Content
			}
			if len(content) > 0 {
				content = content[1:]
			}
			converted, _ := filter(content)
			return []parse.Node{
				&parse.TextNode{
					NodeType: parse.NodeText,
					Text:     []byte(converted),
				},
			}
		}
	} else {
		switch t.Type {
		case TextToken:
			ps := []parse.Node{}
			ps = append(ps, t.OpeningCompile(prefix)...)
			return append(ps, t.ChildCompile(prefix+"  ")...)
		case HTMLToken, ExecToken, TagToken:
			ps := []parse.Node{}
			ps = append(ps, t.OpeningCompile(prefix)...)
			ps = append(ps, t.ChildCompile(prefix+"  ")...)
			return append(ps, t.ClosingCompile(prefix)...)
		case ElseIfToken:
			return t.ChildCompile(prefix)
		case ElseRangeToken:
			return t.ChildCompile(prefix)
		case BlockToken:
			ps := []parse.Node{}
			ps = append(ps, t.OpeningCompile(prefix)...)
			ps = append(ps, t.ChildCompile(prefix)...)
			ps = append(ps, t.ClosingCompile(prefix)...)
			return ps
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
	for _, ot := range t.Closing {
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
	return compileTokens(t.Children, prefix)
}
