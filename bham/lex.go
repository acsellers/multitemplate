package bham

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

func (pt *protoTree) lex() {
	scanner := bufio.NewScanner(bytes.NewBufferString(pt.source))
	var line, content string
	var currentLevel, nowLevel int
	var currentLine int
	for scanner.Scan() {
		currentLine++
		line = scanner.Text()

		if strings.TrimSpace(line) == "" {
			continue
		}

		nowLevel, content = level(line)
		if currentLevel+1 >= nowLevel {
			lineItem := templateLine{nowLevel, content}
			tempLine := currentLine
			for lineItem.needsContent() {
				if scanner.Scan() {
					tempLine++
					lineItem.content = lineItem.appendContent(scanner.Text())
				} else {
					pt.err = fmt.Errorf("Line %d is not completed", currentLine)
					return
				}
			}
			currentLine = tempLine
			pt.lineList = append(pt.lineList, lineItem)
			currentLevel = nowLevel
		} else {
			pt.err = fmt.Errorf("Line %d is overindented", currentLine)
			return
		}
	}
}

func level(s string) (int, string) {
	var currentLevel int
	for {
		switch s[0] {
		case ' ':
			if s[1] == ' ' {
				s = s[2:]
			} else {
				return currentLevel, s
			}
		case '\t':
			s = s[1:]
		default:
			return currentLevel, s
		}
		currentLevel++
	}
}
