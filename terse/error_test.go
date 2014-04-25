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
		Name:     "Unclosed attribute parentheses",
		Source:   "html name=()",
		Contains: "Empty parentheses",
	},
}

type errorTest struct {
	Name     string
	Source   string
	Contains string
}
