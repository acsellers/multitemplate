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
}

var errorToken = &token{Type: ErrorToken}

type tokenType int

const (
	ErrorToken tokenType = iota
	TextToken
	HTMLToken
	CommentToken
)

func (t *token) Compile(prefix string) parse.Node {
	switch t.Type {
	case ErrorToken:
		panic("Error token's should not make it to compilation")
	case TextToken:
		return &parse.TextNode{
			NodeType: parse.NodeText,
			Text:     []byte(prefix + t.Content),
		}
	case HTMLToken:
		return &parse.TextNode{
			NodeType: parse.NodeText,
			Text:     []byte(prefix + t.Content),
		}
	case CommentToken:
	default:
	}
	return &parse.TextNode{
		NodeType: parse.NodeText,
	}
}
