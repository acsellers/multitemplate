package terse

import (
	"fmt"
	ht "html/template"
	"strings"
	"text/template"
	"text/template/parse"
)

type resources struct {
	funcs ht.FuncMap
	vars  []string
}

func (rsc *resources) Prelude() string {
	ps := ""
	for _, v := range rsc.vars {
		ps += "{{ " + v + " := $ }}"
	}
	return ps
}

func surround(s string) string {
	if strings.HasPrefix(strings.TrimSpace(s), "{{") {
		return s
	}
	return "{{ " + s + " }}"
}
func actionNode(code string, rsc *resources) (*parse.ActionNode, error) {
	t := template.New("mule").Funcs(template.FuncMap(rsc.funcs))
	t, err := t.Parse(rsc.Prelude() + surround(code))
	if err != nil {
		return nil, err
	}

	ln := t.Tree.Root.Nodes[len(t.Tree.Root.Nodes)-1]
	if an, ok := ln.(*parse.ActionNode); ok {
		if len(an.Pipe.Decl) > 0 {
		DeclLoop:
			for _, v := range an.Pipe.Decl {
				vs := v.String()
				for _, vd := range rsc.vars {
					if vs == vd {
						continue DeclLoop
					}
				}
				rsc.vars = append(rsc.vars, vs)
			}
		}
		return an, nil
	}
	return nil, fmt.Errorf("Node could not be parsed")
}
