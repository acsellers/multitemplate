package bham

import (
	"fmt"
	"strings"
	"text/template/parse"
)

func (pt *protoTree) compile() {
	cleanName := pt.name
	i := strings.Index(pt.name, ".bham")
	if i >= 0 {
		cleanName = pt.name[:i] + pt.name[i+5:]
	}

	pt.outputTree = newTree(pt.name, cleanName)

	pt.compileToList(pt.outputTree.Root, pt.nodes)
}

func (pt *protoTree) compileToList(arr *parse.ListNode, nodes []protoNode) {
	for _, node := range nodes {
		switch node.identifier {
		case identRaw:
			arr.Nodes = append(arr.Nodes, newTextNode(node.content))
		case identFilter:
			if node.needsRuntimeData() {
			} else {
				content := node.filter.Open + node.filter.Handler(node.content) + node.filter.Close
				arr.Nodes = append(arr.Nodes, newTextNode(content))
			}
		case identExecutable:
			node, err := pt.parseTemplateCode(node.content)
			if err == nil {
				arr.Nodes = append(arr.Nodes, &parse.ActionNode{
					NodeType: parse.NodeAction,
					Pipe:     node,
				})
			} else {
				pt.err = err
			}
		case identTag:
			nodes, err := pt.newStandaloneTag(node.content)
			if err == nil {
				arr.Nodes = append(arr.Nodes, nodes...)
			} else {
				pt.err = err
			}
		case identTagOpen:
			td, c, err := pt.parseTag(node.content)
			if err == nil {
				arr.Nodes = append(arr.Nodes, newTextNode(td.Opening()))
				if c != "" {
					arr.Nodes = append(arr.Nodes, newMaybeTextNode(c)...)
				}
			} else {
				pt.err = err
			}
		case identTagClose:
			td, _, err := pt.parseTag(node.content)
			if err == nil {
				arr.Nodes = append(arr.Nodes, newTextNode(td.Close()))
			} else {
				pt.err = err
			}
		case identText:
			arr.Nodes = append(arr.Nodes, newMaybeTextNode(node.content)...)
		case identIf:
			branching, err := pt.parseTemplateCode(node.content)
			if err == nil {
				in := &parse.IfNode{
					parse.BranchNode{
						NodeType: parse.NodeIf,
						Pipe:     branching,
						List: &parse.ListNode{
							NodeType: parse.NodeList,
						},
						ElseList: &parse.ListNode{
							NodeType: parse.NodeList,
						},
					},
				}
				if len(node.list) > 0 {
					pt.compileToList(in.List, node.list)
				}
				if len(node.elseList) > 0 {
					pt.compileToList(in.ElseList, node.elseList)
				}
				arr.Nodes = append(arr.Nodes, in)
			} else {
				pt.err = err
			}

		case identRange:
			branching, err := pt.parseTemplateCode(node.content)
			if err == nil {
				in := &parse.RangeNode{
					parse.BranchNode{
						NodeType: parse.NodeIf,
						Pipe:     branching,
						List: &parse.ListNode{
							NodeType: parse.NodeList,
						},
						ElseList: &parse.ListNode{
							NodeType: parse.NodeList,
						},
					},
				}
				if len(node.list) > 0 {
					pt.compileToList(in.List, node.list)
				}
				if len(node.elseList) > 0 {
					pt.compileToList(in.ElseList, node.elseList)
				}
				arr.Nodes = append(arr.Nodes, in)
			} else {
				pt.err = err
			}

		case identWith:
			branching, err := pt.parseTemplateCode(node.content)
			if err == nil {
				wn := &parse.WithNode{
					parse.BranchNode{
						NodeType: parse.NodeWith,
						Pipe:     branching,
						List: &parse.ListNode{
							NodeType: parse.NodeList,
						},
						ElseList: &parse.ListNode{
							NodeType: parse.NodeList,
						},
					},
				}
				if len(node.list) > 0 {
					pt.compileToList(wn.List, node.list)
				}
				arr.Nodes = append(arr.Nodes, wn)
			} else {
				pt.err = err
			}

		default:
			fmt.Println(node.identifier)
			fmt.Println(node.content)
		}
	}
}
