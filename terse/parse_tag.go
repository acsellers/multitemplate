package terse

import (
	"fmt"
	"regexp"
	"strings"
)

func parseTag(node *rawNode, children bool) (*tag, error) {
	t := &tag{
		Name:     "div",
		Node:     node,
		Source:   node.Code,
		Attrs:    make(map[string]string),
		DynAttrs: make(map[string]string),
	}
	e := t.Parse(children)
	return t, e
}

type tag struct {
	Node      *rawNode
	Source    string
	Name      string
	ChildTags []string
	Id        string
	Classes   []string
	Attrs     map[string]string
	DynAttrs  map[string]string
	Remaining string
	Enclosing bool
}

var attrStartRegex = regexp.MustCompile(`^([a-zA-Z0-9_-]+)=`)

func (t *tag) Parse(children bool) error {
	t.Enclosing = children
	// pull percentage sign off
	if t.Source[0] == '%' {
		t.Source = t.Source[1:]
	}
	// pull element name
	if t.Source[0] != '#' && t.Source[0] != '.' {
		t.Name = firstTextToken(t.Source)
		t.Source = t.Source[len(t.Name):]
	}
	// pull id and class static's
	for len(t.Source) > 0 && (t.Source[0] == '#' || t.Source[0] == '.') {
		switch t.Source[0] {
		case '#':
			t.Id = firstTextToken(t.Source[1:])
			t.Source = t.Source[len(t.Id)+1:]
		case '.':
			cl := firstTextToken(t.Source[1:])
			t.Source = t.Source[len(cl)+1:]
			t.Classes = append(t.Classes, cl)
		}
	}

	if len(t.Source) > 0 && t.Source[0] == '(' {
		if len(t.Node.Children) == 0 {
			return fmt.Errorf("Missing nested attributes for %s", t.Node.Code)
		}
		for _, child := range t.Node.Children {
			if attrStartRegex.MatchString(strings.TrimSpace(child.Code)) {
				t.Source = child.Code
				e := t.parseAttributes()
				if e != nil {
					return e
				}
			} else if strings.HasPrefix(strings.TrimSpace(child.Code), ")") {
				t.Node.Children = t.Node.Children[1:]
				t.Remaining = strings.TrimSpace(child.Code)[1:]
				break
			}

			t.Node.Children = t.Node.Children[1:]
			if strings.HasPrefix(strings.TrimSpace(t.Source), ")") {
				t.Remaining = strings.TrimSpace(t.Source)[1:]
				break
			}
			if t.Remaining != "" {
				break
			}
		}
	} else if attrStartRegex.MatchString(strings.TrimSpace(t.Source)) {
		e := t.parseAttributes()
		if e != nil {
			return e
		}
	}
	t.Remaining = t.Source
	if len(t.Remaining) > 0 && t.Remaining[0] == ' ' {
		t.Remaining = t.Remaining[1:]
	}
	for len(t.Remaining) > 0 && t.Remaining[0] == '>' {
		check := strings.TrimSpace(t.Remaining[1:])
		switch {
		case len(check) == 0:
			return fmt.Errorf("Missing tag in nested tags for %s", t.Node.Code)
		case check[0] == '%':
			element := firstTextToken(check[1:])
			check = check[len(element)+1:]
			t.Remaining = strings.TrimSpace(check)
			t.ChildTags = append(t.ChildTags, element)
		case ValidElements[firstTextToken(check)]:
			element := firstTextToken(check)
			check = check[len(element):]
			t.Remaining = strings.TrimSpace(check)
			t.ChildTags = append(t.ChildTags, element)
		}
	}
	return nil
}

func (t *tag) Open() *token {
	tc := t.Start()
	if t.Enclosing || t.Remaining != "" || len(t.ChildTags) > 0 {
		tc += ">"
		for _, tag := range t.ChildTags {
			tc += "<" + tag + ">"
		}
		return &token{Type: HTMLToken, Content: tc}
	} else {
		tc += " />"
		return &token{Type: HTMLToken, Content: t.Start() + " />"}
	}
}

