package helpers

import (
	"fmt"
	"html/template"
)

var coreFunctions = template.FuncMap{
	"attr": func(name string, value interface{}) Attr {
		return Attr{name, value}
	},
	"attrs": makeAttrList,

	"data": func(args ...interface{}) (AttrList, error) {
		al, e := makeAttrList(args...)
		if e != nil {
			return AttrList{}, e
		}
		dal := AttrList{}
		for k, v := range al {
			dal["data-"+k] = v
		}
		return dal, nil
	},
}

type AttrList map[string]interface{}

type Attr struct {
	Name  string
	Value interface{}
}

func makeAttrList(args ...interface{}) (AttrList, error) {
	al := AttrList{}
	var name string
	for _, arg := range args {
		switch v := arg.(type) {
		case Attr:
			al[v.Name] = v.Value
			name = ""
		case string:
			if name == "" {
				name = v
				al[name] = true
			} else {
				al[name] = v
				name = ""
			}
		default:
			if name == "" {
				name = fmt.Sprint(v)
				al[name] = true
			} else {
				al[name] = v
				name = ""
			}
		}
	}
	return al, nil
}
