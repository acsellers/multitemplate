package helpers

import (
	"html/template"
)

var (
	jsTagFunctions = map[string]interface{}{
		"link_to_function": func(text, function string, options ...AttrList) template.HTML {
			al := combine("", "", options)
			al["onclick"] = template.HTML(function + "; return false")
			al["href"] = "#"
			return buildTag("a", template.HTML(template.HTMLEscapeString(text)), al)
			return ""
		},
	}
)
