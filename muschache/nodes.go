package mustache

import (
	"fmt"
	"strings"
	"text/template/parse"
)

func newTextNode(s string) *parse.TextNode {
	return &parse.TextNode{
		NodeType: parse.NodeText,
		Text:     []byte(s),
	}
}

func newTemplateNode(w string) *parse.TemplateNode {
	return &parse.TemplateNode{
		NodeType: parse.NodeTemplate,
		Name:     w,
		Pipe: &parse.PipeNode{
			NodeType: parse.NodePipe,
			Cmds: []*parse.CommandNode{
				&parse.CommandNode{
					NodeType: parse.NodeCommand,
					Args: []parse.Node{
						&parse.DotNode{},
					},
				},
			},
		},
	}
}

func newBlockNode(a string) (*parse.Tree, *parse.IfNode, *parse.ListNode) {
	tmplName := fmt.Sprintf("mustacheAnonymous%d", mangleNum)
	mangleNum++
	startList := []parse.Node{
		&parse.ActionNode{
			NodeType: parse.NodeAction,
			Pipe: &parse.PipeNode{
				NodeType: parse.NodePipe,
				Decl: []*parse.VariableNode{
					&parse.VariableNode{
						NodeType: parse.NodeVariable,
						Ident:    []string{"$mustacheCurrent"},
					},
				},
				Cmds: []*parse.CommandNode{
					&parse.CommandNode{
						NodeType: parse.NodeCommand,
						Args:     []parse.Node{&parse.DotNode{}},
					},
				},
			},
		},
	}

	tree := &parse.Tree{
		Name:      tmplName,
		ParseName: a,
		Root: &parse.ListNode{
			NodeType: parse.NodeList,
			Nodes:    startList,
		},
	}
	return tree, newBlockChooseNode(tmplName, a), tree.Root
}

func newElseBlock(f string) (*parse.IfNode, *parse.ListNode) {
	listNode := &parse.ListNode{
		NodeType: parse.NodeList,
	}
	ifNode := &parse.IfNode{
		parse.BranchNode{
			NodeType: parse.NodeIf,
			Pipe: &parse.PipeNode{
				NodeType: parse.NodePipe,
				Cmds: []*parse.CommandNode{
					&parse.CommandNode{
						NodeType: parse.NodeCommand,
						Args: []parse.Node{
							&parse.FieldNode{
								NodeType: parse.NodeField,
								Ident:    strings.Split(f, "."),
							},
						},
					},
				},
			},
			List: &parse.ListNode{
				NodeType: parse.NodeList,
			},
			ElseList: listNode,
		},
	}
	return ifNode, listNode
}

