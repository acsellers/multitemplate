package bham

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"text/template/parse"
	"unicode"
)

var (
	dotVarField = `([\.|\$][^\t^\n^\v^\f^\r^ ]+)+`

	simpleValue    = regexp.MustCompile(`true|false|nil`)
	simpleField    = regexp.MustCompile(fmt.Sprintf(`^%s$`, dotVarField))
	simpleFunction = regexp.MustCompile(fmt.Sprintf(`^([^\.^\t^\n^\v^\f^\r^ ]+)( %s)*$`, dotVarField))
)

func parseTemplateCode(content string) (*parse.PipeNode, error) {
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
		return processCode(content)
	}
}

func processCode(s string) (*parse.PipeNode, error) {
	var chars []rune
	for _, char := range s {
		chars = append(chars, char)
	}

	current, last := 0, len(chars)-1
	continuing := true
	work := new(bytes.Buffer)
	var funcName string
	var args []string
	var nodeArgs []parse.Node
	var nodes []*parse.CommandNode
	declNodes := []*parse.VariableNode{}

	{
	begin_command:
		for current <= last && unicode.IsSpace(chars[current]) {
			current++
		}
		for current <= last && !unicode.IsSpace(chars[current]) {
			work.WriteRune(chars[current])
			current++
		}
		funcName = work.String()
		work.Reset()
		for current <= last && unicode.IsSpace(chars[current]) {
			current++
		}
		if current >= last {
			goto complete
		}
		goto choose_after

	func_arg:
		if chars[current] == '"' {
			work.WriteRune(chars[current])
			current++
			for current <= last && chars[current] != '"' && chars[current-1] != '\\' {
				work.WriteRune(chars[current])
				current++
			}
			if chars[current] != '"' {
				return nil, fmt.Errorf("Unterminated string: %s", work.String())
			}
			work.WriteRune(chars[current])
			current++
		} else {
			for current <= last && !unicode.IsSpace(chars[current]) {
				work.WriteRune(chars[current])
				current++
			}
		}

		for current <= last && unicode.IsSpace(chars[current]) {
			current++
		}
		args = append(args, work.String())
		work.Reset()

		if current >= last {
			goto complete
		}

	choose_after:
		if chars[current] == '|' {
			goto push_func
		} else {
			if chars[current] == ':' && chars[current+1] == '=' {
				current = current + 2
				goto convert_to_assignment
			} else {
				goto func_arg
			}
		}
	convert_to_assignment:
		declNodes = append(declNodes, &parse.VariableNode{
			NodeType: parse.NodeVariable,
			Ident:    strings.Split(funcName, "."),
		})
		funcName = ""
		for _, arg := range args {
			declNodes = append(declNodes, &parse.VariableNode{
				NodeType: parse.NodeVariable,
				Ident:    strings.Split(arg, "."),
			})
		}
		args = []string{}

		goto begin_command

	complete:
		continuing = false

	push_func:
		nodeArgs = []parse.Node{}
		switch funcName[0] {
		case '.':
			if len(funcName) > 1 {
				nodeArgs = append(nodeArgs,
					&parse.FieldNode{
						NodeType: parse.NodeField,
						Ident:    strings.Split(funcName[1:], "."),
					})
			} else {
				nodeArgs = append(nodeArgs, &parse.DotNode{})
			}
		case '$':
			nodeArgs = append(nodeArgs,
				&parse.VariableNode{
					NodeType: parse.NodeVariable,
					Ident:    strings.Split(funcName, "."),
				})
		case '"':
			nodeArgs = append(nodeArgs,
				&parse.StringNode{
					NodeType: parse.NodeString,
					Quoted:   funcName,
					Text:     funcName[1 : len(funcName)-1],
				},
			)
		default:
			nodeArgs = append(nodeArgs,
				&parse.IdentifierNode{
					NodeType: parse.NodeIdentifier,
					Ident:    funcName,
				})
		}
		for _, arg := range args {
			nodeArgs = append(nodeArgs, parseFuncArg(arg))
		}
		funcName = ""
		nodes = append(nodes, &parse.CommandNode{
			NodeType: parse.NodeCommand,
			Args:     nodeArgs,
		})

		args = []string{}

		if continuing {
			goto begin_command
		}
	}

	if current < last {
		return nil, fmt.Errorf("Could not complete parse")
	}
	if funcName != "" {
		return nil, fmt.Errorf("Parse was not able to complete")
	}

	return &parse.PipeNode{
		NodeType: parse.NodePipe,
		Decl:     declNodes,
		Cmds:     nodes,
	}, nil
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
