package multitemplate

import "html/template"

func generateFuncs(t *Template) template.FuncMap {
	return template.FuncMap{
		"yield": func(vals ...interface{}) (string, error) {
			var e error
			switch len(vals) {
			case 0:
				t.ctx.output.Immediate(t.ctx.mainContent)
				return "<\"'.", nil
			case 1:
				if name, ok := vals[0].(string); ok {
					if t.ctx.Yields[name] != "" {
						rb, e := t.ctx.exec(t.ctx.Yields[name], t.ctx.Dot)
						t.ctx.output.Immediate(rb)
						if e != nil {
							return "<\"'.", e
						}
					}
					if rb, ok := t.ctx.Blocks[name]; ok {
						t.ctx.output.Immediate(rb)
						return "<\"'.", nil
					}
				}
				rb, e := t.ctx.exec(t.ctx.Main, vals[0])
				t.ctx.output.Immediate(rb)
				return "<\"'.", e
			case 2:
				if name, ok := vals[0].(string); ok {
					// Use provided fallback if necessary
					if f, ok := vals[0].(fallback); ok {
						t.ctx.output.next, e = t.ctx.execWithFallback(name, f, t.ctx.Dot)
						return "<\"'.", e
						// Provided data to run
					} else {
						if t.ctx.Yields[name] != "" {
							t.ctx.output.next, e = t.ctx.exec(t.ctx.Yields[name], vals[0])
							return "<\"'.", e
						}
						if rb, ok := t.ctx.Blocks[name]; ok {
							t.ctx.output.next = rb
							return "<\"'.", nil
						}
						return "", nil
					}
				}
			default:
				if name, ok := vals[0].(string); ok {
					var d interface{}
					for _, v := range vals {
						if f, ok := v.(fallback); ok {
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
					t.ctx.output.next, e = t.ctx.exec(name, d)
					return "<\"'.", e
				}
			}
			return "", nil
		},
		"content_for": func(name string, templateName string) string {
			if t.ctx.Yields[name] == "" {
				if _, ok := t.ctx.Blocks[name]; !ok {
					t.ctx.Yields[name] = templateName
				}
			}
			return ""
		},
		"root_dot": func() interface{} {
			return t.ctx.Dot
		},
		"exec": func(templateName string, dot interface{}) (string, error) {
			rb, e := t.ctx.exec(templateName, dot)
			t.ctx.output.next = rb
			return "<\"'.", e
		},
		"block": func(name string) (string, error) {
			if t.ctx.openableScope() {
				t.ctx.output.Open(name)
			} else {
				if _, ok := t.ctx.Yields[name]; ok {
					rb, e := t.ctx.exec(t.ctx.Yields[name], t.ctx.Dot)
					t.ctx.output.Nop(rb)
					return "", e
				} else if rb, ok := t.ctx.Blocks[name]; ok {
					t.ctx.output.Nop(rb)
				} else {
					return "", nil
				}
			}
			return "<\"'.", nil
		},
		"exec_block": func(name string) (string, error) {
			if _, ok := t.ctx.Yields[name]; ok {
				rb, e := t.ctx.exec(t.ctx.Yields[name], t.ctx.Dot)
				t.ctx.output.Nop(rb)
				return "<\"'.", e
			} else if rb, ok := t.ctx.Blocks[name]; ok {
				t.ctx.output.Nop(rb)
			}
			return "<\"'.", nil
		},
		"define_block": func(name string) string {
			t.ctx.output.Open(name)
			return "<\"'."
		},
		"end_block": func() string {
			n, rb := t.ctx.output.Close()
			if n == "" {
				return ""
			}
			if _, ok := t.ctx.Blocks[n]; !ok {
				if t.ctx.Yields[n] == "" {
					t.ctx.Blocks[n] = rb
				}
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
	"fallback": func(s string) fallback {
		return fallback(s)
	},
}

// LoadedFuncs is the place to load functions to be loaded.
var LoadedFuncs = template.FuncMap{}
