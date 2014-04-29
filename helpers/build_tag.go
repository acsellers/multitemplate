package helpers

import (
	"fmt"
	"html/template"
	"strings"
	"unicode"
)

func buildTag(tagName string, content template.HTML, al AttrList) template.HTML {
	tag := "<" + tagName

	for name, value := range al {
		if !controlAttribute(name) {
			tag += " " + strings.Map(attrNameFilter, name) + "=\""
			if b, ok := value.(bool); ok && b {
				tag += strings.Map(attrNameFilter, name) + "\""
			} else if hs, ok := value.(template.HTML); ok {
				tag += string(hs) + "\""
			} else {
				tag += template.JSEscapeString(fmt.Sprint(value)) + "\""
			}
		}
	}

	if v, ok := al["MT_skip_close"]; ok {
		if b, ok := v.(bool); ok && b {
			return template.HTML(tag + ">" + string(content))
		}
	}

	if string(content) == "" {
		tag += " />"
	} else {
		tag += ">" + string(content) + "</" + tagName + ">"
	}
	return template.HTML(tag)
}

func controlAttribute(attr string) bool {
	switch attr {
	case "MT_skip_close":
		return true
	default:
		return false
	}
}

func attrNameFilter(r rune) rune {
	switch {
	case r >= 'A' && r <= 'Z':
		return r
	case r >= 'a' && r <= 'z':
		return r
	case r == '-':
		return r
	case r == '_':
		return r
	}
	return -1
}

func nameFilter(r rune) rune {
	switch {
	case unicode.IsUpper(r):
		return unicode.ToLower(r)
	case unicode.IsPunct(r) || unicode.IsSymbol(r):
		return '_'
	case unicode.IsSpace(r):
		return '_'
	case unicode.IsGraphic(r):
		return r
	default:
		return -1
	}
}
