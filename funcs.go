package multitemplate

import (
	"bytes"
	"html/template"
)

func GenerateFuncs(t *Template) template.FuncMap {
	return template.FuncMap{
		"yield": func(name string, vals ...interface{}) (template.HTML, error) {
			switch len(vals) {
			case 0:
				if t.ctx.Yields[name] != "" {
					b := bytes.Buffer{}
					e := t.ExecuteTemplate(&b, t.ctx.Yields[name], t.ctx.Dot)
					return template.HTML(b.String()), e
				}
			case 1:
				// Use provided fallback if necessary
				if f, ok := vals[0].(Fallback); ok {
					return t.ctx.execWithFallback(name, f, t.ctx.Dot)
					// Provided data to run
				} else {
					if t.ctx.Yields[name] != "" {
						return t.ctx.exec(t.ctx.Yields[name], vals[0])
					}
					if t.ctx.Content[name] != "" {
						return t.ctx.Content[name], nil
					}
					return template.HTML(""), nil
				}
			default:
				var d interface{}
				for _, v := range vals {
					if f, ok := v.(Fallback); ok {
						if d == nil {
							t.ctx.execWithFallback(name, f, vals[1])
						} else {
							t.ctx.execWithFallback(name, f, d)
						}
					} else {
						if d == nil {
							d = v
						}
					}
				}
				return t.ctx.exec(name, d)
			}
			return template.HTML(""), nil
		},
		"content_for": func(name string, templateName string) {
			t.ctx.Yields[name] = templateName
		},
		"root_dot": func() interface{} {
			return t.ctx.Dot
		},
		"exec": func(templateName string, dot interface{}) (template.HTML, error) {
			return t.ctx.exec(templateName, dot)
		},
	}
}

var StaticFuncs = template.FuncMap{
	"fallback": func(s string) Fallback {
		return Fallback(s)
	},
}
