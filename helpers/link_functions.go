package helpers

import (
	"fmt"
	"html/template"
	"net/url"
)

var linkFuncs = template.FuncMap{
	"link_to": func(link, name string, options ...AttrList) template.HTML {
		u, e := url.Parse(link)
		if e != nil {
			return ""
		}

		al := combine("", "", options)
		al["href"] = u.String()
		return buildTag("a", template.HTML(template.HTMLEscapeString(name)), al)
	},
	"link_to_function": func(text, function string, options ...AttrList) template.HTML {
		al := combine("", "", options)
		al["onclick"] = template.HTML(function + "; return false")
		al["href"] = "#"
		return buildTag("a", template.HTML(template.HTMLEscapeString(text)), al)
		return ""
	},
	"url_encode": url.QueryEscape,
	"urlize": func(link string) template.HTML {
		u, e := url.Parse(link)
		if e != nil {
			return ""
		}

		return template.HTML(
			fmt.Sprintf(`<a href="%s">%s</a>`, u.String(), link),
		)
	},
	"urlize_truncate": func(link string, num int) template.HTML {
		chars := []rune(link)
		u, e := url.Parse(link)
		if e != nil {
			return ""
		}
		if num >= len(chars) {

			return template.HTML(
				fmt.Sprintf(`<a href="%s">%s</a>`, u.String(), link),
			)
		}
		return template.HTML(
			fmt.Sprintf(`<a href="%s">%s</a>`, u.String(), chars[:num]),
		)
	},
}
