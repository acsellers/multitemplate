package helpers

import (
	"fmt"
	"html/template"
)

var selectTagFunctions = template.FuncMap{
	"select_tag": func(name string, optionThing interface{}, options ...AttrList) template.HTML {
		al := combine(name, "", options)
		if _, ok := al["name"]; !ok {
			al["name"] = name
		}
		oc := optionFor(optionThing)
		return buildTag("select", oc.ToHTML(), al)
	},
	"option": func(text, value interface{}) Option {
		if tc, ok := text.(template.HTML); ok {
			return Option{Name: tc, Value: fmt.Sprint(value)}
		}
		return Option{
			Name:  template.HTML(template.HTMLEscapeString(fmt.Sprint(text))),
			Value: fmt.Sprint(value),
		}
	},
	"options": func(vals ...interface{}) OptionList {
		ol := OptionList{}
		for _, val := range vals {
			if vc, ok := val.(template.HTML); ok {
				ol = append(ol, Option{
					Name:  vc,
					Value: fmt.Sprint(val),
				})
			} else {
				vs := fmt.Sprint(val)
				ol = append(ol, Option{
					Name:  template.HTML(template.HTMLEscapeString(vs)),
					Value: vs,
				})
			}
		}
		return ol
	},
	"options_with_values": func(vals ...interface{}) OptionList {
		ol := OptionList{}
		var waiting template.HTML
		for _, val := range vals {
			switch tav := val.(type) {
			case template.HTML:
				if waiting != "" {
					ol = append(ol, Option{Name: waiting, Value: string(tav)})
					waiting = ""
				} else {
					waiting = tav
				}
			case string:
				if waiting != "" {
					ol = append(ol, Option{Name: waiting, Value: string(tav)})
					waiting = ""
				} else {
					waiting = template.HTML(template.HTMLEscapeString(tav))
				}
			case Option:
				ol = append(ol, tav)
			}
		}
		return ol
	},
}

type Option struct {
	Name  template.HTML
	Value string
}

func (o Option) ToHTML() template.HTML {
	content := "<option value=\"" + template.HTMLEscapeString(o.Value) + "\">"
	content += string(o.Name)
	content += "</option>"
	return template.HTML(content)
}

func (o Option) Options() []Option {
	return []Option{o}
}

type OptionLike interface {
	ToHTML() template.HTML
	Options() []Option
}
type OptionList []OptionLike

func (ol OptionList) ToHTML() template.HTML {
	var content template.HTML
	for _, o := range ol {
		content += o.ToHTML()
	}
	return content
}

type OptionGroup struct {
	Label        string
	InnerOptions []Option
}

func (og OptionGroup) ToHTML() template.HTML {
	content := template.HTML("<optgroup label=\"" + template.HTMLEscapeString(og.Label) + "\">")
	for _, opt := range og.InnerOptions {
		content += opt.ToHTML()
	}
	content += "</optgroup>"
	return content
}

func (og OptionGroup) Options() []Option {
	return og.InnerOptions
}

func optionFor(ot interface{}) OptionList {
	if ol, ok := ot.(OptionList); ok {
		return ol
	}

	return OptionList{}
}

func init() {
	for k, v := range selectTagFunctions {
		formTagFunctions[k] = v
	}
}