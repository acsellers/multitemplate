package terse

import (
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
		Contains: "Empty attribute",
	},
	errorTest{
		Name:     "Headless Totem Pole",
		Source:   "html > head >",
		Contains: "Missing tag",
	},
}

type errorTest struct {
	Name     string
	Source   string
	Contains string
}
