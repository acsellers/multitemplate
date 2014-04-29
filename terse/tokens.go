package terse

import (
	"reflect"
	"strings"
	"text/template/parse"
)

type tokenTree struct {
	roots []*token
	defs  []*token
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
	WithToken
	ElseWithToken
	TagToken
	TemplateToken
	FilterToken
	FilterContentToken
	DefineToken
)

func (t *token) Compile(prefix string) []parse.Node {
	if t.Type == ErrorToken {
		panic("Error token's should not make it to compilation")
	}

	if t.Content != "" {
		switch t.Type {
		case TemplateToken:
			nodes := []parse.Node{}
			if prefix != "" {
				nodes = append(nodes,
					&parse.TextNode{
						NodeType: parse.NodeText,
						Text:     []byte(prefix),
						Pos:      parse.Pos(t.Pos),
					},
				)
			}
			an, e := actionNode(t.Children[0].Content, t.Rsc, t.Children[0].Pos)
			if e != nil {
				t.Rsc.err = e
				return []parse.Node{}
			}
			return append(nodes,
				&parse.TemplateNode{
					Pos:      parse.Pos(t.Pos),
					NodeType: parse.NodeTemplate,
					Name:     t.Content,
					Pipe:     an.Pipe,
				},
			)
		case DefineToken:
			t.Rsc.tt.defs = append(t.Rsc.tt.defs, t)
		case TextToken:
			return textNodes(prefix+t.Content, t.Rsc, t.Pos)
		case HTMLToken:
			return textNodes(prefix+t.Content, t.Rsc, t.Pos)
		case IfToken:
			bn := parse.BranchNode{
				NodeType: parse.NodeIf,
				Pos:      parse.Pos(t.Pos),
			}

			an, e := actionNode(t.Content, t.Rsc, t.Pos)
			if e != nil {
				t.Rsc.err = e
				return []parse.Node{}
			}
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
			var an *parse.ActionNode
			var e error
			var vars []string
			if doubleRangeRegex.MatchString(t.Content) {
				sm := doubleRangeRegex.FindStringSubmatch(t.Content)[1:]
				vars = []string{sm[2], sm[1]}
				t.Rsc.vars = append(t.Rsc.vars, sm[1:]...)
				an, e = actionNode(sm[0], t.Rsc, t.Pos)
			} else if singleRangeRegex.MatchString(t.Content) {
				sm := singleRangeRegex.FindStringSubmatch(t.Content)[1:]
				vars = sm[1:]
				t.Rsc.vars = append(t.Rsc.vars, sm[1:]...)
				an, e = actionNode(sm[0], t.Rsc, t.Pos)
			} else {
				an, e = actionNode(t.Content, t.Rsc, t.Pos)
			}

			if e != nil {
				t.Rsc.err = e
				return []parse.Node{}
			}
			for _, vd := range vars {
				an.Pipe.Decl = append(an.Pipe.Decl,
					&parse.VariableNode{
						NodeType: parse.NodeVariable,
						Ident:    strings.Split(vd, "."),
						Pos:      parse.Pos(t.Pos),
					},
				)
			}
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
		case WithToken:
			bn := parse.BranchNode{
				NodeType: parse.NodeWith,
				Pos:      parse.Pos(t.Pos),
			}
			var an *parse.ActionNode
			var e error
			var vars []string
			if singleRangeRegex.MatchString(t.Content) {
				sm := singleRangeRegex.FindStringSubmatch(t.Content)[1:]
				vars = sm[1:]
				t.Rsc.vars = append(t.Rsc.vars, sm[1:]...)
				an, e = actionNode(sm[0], t.Rsc, t.Pos)
			} else {
				an, e = actionNode(t.Content, t.Rsc, t.Pos)
			}
			if e != nil {
				return []parse.Node{}
			}

			for _, vd := range vars {
				an.Pipe.Decl = append(an.Pipe.Decl,
					&parse.VariableNode{
						NodeType: parse.NodeVariable,
						Ident:    strings.Split(vd, "."),
						Pos:      parse.Pos(t.Pos),
					},
				)
			}
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
				&parse.WithNode{bn},
			}

		case ExecToken:
			n, e := actionNode(t.Content, t.Rsc, t.Pos)
			if e != nil {
				t.Rsc.err = e
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
							Pos:      parse.Pos(t.Pos),
						})

						if rv.NumIn() == 0 {
							n, e := actionNode("end_"+in.Ident, t.Rsc, t.Pos)
							if e != nil {
								t.Rsc.err = e
							}
							na = append(na, n)
						} else {
							n, e := actionNode("end_"+strings.TrimSpace(t.Content), t.Rsc, t.Pos)
							if e != nil {
								t.Rsc.err = e
							} else {
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
			converted, interpolate := filter(content)
			if interpolate {
				return textNodes(converted, t.Rsc, t.Pos)
			} else {
				return []parse.Node{
					&parse.TextNode{
						NodeType: parse.NodeText,
						Text:     []byte(converted),
						Pos:      parse.Pos(t.Pos),
					},
				}
			}
		}
	} else {
		switch t.Type {
		case TextToken:
			ps := []parse.Node{}
			ps = append(ps, t.OpeningCompile(prefix)...)
			return append(ps, t.ChildCompile(prefix+"  ")...)
		case HTMLToken, ExecToken:
			ps := []parse.Node{}
			ps = append(ps, t.OpeningCompile(prefix)...)
			ps = append(ps, t.ChildCompile(prefix+"  ")...)
			return append(ps, t.ClosingCompile(prefix)...)
		case ElseIfToken:
			return t.ChildCompile(prefix)
		case ElseRangeToken:
			return t.ChildCompile(prefix)
		case ElseWithToken:
			return t.ChildCompile(prefix)
		case BlockToken, TagToken:
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
	case WithToken:
		return true, ElseWithToken
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
