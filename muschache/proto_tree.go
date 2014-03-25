package mustache

import (
	"text/template/parse"
)

var mangleNum int

type protoTree struct {
	source     string
	tree       *parse.Tree
	childTrees []*parse.Tree
	list       *parse.ListNode
	stack      []*parse.ListNode
	err        error
	localLeft  string
	localRight string
}

func (pt *protoTree) templates() map[string]*parse.Tree {
	output := make(map[string]*parse.Tree)
	for _, tree := range pt.childTrees {
		output[tree.Name] = tree
	}
	output[pt.tree.Name] = pt.tree

	return output
}

func (pt *protoTree) pop() *parse.ListNode {
	if len(pt.stack) == 0 {
		return pt.tree.Root
	}
	ln := pt.stack[len(pt.stack)-1]
	if len(pt.stack) == 1 {
		pt.stack = pt.stack[1:]
	} else {
		pt.stack = pt.stack[:len(pt.stack)-2]
	}
	return ln
}
func (pt *protoTree) push(ln *parse.ListNode) {
	pt.stack = append(pt.stack, ln)
}
