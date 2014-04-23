package terse

import "strings"

func parseTag(tagCode string, children bool) (*token, string, *token) {
	t := &tag{Name: "div", Source: tagCode}
	t.Parse(children)
	return t.Open(), t.Remaining, t.Close()
}

type tag struct {
	Source    string
	Name      string
	Id        string
	Classes   []string
	Attrs     map[string]string
	Remaining string
	Enclosing bool
}

func (t *tag) Parse(children bool) {
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
		if t.Source[0] == '#' {
			t.Id = firstTextToken(t.Source[1:])
			t.Source = t.Source[len(t.Id)+1:]
		}
		if t.Source[0] == '.' {
			cl := firstTextToken(t.Source[1:])
			t.Source = t.Source[len(cl)+1:]
			t.Classes = append(t.Classes, cl)
		}
	}
	if len(t.Source) > 0 && t.Source[0] != '(' {
	}
	t.Remaining = t.Source
}

func (t *tag) Open() *token {
	if t.Enclosing || t.Remaining != "" {
		return &token{Type: HTMLToken, Content: t.Start() + ">"}
	} else {
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
	return tc
}

func (t *tag) Close() *token {
	if t.Enclosing || t.Remaining != "" {
		return &token{Type: HTMLToken, Content: "</" + t.Name + ">"}
	}
	return &token{Type: HTMLToken, Content: ""}
}
