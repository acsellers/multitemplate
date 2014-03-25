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

func safeAction(s string) (*parse.ActionNode, error) {
	// take the simplest way of getting text/template to parse it
	// and then steal the result
	if varUse.MatchString(s) {
		for _, varUser := range varUse.FindAllStringSubmatch(s, -1) {
			s = "{{ " + varUser[0] + " := 0 }}" + s
		}
	}
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
