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
		ps += LeftDelim + v + " := $ " + RightDelim
	}
	return ps
}

func (rsc *resources) UpdateVars(n parse.Node) {
	switch an := n.(type) {
	case *parse.ActionNode:
		for _, decl := range an.Pipe.Decl {
			rsc.vars = append(rsc.vars, decl.String())
		}
	case *parse.IfNode:
		for _, ln := range an.List.Nodes {
			rsc.UpdateVars(ln)
		}
		for _, ln := range an.ElseList.Nodes {
			rsc.UpdateVars(ln)
		}
	case *parse.RangeNode:
		for _, ln := range an.List.Nodes {
			rsc.UpdateVars(ln)
		}
		for _, ln := range an.ElseList.Nodes {
			rsc.UpdateVars(ln)
		}
	case *parse.WithNode:
		for _, ln := range an.List.Nodes {
			rsc.UpdateVars(ln)
		}
		for _, ln := range an.ElseList.Nodes {
			rsc.UpdateVars(ln)
		}
	}
}

func surround(s string) string {
	if strings.HasPrefix(strings.TrimSpace(s), LeftDelim) {
		return s
	}
	return LeftDelim + s + RightDelim
}
func actionNode(code string, rsc *resources) (*parse.ActionNode, error) {
	t := template.New("mule").Funcs(template.FuncMap(rsc.funcs)).Delims(LeftDelim, RightDelim)
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

func textNodes(text string, rsc *resources) []parse.Node {
	t := template.New("mule").Funcs(template.FuncMap(rsc.funcs)).Delims(LeftDelim, RightDelim)
	t, err := t.Parse(rsc.Prelude() + text)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	ln := t.Tree.Root.Nodes[len(rsc.vars):]
	for _, n := range ln {
		rsc.UpdateVars(n)
	}
	return ln
}
