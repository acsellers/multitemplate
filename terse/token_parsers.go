package terse

import (
	"fmt"
	"strings"
)

var Doctypes = map[string]string{
	"":             `<!DOCTYPE html>`,
	"Transitional": `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">`,
	"Strict":       `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">`,
	"Frameset":     `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Frameset//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-frameset.dtd">`,
	"5":            `<!DOCTYPE html>`,
	"1.1":          `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN" "http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd">`,
	"Basic":        `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML Basic 1.1//EN" "http://www.w3.org/TR/xhtml-basic/xhtml-basic11.dtd">`,
	"Mobile":       `<!DOCTYPE html PUBLIC "-//WAPFORUM//DTD XHTML Mobile 1.2//EN" "http://www.openmobilealliance.org/tech/DTD/xhtml-mobile12.dtd">`,
	"RDFa":         `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML+RDFa 1.0//EN" "http://www.w3.org/MarkUp/DTD/xhtml-rdfa-1.dtd">`,
}

func doctypeToken(node *rawNode) (*token, error) {
	if len(node.Children) > 0 {
		for _, cNode := range node.Children {
			if !commentCode(cNode.Code) {
				return errorToken, fmt.Errorf("Doctypes may not have nested non-comment lines")
			}
		}
	}
	line := strippedLine(node.Code)
	if line[0:2] != "!!" {
		return errorToken, fmt.Errorf("Doctype token must have two exclamation points")
	}
	line = strippedLine(line[2:])
	if dt, ok := Doctypes[line]; ok {
		td := &token{Content: dt, Type: HTMLToken, Pos: node.Pos}
		ct := &token{Type: CommentToken}
		for _, child := range node.Children {
			if ct.Pos == 0 {
				ct.Pos = node.Pos
			}
			ct.Content += child.Print("\n")
		}
		td.Children = []*token{ct}
		return td, nil
	}
	return errorToken, fmt.Errorf("Doctype token was not found in the list of Doctypes was: '%s'", line)
}

func execToken(node *rawNode) (*token, error) {
	return errorToken, fmt.Errorf("Not Implemented")
}

func verbatimToken(node *rawNode) (*token, error) {
	ct := &token{Type: HTMLToken, Pos: node.Pos}
	if node.Code[0] == '/' {
		ct.Content = node.Code[1:]
		if ct.Content[0] == ' ' {
			ct.Content = ct.Content[1:]
		}
	} else {
		ct.Content = node.Code
	}
	for _, child := range node.Children {
		ct.Content += "\n" + child.Print("  ")
	}
	return ct, nil
}

func commentToken(node *rawNode) (*token, error) {
	ct := &token{Type: CommentToken, Content: node.Code, Pos: node.Pos}
	for _, c := range node.Children {
		ct.Content += "\n" + c.Code
	}

	return ct, nil
}

func tagToken(node *rawNode) (*token, error) {

	return errorToken, fmt.Errorf("Not Implemented")
}

func filterToken(node *rawNode) (*token, error) {

	return errorToken, fmt.Errorf("Not Implemented")
}

func blockToken(node *rawNode) (*token, error) {

	return errorToken, fmt.Errorf("Not Implemented")
}

func defineBlockToken(node *rawNode) (*token, error) {

	return errorToken, fmt.Errorf("Not Implemented")
}

func execBlockToken(node *rawNode) (*token, error) {

	return errorToken, fmt.Errorf("Not Implemented")
}

func yieldToken(node *rawNode) (*token, error) {

	return errorToken, fmt.Errorf("Not Implemented")
}

func ifToken(node *rawNode) (*token, error) {
	if node.Code[0] == '?' {
		node.Code = node.Code[1:]
	}

	it := &token{Type: IfToken, Pos: node.Pos, Content: node.Code}

	var e error
	it.Children, e = childTokenize(node)
	if e != nil {
		return errorToken, e
	}

	return it, nil
}

func elseToken(node *rawNode) (*token, error) {
	et := &token{Type: ElseIfToken, Pos: node.Pos}

	var e error
	et.Children, e = childTokenize(node)
	if e != nil {
		return errorToken, e
	}

	return et, nil
}

func rangeToken(node *rawNode) (*token, error) {
	if node.Code[0] == '&' {
		node.Code = node.Code[1:]
	}

	it := &token{Type: RangeToken, Pos: node.Pos, Content: node.Code}

	var e error
	it.Children, e = childTokenize(node)
	if e != nil {
		return errorToken, e
	}

	return it, nil
}

func rangeElseToken(node *rawNode) (*token, error) {
	et := &token{Type: ElseRangeToken, Pos: node.Pos}

	var e error
	et.Children, e = childTokenize(node)
	if e != nil {
		return errorToken, e
	}

	return et, nil
}

func idClassToken(node *rawNode) (*token, error) {

	return errorToken, fmt.Errorf("Not Implemented")
}

func textToken(node *rawNode) (*token, error) {
	// interpolated text isn't supported yet, so bail out
	if strings.Contains(node.Code, LeftDelim) && strings.Contains(node.Code, RightDelim) {
		return errorToken, fmt.Errorf("Not Implemented")
	}

	// simplest possible text
	if len(node.Children) == 0 {
		return &token{Type: TextToken, Content: node.Code, Pos: node.Pos}, nil
	}

	td := &token{Type: TextToken}
	td.Opening = []*token{
		&token{Type: TextToken, Content: node.Code, Pos: node.Pos},
	}

	for _, child := range node.Children {
		if commentCode(child.Code) {
			t, e := commentToken(child)
			if e != nil {
				return nil, e
			}
			td.Children = append(td.Children, t)
		} else {
			t, e := textToken(child)
			if e != nil {
				return errorToken, fmt.Errorf("Not Implemented")
			}
			td.Children = append(td.Children, t)
		}
	}
	return td, nil
}
