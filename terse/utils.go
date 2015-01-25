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
	tt    *tokenTree
	vars  []string
	err   error
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
		if an.ElseList != nil {
			for _, ln := range an.ElseList.Nodes {
				rsc.UpdateVars(ln)
			}
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
func actionNode(code string, rsc *resources, pos int) (*parse.ActionNode, error) {
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
		normalize(an, pos)
		return an, nil
	}
	return nil, fmt.Errorf("Node could not be parsed: %s", code)
}

func normalize(an *parse.ActionNode, pos int) {
	drop := int(an.Pos)
	an.Pos = parse.Pos(pos)
	an.Pipe.Pos = parse.Pos(int(an.Pipe.Pos) - drop + pos)
	for _, decl := range an.Pipe.Decl {
		decl.Pos = parse.Pos(int(decl.Pos) - drop + pos)
	}
	for _, cmd := range an.Pipe.Cmds {
		cmd.Pos = parse.Pos(int(cmd.Pos) - drop + pos)
	}
}

func textNodes(text string, rsc *resources, pos int) []parse.Node {
	t := template.New("mule").Funcs(template.FuncMap(rsc.funcs)).Delims(LeftDelim, RightDelim)
	t, err := t.Parse(rsc.Prelude() + text)
	if err != nil {
		rsc.err = err
		return nil
	}

	ln := t.Tree.Root.Nodes[len(rsc.vars):]
	drop := int(ln[0].Position())
	for _, n := range ln {
		setPos(n, int(n.Position())-drop+pos)
		rsc.UpdateVars(n)
	}
	return ln
}

func setPos(n parse.Node, pos int) {
	switch an := n.(type) {
	case *parse.ActionNode:
		an.Pos = parse.Pos(pos)
	case *parse.TextNode:
		an.Pos = parse.Pos(pos)
	case *parse.TemplateNode:
		an.Pos = parse.Pos(pos)
	}
}
