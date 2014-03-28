package multitemplate

import "html/template"

// Generate the Context dependent functions for a template execution
func GenerateFuncs(t *Template) template.FuncMap {
	return template.FuncMap{
		"yield": func(vals ...interface{}) (template.HTML, error) {
			switch len(vals) {
			case 0:
				return t.ctx.mainContent, nil
			case 1:
				if name, ok := vals[0].(string); ok {
					if t.ctx.Yields[name] != "" {
						return t.ctx.exec(t.ctx.Yields[name], t.ctx.Dot)
					}
					if t.ctx.Content[name] != "" {
						return t.ctx.Content[name], nil
					}
				}
				return t.ctx.exec(t.ctx.Main, vals[0])
			case 2:
				if name, ok := vals[0].(string); ok {
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
				}
			default:
				if name, ok := vals[0].(string); ok {
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
			}
			return template.HTML(""), nil
		},
		"content_for": func(name string, templateName string) string {
			t.ctx.Yields[name] = templateName
			return ""
		},
		"root_dot": func() interface{} {
			return t.ctx.Dot
		},
		"exec": func(templateName string, dot interface{}) (template.HTML, error) {
			return t.ctx.exec(templateName, dot)
		},
		"block": func(name string) (string, error) {
			if t.ctx.parent != "" {
				t.ctx.output.Open(name)
			} else {
				if c, ok := t.ctx.Content[name]; ok {
					t.ctx.output.Write([]byte(c))
					t.ctx.output.Nop()
				}
				if _, ok := t.ctx.Yields[name]; ok {
					c, e := t.ctx.exec(t.ctx.Yields[name], t.ctx.Dot)
					t.ctx.output.Write([]byte(c))
					t.ctx.output.Nop()
					return "", e
				}
			}
			return "", nil
		},
		"end_block": func() string {
			n, c := t.ctx.output.Close()
			if n == "" {
				return ""
			}
			if _, ok := t.ctx.Content[n]; !ok {
				t.ctx.Content[n] = template.HTML(c)
			}
			return ""
		},
		"extends": func(parent string) string {
			t.ctx.output.NoRoot()
			t.ctx.parent = parent
			return ""
		},
	}
}

// Functions that are not tied to a context, but are part of the core
// multitemplate system
var StaticFuncs = template.FuncMap{
	"fallback": func(s string) Fallback {
		return Fallback(s)
	},
}
