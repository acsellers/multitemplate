package bham

import (
	"strings"
)

type templateLine struct {
	indentation int
	content     string
}

func (t templateLine) accept(chars string) bool {
	for _, c := range t.content {
		for _, s := range chars {
			if s == c {
				return true
			}
		}
		return false
	}
	return false
}
func (t templateLine) prefix(str string) bool {
	return len(t.content) >= len(str) && t.content[:len(str)] == str
}
func (t templateLine) suffix(str string) bool {
	return len(t.content) >= len(str) && t.content[len(t.content)-len(str):] == str
}
func (t templateLine) after(str string) templateLine {
	if t.accept(str) {
		return templateLine{
			t.indentation,
			strings.TrimSpace(t.content[1:]),
		}
	}

	return t
}

func (t templateLine) without(str string) templateLine {
	if t.prefix(str) {
		return templateLine{
			t.indentation,
			strings.TrimSpace(t.content[len(str):]),
		}
	}

	return t
}

func (t templateLine) String() string {
	return t.content
}

func (t templateLine) blockParameter() bool {
	action := t.after("=-")
	switch {
	case action.prefix("if "):
		return true
	case action.prefix("unless "):
		return true
	case action.prefix("else "):
		return true
	case action.prefix("with "):
		return true
	case action.prefix("range "):
		return true
	}
	return false
}

func (t templateLine) mightHaveElse() bool {
	action := t.after("=-")
	switch {
	case action.prefix("if "):
		return true
	case action.prefix("range "):
		return true
	}
	return false
}

func (t templateLine) isIf() bool {
	return t.after("=-").prefix("if ")
}
func (t templateLine) isUnless() bool {
	return t.after("=-").prefix("unless ")
}
func (t templateLine) isRange() bool {
	return t.after("=-").prefix("range ")
}
func (t templateLine) isWith() bool {
	return t.after("=-").prefix("with ")
}
func (t templateLine) isElse() bool {
	return t.after("=-").prefix("else")
}
func (t templateLine) isTag() bool {
	return t.accept("#.%")
}
func (t templateLine) needsContent() bool {
	if t.accept("=-") && t.after("=-").suffix("\\") {
		return true
	}
	if t.isTag() && t.suffix("\\") {
		return true
	}

	return false
}

func (t templateLine) appendContent(s string) string {
	if t.suffix("\\") {
		return strings.TrimSpace(t.content[:len(t.content)-1]) + " " + strings.TrimSpace(s)
	} else {
		return strings.TrimSpace(t.content) + " " + strings.TrimSpace(s)
	}
}
