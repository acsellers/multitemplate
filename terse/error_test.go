package terse

import (
	"strings"
	"testing"

	"github.com/acsellers/multitemplate"
	"github.com/acsellers/multitemplate/helpers"
)

func TestErrors(t *testing.T) {
	helpers.LoadHelpers("all")
	for _, test := range errorTests {
		tmpl := multitemplate.New("terse")
		var e error
		tmpl, e = tmpl.Parse("parse", test.Source, "terse")
		if e == nil {
			t.Error("No Error:", test.Name)
			continue
		} else if test.Contains != "" {
			if !strings.Contains(e.Error(), test.Contains) {
				t.Error("Incorrect error for", test.Name, "Expected:", test.Contains, "Was:", e)
			}
		}
	}
}

var errorTests = []errorTest{
	errorTest{
		Name:     "Unclosed attribute parentheses",
		Source:   "html name=(link_to",
		Contains: "Unclosed parentheses",
	},
	errorTest{
		Name:     "Empty attribute parentheses",
		Source:   "html name=()",
		Contains: "Empty parentheses",
	},
	errorTest{
		Name:     "Unclosed attribute single quotes ",
		Source:   "html name='link_to",
		Contains: "Unclosed single quotes",
	},
	errorTest{
		Name:     "Empty attribute single quotes",
		Source:   "html name=''",
		Contains: "Empty single quotes",
	},
	errorTest{
		Name:     "Unclosed attribute double quotes ",
		Source:   "html name=\"link_to",
		Contains: "Unclosed double quotes",
	},
	errorTest{
		Name:     "Empty attribute double quotes",
		Source:   "html name=\"\"",
		Contains: "Empty double quotes",
	},
	errorTest{
		Name:     "Blank attribute variable",
		Source:   "html name=",
		Contains: "Blank attribute",
	},
	errorTest{
		Name:     "Headless Totem Pole",
		Source:   "html > head >",
		Contains: "Missing tag",
	},
	errorTest{
		Name:   "Malformed code test",
		Source: "= link_to (",
	},
	errorTest{
		Name:     "Multiple errors test",
		Source:   "head >\n= link_to (",
		Contains: "Missing tag",
	},
	errorTest{
		Name:   "Bad function name for if",
		Source: "?asdfasdfsafd\n  no",
	},
	errorTest{
		Name:   "Bad function name for if",
		Source: "&asdfasdfsafd\n  no",
	},
	errorTest{
		Name:   "Bad interpolated code 1",
		Source: "asdf {{ asdfasdf }}",
	},
	errorTest{
		Name:   "Bad interpolated code 2",
		Source: "asdf\n  {{ asdfasdf }}",
	},
	errorTest{
		Name:   "Bad interpolated code 3",
		Source: "asdf\n  {{ {}{}} }}",
	},
	errorTest{
		Name:   "Bad interpolated code 3",
		Source: "asdf}} {{  }}",
	},
	errorTest{
		Name:   "Weird if statement",
		Source: "?!?",
	},
}

type errorTest struct {
	Name     string
	Source   string
	Contains string
}