func newBlockChooseNode(tmpl, field string) *parse.IfNode {
	return &parse.IfNode{
		parse.BranchNode{
			NodeType: parse.NodeIf,
			Pipe:     newIdentNode(field).Pipe,
			List: &parse.ListNode{
				NodeType: parse.NodeList,
				Nodes: []parse.Node{
					&parse.IfNode{
						parse.BranchNode{
							NodeType: parse.NodeIf,
							Pipe: &parse.PipeNode{
								NodeType: parse.NodePipe,
								Cmds: []*parse.CommandNode{
									&parse.CommandNode{
										NodeType: parse.NodeCommand,
										Args: []parse.Node{
											&parse.FieldNode{
												NodeType: parse.NodeField,
												Ident:    strings.Split(field, "."),
											},
										},
									},
									&parse.CommandNode{
										NodeType: parse.NodeCommand,
										Args: []parse.Node{
											&parse.IdentifierNode{
												NodeType: parse.NodeIdentifier,
												Ident:    "mustacheIsCollection",
											},
										},
									},
								},
							},
							List: &parse.ListNode{
								NodeType: parse.NodeList,
								Nodes: []parse.Node{
									&parse.RangeNode{
										parse.BranchNode{
											NodeType: parse.NodeRange,
											Pipe: &parse.PipeNode{
												NodeType: parse.NodePipe,
												Cmds: []*parse.CommandNode{
													&parse.CommandNode{
														NodeType: parse.NodeCommand,
														Args: []parse.Node{
															&parse.FieldNode{
																NodeType: parse.NodeField,
																Ident:    []string{field},
															},
														},
													},
												},
											},
											List: &parse.ListNode{
												NodeType: parse.NodeList,
												Nodes: []parse.Node{
													&parse.TemplateNode{
														NodeType: parse.NodeTemplate,
														Name:     tmpl,
														Pipe: &parse.PipeNode{
															NodeType: parse.NodePipe,
															Cmds: []*parse.CommandNode{
																&parse.CommandNode{
																	NodeType: parse.NodeCommand,
																	Args: []parse.Node{
																		&parse.IdentifierNode{
																			NodeType: parse.NodeIdentifier,
																			Ident:    "mustacheUpscope",
																		},
																		&parse.VariableNode{
																			NodeType: parse.NodeVariable,
																			Ident:    []string{"$mustacheCurrent"},
																		},
																		&parse.DotNode{},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
							ElseList: &parse.ListNode{
								NodeType: parse.NodeList,
								Nodes: []parse.Node{
									&parse.TemplateNode{
										NodeType: parse.NodeTemplate,
										Name:     tmpl,
										Pipe: &parse.PipeNode{
											NodeType: parse.NodePipe,
											Cmds: []*parse.CommandNode{
												&parse.CommandNode{
													NodeType: parse.NodeCommand,
													Args: []parse.Node{
														&parse.IdentifierNode{
															NodeType: parse.NodeIdentifier,
															Ident:    "mustacheUpscope",
														},
														&parse.VariableNode{
															NodeType: parse.NodeVariable,
															Ident:    []string{"$mustacheCurrent"},
														},
														&parse.FieldNode{
															NodeType: parse.NodeField,
															Ident:    strings.Split(field, "."),
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func newIdentNode(field string) *parse.ActionNode {
	return newActionNodeForCommands(
		newCommandFieldNode(
			strings.Split(field, ".")...,
		),
	)
}

func newYieldNode(w string) *parse.ActionNode {
	args := []parse.Node{
		&parse.IdentifierNode{
			NodeType: parse.NodeIdentifier,
			Ident:    "yield",
		},
	}
	tw := strings.TrimSpace(w)
	if tw != "" {
		args = append(args, &parse.StringNode{
			NodeType: parse.NodeString,
			Quoted:   tw,
			Text:     tw,
		})
	}
	args = append(args, &parse.DotNode{})

	return newActionNodeForCommands(&parse.CommandNode{
		NodeType: parse.NodeCommand,
		Args:     args,
	})
}

func newUnescapedIdentNode(field string) *parse.ActionNode {
	return newActionNodeForCommands(
		newCommandFieldNode(strings.Split(field, ".")...),
		newCommandIdentifierNode("mustacheUnescape"),
	)
}

func newActionNodeForCommands(cn ...*parse.CommandNode) *parse.ActionNode {
	return &parse.ActionNode{
		NodeType: parse.NodeAction,
		Pipe: &parse.PipeNode{
			NodeType: parse.NodePipe,
			Cmds:     cn,
		},
	}
}

func newCommandFieldNode(fields ...string) *parse.CommandNode {
	return &parse.CommandNode{
		NodeType: parse.NodeCommand,
		Args: []parse.Node{
			&parse.FieldNode{
				NodeType: parse.NodeField,
				Ident:    fields,
			},
		},
	}
}
func newCommandIdentifierNode(ident string) *parse.CommandNode {
	return &parse.CommandNode{
		NodeType: parse.NodeCommand,
		Args: []parse.Node{
			&parse.IdentifierNode{
				NodeType: parse.NodeIdentifier,
				Ident:    ident,
			},
		},
	}
}
