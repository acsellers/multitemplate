package bham

import (
	"fmt"
	"strings"
)

func (pt *protoTree) analyze() {
	pt.doAnalyze(0, len(pt.lineList)-1)
	pt.nodes = pt.currNodes
}

func (pt *protoTree) doAnalyze(currentIndex, finalIndex int) {
	for currentIndex <= finalIndex {
		if pt.err != nil {
			return
		}

		line := pt.lineList[currentIndex]

		switch {
		case line.accept("%.#"):
			currentIndex = pt.tagLike(currentIndex, finalIndex)
			continue
		case line.accept("=-"):
			currentIndex = pt.actionableLine(currentIndex, finalIndex)
			continue
		case line.accept(":"):
			for _, handler := range Filters {
				if line.content == handler.Trigger {
					currentIndex = pt.followHandler(currentIndex+1, finalIndex, handler)
					return
				}
			}
			pt.err = fmt.Errorf("Bad handler: %s", line.content)
			return
		case line.prefix("!!!"):
			pt.insertDoctype(line)
			currentIndex++
		default:
			pt.insertText(line)
			currentIndex++
		}
	}
}

func (pt *protoTree) insertDoctype(line templateLine) {
	doctype, ok := Doctypes[strings.TrimSpace(line.content[3:])]
	if ok {
		pt.insertRaw(doctype, line.indentation)
	} else {
		pt.err = fmt.Errorf("Bad doctype, details: '%s'", line.content)
	}
}

func (pt *protoTree) followHandler(startIndex, finalIndex int, handler FilterHandler) int {
	lines := make([]string, 0)
	index := startIndex
	base := pt.lineList[startIndex].indentation
	for index <= finalIndex && base <= pt.lineList[index].indentation {
		diff := base - pt.lineList[index].indentation
		lines = append(lines, pad(diff)+pt.lineList[index].content)
		index++
	}
	pt.insertFilter(
		strings.Join(lines, "\n"),
		pt.lineList[startIndex].indentation,
		handler,
	)

	return index
}

func pad(indent int) string {
	var output string
	for i := 0; i < indent; i++ {
		output = output + "  "
	}
	return output
}

func (pt *protoTree) actionableLine(startIndex, finalIndex int) int {
	currentIndex := startIndex + 1
	var endIndex int
	if pt.lineList[startIndex].blockParameter() {
		for currentIndex <= finalIndex && pt.lineList[startIndex].indentation < pt.lineList[currentIndex].indentation {
			currentIndex++
		}
		parentNodes := pt.currNodes
		pt.currNodes = []protoNode{}
		pt.doAnalyze(startIndex+1, currentIndex-1)
		primaryNodes, secondaryNodes := pt.currNodes, []protoNode{}

		if pt.lineList[startIndex].mightHaveElse() && currentIndex < finalIndex {
			if pt.lineList[currentIndex].isElse() {
				endIndex = currentIndex
				currentIndex++
				for pt.lineList[endIndex].indentation < pt.lineList[currentIndex].indentation && currentIndex < finalIndex {
					currentIndex++
				}
				pt.currNodes = []protoNode{}
				pt.doAnalyze(endIndex+1, currentIndex)
				secondaryNodes = pt.currNodes
				currentIndex++
			}
		}
		pt.currNodes = parentNodes
		switch {
		case pt.lineList[startIndex].isIf():
			pt.insertIf(
				pt.lineList[startIndex].after("-=").without("if ").String(),
				pt.lineList[startIndex].indentation,
				primaryNodes,
				secondaryNodes,
			)
		case pt.lineList[startIndex].isUnless():
			pt.insertIf(
				pt.lineList[startIndex].after("-=").without("unless ").String(),
				pt.lineList[startIndex].indentation,
				primaryNodes,
				secondaryNodes,
			)
		case pt.lineList[startIndex].isRange():
			pt.insertRange(
				pt.lineList[startIndex].after("-=").without("range ").String(),
				pt.lineList[startIndex].indentation,
				primaryNodes,
				secondaryNodes,
			)
		case pt.lineList[startIndex].isWith():
			pt.insertWith(
				pt.lineList[startIndex].after("-=").without("with ").String(),
				pt.lineList[startIndex].indentation,
				primaryNodes,
			)
		}
	} else {
		pt.insertExecutable(
			pt.lineList[startIndex].after("-=").String(),
			pt.lineList[startIndex].indentation,
		)
	}

	return currentIndex
}

func (pt *protoTree) tagLike(currentIndex, finalIndex int) int {
	if finalIndex == currentIndex || pt.lineList[currentIndex+1].indentation <= pt.lineList[currentIndex].indentation {
		pt.currNodes = append(pt.currNodes, protoNode{
			level:      pt.lineList[currentIndex].indentation,
			identifier: identTag,
			content:    pt.lineList[currentIndex].content,
		})
		return currentIndex + 1
	} else {
		pt.currNodes = append(pt.currNodes, protoNode{
			level:      pt.lineList[currentIndex].indentation,
			identifier: identTagOpen,
			content:    pt.lineList[currentIndex].content,
		})
		tagIndex := currentIndex + 1
		for tagIndex < finalIndex && pt.lineList[tagIndex].indentation > pt.lineList[currentIndex].indentation {
			tagIndex++
		}
		if pt.lineList[tagIndex].indentation <= pt.lineList[currentIndex].indentation {
			tagIndex--
		}
		pt.doAnalyze(currentIndex+1, tagIndex)
		pt.currNodes = append(pt.currNodes, protoNode{
			level:      pt.lineList[currentIndex].indentation,
			identifier: identTagClose,
			content:    pt.lineList[currentIndex].content,
		})
		return tagIndex + 1
	}
}
