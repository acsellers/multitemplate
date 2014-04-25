package helpers

import (
	"fmt"
	"html/template"
)

var selectTagFuncs = template.FuncMap{
	"select_tag": func(name string, optionThing interface{}, options ...AttrList) template.HTML {
		al := combine(name, "", options)
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
			switch av := val.(type) {
			case template.HTML:
				ol = append(ol, Option{
					Name:  av,
					Value: template.JSEscapeString(string(av)),
				})
			case string:
				ol = append(ol, Option{
					Name:  template.HTML(template.HTMLEscapeString(av)),
					Value: template.JSEscapeString(av),
				})
			case Option:
				ol = append(ol, av)
			case OptionLike:
				for _, o := range av.Options() {
					ol = append(ol, o)
				}
			default:
				vs := fmt.Sprint(val)
				ol = append(ol, Option{
					Name:  template.HTML(template.HTMLEscapeString(vs)),
					Value: template.JSEscapeString(vs),
				})
			}
		}
		return ol
	},
	"group_options": func(label string, ol ...OptionLike) OptionGroup {
		options := []Option{}
		for _, oi := range ol {
			options = append(options, oi.Options()...)
		}
		return OptionGroup{label, options}
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

func (og OptionList) Options() []Option {
	options := []Option{}
	for _, oi := range og {
		options = append(options, oi.Options()...)
	}
	return options
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
