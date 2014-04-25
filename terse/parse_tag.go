package terse

import (
	"fmt"
	"regexp"
	"strings"
)

func parseTag(tagCode string, children bool) (*tag, error) {
	t := &tag{
		Name:     "div",
		Source:   tagCode,
		Attrs:    make(map[string]string),
		DynAttrs: make(map[string]string),
	}
	e := t.Parse(children)
	return t, e
}

type tag struct {
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

var attrStartRegex = regexp.MustCompile(`^([a-zA-Z0-9]+)=`)

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
		fmt.Println("Multi line")
	} else if attrStartRegex.MatchString(strings.TrimSpace(t.Source)) {
		for attrStartRegex.MatchString(strings.TrimSpace(t.Source)) {
			t.Source = strings.TrimSpace(t.Source)
			attr := attrStartRegex.FindStringSubmatch(t.Source)[1]
			t.Source = t.Source[len(attr)+1:]
			switch t.Source[0] {
			case '"':
				t.Attrs[attr] = t.Source[1 : 1+strings.Index(t.Source[1:], "\"")]
				t.Source = t.Source[2+len(t.Attrs[attr]):]
			case '\'':
				t.Attrs[attr] = t.Source[1 : 1+strings.Index(t.Source[1:], "'")]
				t.Source = t.Source[2+len(t.Attrs[attr]):]
			case '(':
				t.DynAttrs[attr] = t.Source[1 : 1+strings.Index(t.Source[1:], ")")]
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
		}
	}
	t.Remaining = t.Source
	if len(t.Remaining) > 0 && t.Remaining[0] == ' ' {
		t.Remaining = t.Remaining[1:]
	}
	for len(t.Remaining) > 0 && t.Remaining[0] == '>' {
		check := strings.TrimSpace(t.Remaining[1:])
		switch {
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