func (t *tag) Start() string {
	tc := "<" + t.Name
	if len(t.Classes) > 0 {
		tc += " class=\"" + strings.Join(t.Classes, " ") + "\""
	}
	if t.Id != "" {
		tc += " id=\"" + t.Id + "\""
	}
	for n, v := range t.Attrs {
		tc += fmt.Sprintf(` %s="%s"`, n, v)
	}
	for n, v := range t.DynAttrs {
		tc += fmt.Sprintf(` %s="%s%s%s"`, n, LeftDelim, v, RightDelim)
	}

	return tc
}

func (t *tag) parseAttributes() error {
	for attrStartRegex.MatchString(strings.TrimSpace(t.Source)) {
		t.Source = strings.TrimSpace(t.Source)
		attr := attrStartRegex.FindStringSubmatch(t.Source)[1]
		t.Source = t.Source[len(attr)+1:]
		if len(t.Source) == 0 {
			return fmt.Errorf("Blank attribute value for attribute %s in %s", attr, t.Node.Code)
		}
		err := t.parseAttribute(attr)
		if err != nil {
			return err
		}
	}
	return nil
}
func (t *tag) parseAttribute(attr string) error {
	switch t.Source[0] {
	case '"':
		i := strings.Index(t.Source[1:], "\"")
		if i == -1 {
			return fmt.Errorf("Unclosed double quotes for attribute %s in %s", attr, t.Node.Code)
		}
		if i == 0 {
			return fmt.Errorf("Empty double quotes for attribute %s in %s", attr, t.Node.Code)
		}
		if strings.Contains(t.Source[1:1+i], LeftDelim) && !strings.Contains(t.Source[1:1+1], RightDelim) {
			di := strings.Index(t.Source[1+i:], RightDelim)
			if di != -1 {
				i += di + strings.Index(t.Source[1+i+di:], "\"")
			}
		}
		t.Attrs[attr] = t.Source[1 : 1+i]
		t.Source = t.Source[2+len(t.Attrs[attr]):]
	case '\'':
		i := strings.Index(t.Source[1:], "'")
		if i == -1 {
			return fmt.Errorf("Unclosed single quotes for attribute %s in %s", attr, t.Node.Code)
		}
		if i == 0 {
			return fmt.Errorf("Empty single quotes for attribute %s in %s", attr, t.Node.Code)
		}
		if strings.Contains(t.Source[1:1+i], LeftDelim) && !strings.Contains(t.Source[1:1+1], RightDelim) {
			di := strings.Index(t.Source[1+i:], RightDelim)
			if di != -1 {
				i += di + strings.Index(t.Source[1+i+di:], "'")
			}
		}
		t.Attrs[attr] = t.Source[1 : 1+i]
		t.Source = t.Source[2+len(t.Attrs[attr]):]
	case '(':
		i := strings.Index(t.Source[1:], ")")
		if i == -1 {
			return fmt.Errorf("Unclosed parentheses for attribute %s in %s", attr, t.Node.Code)
		}
		if i == 0 {
			return fmt.Errorf("Empty parentheses for attribute %s in %s", attr, t.Node.Code)
		}
		t.DynAttrs[attr] = t.Source[1 : 1+i]
		t.Source = t.Source[2+len(t.DynAttrs[attr]):]
	case '$':
		index := strings.Index(t.Source[1:], " ")
		if index == -1 {
			t.DynAttrs[attr] = t.Source
			t.Source = ""
		} else {
			t.DynAttrs[attr] = t.Source[:index]
			t.Source = t.Source[index:]
		}
	case '.':
		index := strings.Index(t.Source, " ")
		if index == -1 {
			t.DynAttrs[attr] = t.Source
			t.Source = ""
		} else {
			t.DynAttrs[attr] = t.Source[:index]
			t.Source = t.Source[index:]
		}
	default:
		index := strings.Index(t.Source[1:], " ")
		if index == -1 {
			t.DynAttrs[attr] = t.Source
			t.Source = ""
		} else {
			t.DynAttrs[attr] = t.Source[:index]
			t.Source = t.Source[index:]
		}
	}
	return nil
}

func (t *tag) Close() *token {
	if t.Enclosing || t.Remaining != "" || len(t.ChildTags) > 0 {
		tc := "</" + t.Name + ">"
		for _, tag := range t.ChildTags {
			tc = "</" + tag + ">" + tc
		}
		return &token{Type: HTMLToken, Content: tc}
	}
	return &token{Type: HTMLToken, Content: ""}
}
