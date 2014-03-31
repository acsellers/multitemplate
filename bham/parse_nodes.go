package bham

import (
	"fmt"
	"strings"
	"text/template/parse"
)

func newTree(source, name string) *parse.Tree {
	return &parse.Tree{
		Name:      name,
		ParseName: source,
		Root: &parse.ListNode{
			NodeType: parse.NodeList,
		},
	}
}
func newTextNode(text string) parse.Node {
	return &parse.TextNode{
		NodeType: parse.NodeText,
		Text:     []byte(text),
	}
}

func (pt *protoTree) newMaybeTextNode(text string) []parse.Node {
	if strings.Contains(text, LeftDelim) && strings.Contains(text, RightDelim) {
		output := make([]parse.Node, 0)
		workingText := text
		for containsDelimeters(workingText) {
			index := strings.Index(workingText, LeftDelim)
			output = append(output, newTextNode(workingText[:index]))
			workingText = workingText[index:]

			index = strings.Index(workingText, RightDelim)
			pipeText := workingText[:index+len(RightDelim)]
			workingText = workingText[index+len(RightDelim):]

			action, e := pt.safeAction(pipeText)
			if e != nil {
				output = append(output, newTextNode(pipeText+" "))
			} else {
				output = append(output, action)
			}
		}

		if workingText != "" {
			return append(output, newTextNode(workingText+" "))
		}
		return output
	} else {
		return []parse.Node{
			newTextNode(text + " "),
		}
	}
}

