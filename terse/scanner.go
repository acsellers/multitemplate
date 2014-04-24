package terse

import (
	"bufio"
	"strings"
)

type rawTree struct {
	Children []*rawNode
}

func (rt rawTree) String() string {
	switch len(rt.Children) {
	case 0:
		return ""
	case 1:
		return rt.Children[0].Print("")
	}
	content := rt.Children[0].Print("")
	for _, node := range rt.Children[1:] {
		content += "\n" + node.Print("")
	}
	return content
}

type rawNode struct {
	Code     string
	Children []*rawNode
	Pos      int
}

func (rn rawNode) Print(prefix string) string {
	content := prefix + rn.Code
	for _, node := range rn.Children {
		content += "\n" + node.Print(prefix+"  ")
	}
	return content
}

func scan(src string) rawTree {
	s := bufio.NewScanner(strings.NewReader(src))
	rt := rawTree{}
	for s.Scan() {
		line := s.Text()
		if blankLine(line) {
			continue
		}
		if unindentedLine(line) || len(rt.Children) == 0 {
			rt.Children = append(rt.Children, &rawNode{Code: line})
			continue
		}
		current := rt.Children[len(rt.Children)-1]
		line = unindentLine(line)
		for len(current.Children) != 0 && !unindentedLine(line) {
			current = current.Children[len(current.Children)-1]
			line = unindentLine(line)
		}
		current.Children = append(current.Children, &rawNode{Code: line})
	}
	return rt
}

func blankLine(line string) bool {
	return strings.TrimSpace(line) == ""
}

func unindentedLine(line string) bool {
	return line[0:2] != "  " && line[0] != '\t'
}

func unindentLine(line string) string {
	if line[0:2] == "  " {
		return line[2:]
	}
	if line[0] == '/' {
		return line[1:]
	}
	return line
}
