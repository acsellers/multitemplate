package bham

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	textTmpl "text/template"
	"text/template/parse"
)

var (
	dotVarField = `([\.|\$][^\t^\n^\v^\f^\r^ ]+)+`

	simpleValue    = regexp.MustCompile(`^true|false|nil$`)
	simpleField    = regexp.MustCompile(fmt.Sprintf(`^%s$`, dotVarField))
	simpleFunction = regexp.MustCompile(fmt.Sprintf(`^([^\.^\t^\n^\v^\f^\r^ ]+)( %s)*$`, dotVarField))
)

func (pt *protoTree) parseTemplateCode(content string) (*parse.PipeNode, error) {
	switch {
	case content == ".":
		return &parse.PipeNode{
			NodeType: parse.NodePipe,
			Cmds: []*parse.CommandNode{
				&parse.CommandNode{
					NodeType: parse.NodeCommand,
					Args: []parse.Node{
						&parse.DotNode{},
					},
				},
			},
		}, nil
	case simpleValue.MatchString(content):
		return &parse.PipeNode{
			NodeType: parse.NodePipe,
			Cmds: []*parse.CommandNode{
				&parse.CommandNode{
					NodeType: parse.NodeCommand,
					Args: []parse.Node{
						newValueNode(content),
					},
				},
			},
		}, nil
	case simpleField.MatchString(content):
		var arg parse.Node
		if content[0] == '$' {
			arg = newBareVariableNode(content)
		} else {
			arg = newBareFieldNode(content)
		}
		return &parse.PipeNode{
			NodeType: parse.NodePipe,
			Cmds: []*parse.CommandNode{
				&parse.CommandNode{
					NodeType: parse.NodeCommand,
					Args: []parse.Node{
						arg,
					},
				},
			},
		}, nil
	case simpleFunction.MatchString(content):
		return &parse.PipeNode{
			NodeType: parse.NodePipe,
			Cmds: []*parse.CommandNode{
				&parse.CommandNode{
					NodeType: parse.NodeCommand,
					Args:     newBareFunctionNode(content),
				},
			},
		}, nil
	default:
		return pt.processCode(content)
	}
}

func (pt *protoTree) processCode(s string) (*parse.PipeNode, error) {
	t := textTmpl.New("mule").Funcs(textTmpl.FuncMap(pt.funcs))
	t, err := t.Parse("{{" + s + "}}")
	if err != nil {
		return nil, err
	}
	n := t.Tree.Root.Nodes[0]
	if an, ok := n.(*parse.ActionNode); ok {
		return an.Pipe, nil
	}
	return nil, fmt.Errorf("Couldn't extract code for:'%s'", s)
}

func parseFuncArg(arg string) parse.Node {
	// fields (.first.second)
	// variables ($first.second)
	// strings ("value")
	switch arg[0] {
	case '.':
		if arg == "." {
			return &parse.DotNode{}
		} else {
			return &parse.FieldNode{
				NodeType: parse.NodeField,
				Ident:    strings.Split(arg[1:], "."),
			}
		}
	case '$':
		return &parse.VariableNode{
			NodeType: parse.NodeVariable,
			Ident:    strings.Split(arg, "."),
		}
	case '"':
		return &parse.StringNode{
			NodeType: parse.NodeString,
			Quoted:   arg,
			Text:     arg[1 : len(arg)-1],
		}
	}

	// bool variables (true || false)
	// the value nil
	switch arg {
	case "true":
		return &parse.BoolNode{
			NodeType: parse.NodeBool,
			True:     true,
		}
	case "false":
		return &parse.BoolNode{
			NodeType: parse.NodeBool,
			True:     false,
		}
	case "nil":
		return &parse.NilNode{}
	}

	// numeric
	if node, ok := parseNumeric(arg); ok {
		return node

		// function names (blah, url, etc)
	} else {
		return &parse.IdentifierNode{
			NodeType: parse.NodeIdentifier,
			Ident:    arg,
		}
	}
}

func parseNumeric(num string) (*parse.NumberNode, bool) {
	node := &parse.NumberNode{
		NodeType: parse.NodeNumber,
		Text:     num,
	}

	// Following code is adapted from text/template's newNumber function
	// Handle all int's first
	if num[0] != '-' {
		u, err := strconv.ParseUint(num, 0, 64) // will fail for -0; fixed below.
		if err == nil {
			node.IsUint = true
			node.Uint64 = u
			node.IsFloat = true
			node.Float64 = float64(u)
		}
	}

	i, err := strconv.ParseInt(num, 0, 64)
	if err == nil {
		node.IsInt = true
		node.Int64 = i
		if i == 0 {
			node.IsUint = true // in case of -0.
			node.Uint64 = 0
		}
		node.IsFloat = true
		node.Float64 = float64(i)
		return node, true
	}
	if node.IsUint {
		return node, true
	}

	// handle all floats as floats only
	// text template will allow you to turn integer floats into
	// ints/uints, I'm just not feeling it
	f, err := strconv.ParseFloat(num, 64)
	if err == nil {
		node.IsFloat = true
		node.Float64 = f
		return node, true
	}

	return nil, false
}
