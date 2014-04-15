package bham

import (
	"fmt"
	"strings"
	"text/template"
	"text/template/parse"
)

func findAttrs(s string) (string, string) {
	var openings int
	for i, r := range s {
		switch r {
		case '(':
			openings++
		case ')':
			openings--
			if openings == 0 {
				return s[1:i], s[i+1:]
			}
		}
	}
	return "", s
}

func containsDelimeters(s string) bool {
	return strings.Contains(string(s), RightDelim) &&
		strings.Contains(string(s), LeftDelim)
}

func (pt *protoTree) safeAction(s string) (*parse.ActionNode, error) {
	t := template.New("mule").Funcs(template.FuncMap(pt.funcs))
	t, err := t.Parse(pt.prelude + s)
	if err != nil {
		return nil, err
	}

	n := t.Tree.Root.Nodes[len(t.Tree.Root.Nodes)-1]
	if an, ok := n.(*parse.ActionNode); ok {
		return an, nil
	}
	return nil, fmt.Errorf("Couldn't extract code for:'%s'", s)
}
