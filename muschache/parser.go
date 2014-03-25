package mustache

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

const (
	openBlock = iota
	closeBlock
	elseBlock
	ident
	unescaped
	templateCall
	yield
	noop
	erroring
)

func scanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, data[0 : i+1], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}

func (pt *protoTree) parse() {
	currentWork := &stash{tree: pt}
	scanner := bufio.NewScanner(pt.Reader())
	scanner.Split(scanLines)

	for scanner.Scan() {
		currentWork.Append(scanner.Text())
		if currentWork.hasAction() {
			if currentWork.needsMoreText() {
				continue
			}
		} else {
			continue
		}
		for currentWork.hasAction() && !currentWork.needsMoreText() {
			precedingText, action := currentWork.pullToAction()
			pt.list.Nodes = append(pt.list.Nodes, newTextNode(precedingText))

			switch pt.actionPurpose(action) {
			case ident:
				pt.insertIdentNode(action)
			case templateCall:
				pt.insertTemplateNode(action)
			case yield:
				pt.insertYieldNode(action)
			case openBlock:
				pt.startBlock(action)
			case closeBlock:
				pt.endBlock(action)
			case elseBlock:
				pt.startElseBlock(action)
			}
		}
	}

	if currentWork.hasAction() && currentWork.needsMoreText() {
		pt.err = fmt.Errorf("unterminated delimeter")
	} else {
		if len(currentWork.content) > 0 {
			pt.list.Nodes = append(pt.list.Nodes, newTextNode(currentWork.content))
		}
	}
}

func (pt *protoTree) hasDelims(s string) bool {
	return strings.Index(s, pt.localLeft) < strings.Index(s, pt.localRight) &&
		strings.Index(s, pt.localLeft) >= 0
}

func (pt *protoTree) actionPurpose(w string) int {
	if strings.Contains(w, LeftEscapeDelim) {
		return ident
	}
	w = w[len(pt.localLeft) : len(w)-len(pt.localRight)]
	tw := strings.TrimSpace(w)
	switch tw[0] {
	// start a range/call/if block
	case '#':
		return openBlock

		// end a block
	case '/':
		return closeBlock

		// start an else block
	case '^':
		return elseBlock

		// template/yield
	case '>':
		return templateCall

		// yield block
	case '<':
		return yield

		// switch delimeters
	case '=':
		delims := strings.SplitN(strings.TrimSpace(tw[1:len(tw)-1]), " ", 2)
		if len(delims) != 2 {
			if len(delims)%2 == 0 {
				delims = []string{delims[0][0 : len(delims[0])/2], delims[0][len(delims[0])/2 : len(delims[0])]}
			} else {
				return noop
			}
		}
		pt.localLeft = strings.TrimSpace(delims[0])
		pt.localRight = strings.TrimSpace(delims[1])

		return noop

		// comment block
	case '!':
		return noop

		// .ident block
	case '&':
		return ident

	default:
		return ident
	}
}

func (pt *protoTree) Reader() io.Reader {
	return bytes.NewBufferString(pt.source)
}

func (pt *protoTree) extract(s string) string {
	if strings.HasPrefix(s, LeftEscapeDelim) &&
		strings.HasSuffix(s, RightEscapeDelim) {
		s = s[len(LeftEscapeDelim):]
		s = s[:len(s)-len(RightEscapeDelim)]
	}
	if strings.HasPrefix(s, pt.localLeft) &&
		strings.HasSuffix(s, pt.localRight) {
		s = s[len(pt.localLeft):]
		s = s[:len(s)-len(pt.localRight)]
	}
	switch s[0] {
	case '#', '/', '^', '>', '<', '=', '!', '&':
		s = s[1:]
	}

	return strings.TrimSpace(s)
}

func (pt *protoTree) insertIdentNode(a string) {
	if pt.unescapedAction(a) {
		un := newUnescapedIdentNode(pt.extract(a))
		pt.list.Nodes = append(pt.list.Nodes, un)
	} else {
		ax := pt.extract(a)
		if ax == "." {
			ax = "mustacheItem"
		}
		an := newIdentNode(ax)
		pt.list.Nodes = append(pt.list.Nodes, an)
	}
}

func (pt *protoTree) insertTemplateNode(a string) {
	tn := newTemplateNode(pt.extract(a))
	pt.list.Nodes = append(pt.list.Nodes, tn)
}

func (pt *protoTree) insertYieldNode(a string) {
	yn := newYieldNode(pt.extract(a))
	pt.list.Nodes = append(pt.list.Nodes, yn)
}

func (pt *protoTree) startBlock(a string) {
	tmpl, call, list := newBlockNode(pt.extract(a))
	pt.childTrees = append(pt.childTrees, tmpl)
	pt.list.Nodes = append(pt.list.Nodes, call)
	pt.push(pt.list)
	pt.list = list
}

func (pt *protoTree) endBlock(a string) {
	pt.list = pt.pop()
}

func (pt *protoTree) startElseBlock(a string) {
	ifNode, list := newElseBlock(pt.extract(a))
	pt.list.Nodes = append(pt.list.Nodes, ifNode)
	pt.push(pt.list)
	pt.list = list
}

func (pt *protoTree) unescapedAction(s string) bool {
	return strings.HasPrefix(s, LeftEscapeDelim) ||
		strings.HasPrefix(s, pt.localLeft+"&")
}
