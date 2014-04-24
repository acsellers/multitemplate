package helpers

import (
	"html/template"
	"testing"
	. "github.com/acsellers/assert"
)

func TestFormTag(t *testing.T) {
	Within(t, func(test *Test) {
		for _, formTagTest := range formTagTests {
			test.Section("Helper: " + formTagTest.Helper)
			if f, ok := formTagFuncs[formTagTest.Helper]; ok {
				var r string
				switch af := f.(type) {
				case func(string, interface{}, ...AttrList) template.HTML:
					r = string(af(formTagTest.Args[0], formTagTest.Args[1]))
				case func(string, ...AttrList) template.HTML:
					r = string(af(formTagTest.Args[0], formTagTest.Attrs...))
				default:
					t.Fatalf("Function %s was not type asserted", formTagTest.Helper)
				}
				test.AreEqual(r, formTagTest.Expected)
			} else {
				t.Errorf("Could not find function %s", formTagTest.Helper)
			}
		}
	})
}

var formTagTests = []helperTest{
	helperTest{
		Helper:   "button_tag",
		Args:     []string{"Hello"},
		Expected: `<button name="hello">Hello</button>`,
	},
	helperTest{
		Helper:   "check_box_tag",
		Args:     []string{"is_active", "true"},
		Expected: `<input name="is_active" id="is_active" value="true" type="checkbox" />`,
	},
}

/*
type helperTest struct {
	Helper   string
	Args     []string
	Attrs    []AttrList
	Expected string
}
*/