func newValueNode(val string) parse.Node {
	switch val {
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
	panic("Can only call value node for true, false, nil")
}

func newFieldNode(field string) parse.Node {

	return &parse.ActionNode{
		NodeType: parse.NodeAction,
		Pipe: &parse.PipeNode{
			NodeType: parse.NodePipe,
			Cmds: []*parse.CommandNode{
				&parse.CommandNode{
					NodeType: parse.NodeCommand,
					Args: []parse.Node{
						newBareFieldNode(field),
					},
				},
			},
		},
	}
}
func newBareFieldNode(field string) *parse.FieldNode {
	if field[0] == '.' {
		field = field[1:]
	}

	return &parse.FieldNode{
		NodeType: parse.NodeField,
		Ident:    strings.Split(field, "."),
	}
}

func newBareVariableNode(field string) *parse.VariableNode {
	return &parse.VariableNode{
		NodeType: parse.NodeVariable,
		Ident:    strings.Split(field, "."),
	}
}

func newFunctionNode(command string) parse.Node {
	return &parse.ActionNode{
		NodeType: parse.NodeAction,
		Pipe: &parse.PipeNode{
			NodeType: parse.NodePipe,
			Cmds: []*parse.CommandNode{
				&parse.CommandNode{
					NodeType: parse.NodeCommand,
					Args:     newBareFunctionNode(command),
				},
			},
		},
	}
}
func newBareFunctionNode(command string) []parse.Node {
	args := strings.Split(command, " ")
	nodeArgs := []parse.Node{}
	for _, arg := range args {
		switch arg[:1] {
		case ".":
			nodeArgs = append(nodeArgs,
				&parse.FieldNode{
					NodeType: parse.NodeField,
					Ident:    strings.Split(arg[1:], "."),
				},
			)
		case "$":
			nodeArgs = append(nodeArgs,
				&parse.VariableNode{
					NodeType: parse.NodeVariable,
					Ident:    strings.Split(arg, "."),
				},
			)
		default:
			nodeArgs = append(nodeArgs,
				&parse.IdentifierNode{
					NodeType: parse.NodeIdentifier,
					Ident:    arg,
				},
			)
		}
	}

	return nodeArgs
}

func (pt *protoTree) newStandaloneTag(content string) ([]parse.Node, error) {
	td, c, err := pt.parseTag(content)
	if err != nil {
		return []parse.Node{}, err
	}
	return td.Nodes(pt, c)
}

const (
	stateNull = iota
	stateTag
	stateId
	stateClass
	stateAttr
	stateValue
	stateBridge
	stateSpace
	stateDone
)

func (pt *protoTree) parseTag(content string) (tagDescription, string, error) {
	chars := []rune{}
	td := newTagDescription(pt)
	var current, state int
	for _, char := range content {
		chars = append(chars, char)
	}
	var value string
	limit := len(chars)
	{
		switch chars[current] {
		case '.':
			current++
			state = stateClass
			goto preamble_consume
		case '#':
			current++
			state = stateId
			goto preamble_consume
		case '%':
			current++
			state = stateTag
			goto preamble_consume
		}
	preamble_choose:
		if current == limit {
			value = ""
			goto ending
		} else {
			switch chars[current] {
			case '.':
				current++
				state = stateClass
			case '#':
				current++
				state = stateId
			case '(':
				current++
				state = stateAttr
				goto attr_consume
			case '=', '-', ' ':
				goto extra_consume
			}
		}
	preamble_consume:
		for current < limit {
			switch chars[current] {
			case '.', '#':
				td.Add(value, state)
				value = ""
				goto preamble_choose
			case '(':
				current++
				td.Add(value, state)
				value = ""
				goto attr_consume
			case '-', '=', ' ':
				td.Add(value, state)
				value = ""
				goto extra_consume
			default:
				value = value + string(chars[current])
			}
			current++
		}
		td.Add(value, state)
		value = ""
		goto ending
	attr_consume:
		for current < limit {
			switch chars[current] {
			case ' ':
				td.attributes = append(td.attributes, value)
				value = ""
			case '=':
				value = value + "="
				current++
				goto value_consume
			case ')':
				if value != "" {
					td.attributes = append(td.attributes, value)
				}
				current++
				value = ""
				goto extra_consume
			default:
				value = value + string(chars[current])
			}
			current++
		}
		return td, "", fmt.Errorf("Unterminated attributes for %s", content)
	value_consume:
		if current < limit && chars[current] == '"' {
			value = value + "\""
			current++
			quoteIndex := strings.Index(string(chars[current:]), "\"")
			delimIndex := strings.Index(string(chars[current:]), LeftDelim)
			if delimIndex > 0 && quoteIndex > delimIndex {
				leftLen := len([]rune(LeftDelim))
				fmt.Println("here", string(chars[current:]))
				for string(chars[current:current+leftLen]) != LeftDelim {
					value = value + string(chars[current])
					current += 1
				}
				value = value + LeftDelim
				current += len(LeftDelim)
				rightLen := len([]rune(RightDelim))
				searching := true
				for current < limit-rightLen && searching {
					if string(chars[current:current+rightLen]) == RightDelim {
						searching = false
						value = value + RightDelim
						current += rightLen
						td.executableOpen = true
					} else {
						value = value + string(chars[current])
						current++
					}
				}
			}
			for current < limit {
				if chars[current] == '"' {
					td.attributes = append(td.attributes, value+"\"")
					value = ""
					current++
					goto attr_consume
				} else {
					value = value + string(chars[current])
				}
				current++
			}
			return td, "", fmt.Errorf("HTML attribute values must have closing quotation marks for %s", content)
		}
		return td, "", fmt.Errorf("HTML attribute values must have quotation marks for %s", content)
	extra_consume:
		for current < limit {
			value = value + string(chars[current])
			current++
		}
	ending:
		return td, value, nil
	}
}
func newTagDescription(pt *protoTree) tagDescription {
	return tagDescription{
		tag:  "div",
		tree: pt,
	}
}

type tagDescription struct {
	tag            string
	tree           *protoTree
	executableOpen bool
	classes        []string
	idParts        []string
	attributes     []string
}

func (td *tagDescription) Add(content string, state int) {
	if content == "" {
		return
	}
	switch state {
	case stateClass:
		td.classes = append(td.classes, content)
	case stateId:
		td.idParts = append(td.idParts, content)
	case stateTag:
		td.tag = content
	}
}
func (td tagDescription) Opening() string {
	output := fmt.Sprintf("<%s", td.tag)
	if len(td.attributes) > 0 {
		for _, attr := range td.attributes {
			if len(td.classes) > 0 && strings.HasPrefix(attr, "class") {
				output = output +
					" class=\"" +
					strings.Join(td.classes, " ") +
					" " + attr[7:]
				td.classes = []string{}
			} else {
				if len(td.idParts) > 0 && strings.HasPrefix(attr, "id=") {
					output = output + " id=\"" + strings.Join(td.idParts, "_") + "_" + attr[4:]
					td.idParts = []string{}
				} else {
					output = output + " " + attr
				}
			}
		}
	}
	if len(td.classes) > 0 {
		output = output + " class=\"" + strings.Join(td.classes, " ") + "\""
	}
	if len(td.idParts) > 0 {
		output = output + " id=\"" + strings.Join(td.idParts, IdJoin) + "\""
	}
	return output + ">"
}

func (td tagDescription) Close() string {
	return fmt.Sprintf("</%s>", td.tag)
}

func (td tagDescription) Nodes(pt *protoTree, content string) ([]parse.Node, error) {
	if content != "" {
		if content[0] == '=' {
			content = strings.TrimSpace(content[1:])
			node, err := td.tree.parseTemplateCode(content)
			return append(append(pt.newMaybeTextNode(td.Opening()),
				&parse.ActionNode{
					NodeType: parse.NodeAction,
					Pipe:     node,
				}),
				pt.newMaybeTextNode(td.Close())...,
			), err
		} else {
			output := pt.newMaybeTextNode(td.Opening())
			output = append(output, pt.newMaybeTextNode(content)...)
			return append(output, pt.newMaybeTextNode(td.Close())...), nil
		}
	} else {
		return []parse.Node{
			newTextNode(td.Opening()),
			newTextNode(td.Close()),
		}, nil
	}
}
