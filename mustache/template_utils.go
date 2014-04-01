package mustache

import (
	"fmt"
	"strings"
	"text/template"
	"text/template/parse"
)

func newpipeline(s string) *parse.PipeNode {
	// take the simplest way of getting text/template to parse it
	// and then steal the result
	t, _ := template.New("mule").Parse("{{" + s + "}}")
	main := t.Tree.Root.Nodes[0]
	if an, ok := main.(*parse.ActionNode); ok {
		return an.Pipe
	} else {
		return nil
	}
}

func newBranchNode(nodeType parse.NodeType, pipe string) parse.BranchNode {
	return parse.BranchNode{
		NodeType: nodeType,
		Pipe:     newpipeline(pipe),
	}
}

func safeAction(s string) (*parse.ActionNode, error) {
	t, e := template.New("mule").Parse(s)
	if e != nil {
		return nil, e
	}
	main := t.Tree.Root.Nodes[len(t.Tree.Root.Nodes)-1]
	if an, ok := main.(*parse.ActionNode); ok {
		return an, nil
	} else {
		return nil, fmt.Errorf("Couldn't find action node")
	}
}

func addEmbeddable(tn *parse.TextNode) []parse.Node {
	output := make([]parse.Node, 0)
	workingText := string(tn.Text)
	for containsDelimeters(workingText) {
		index := strings.Index(workingText, LeftDelim)
		output = append(output, newTextNode(workingText[:index]))
		workingText = workingText[index:]

		index = strings.Index(workingText, RightDelim)
		pipeText := workingText[:index+len(RightDelim)]
		workingText = workingText[index+len(RightDelim):]

		action, e := safeAction(pipeText)
		if e != nil {
			output = append(output, newTextNode(pipeText))
		} else {
			output = append(output, action)
		}
	}

	if workingText != "" {
		tn.Text = []byte(workingText)
		output = append(output, tn)
	}
	return output
}

func containsDelimeters(s string) bool {
	return strings.Contains(string(s), RightDelim) &&
		strings.Contains(string(s), LeftDelim)
}
