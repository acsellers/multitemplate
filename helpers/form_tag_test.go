package helpers

/*

import (
	"html/template"
	"net/url"
	"testing"
	. "github.com/acsellers/assert"
)

func TestFormTag(t *testing.T) {
	Within(t, func(test *Test) {
		for _, formTagTest := range assetTests {
			test.Section("Helper: " + assetTest.Helper)
			if f, ok := formTagFuncs[assetTest.Helper]; ok {
				var r string
				switch af := f.(type) {
				case func(...string) template.HTML:
					r = string(af(assetTest.Args...))
				case func(string) template.HTML:
					r = string(af(assetTest.Args[0]))
				case func(string) string:
					r = af(assetTest.Args[0])
				case func(string, ...AttrList) template.HTML:
					r = string(af(assetTest.Args[0], assetTest.Attrs...))
				default:
					t.Fatalf("Function %s was not type asserted", assetTest.Helper)
				}
				test.AreEqual(r, assetTest.Expected)
			} else {
				t.Errorf("Could not find function %s", assetTest.Helper)
			}
		}
	})
}

var formTagTests = []helperTest{
	helperTest{
		Helper:   "atom_link",
		Args:     []string{"/api/atom"},
		Expected: `<link rel="alternate" type="application/atom+xml" title="ATOM" href="http://localhost/api/atom" />`,
	},
	helperTest{
		Helper: "favicon_link",
		Args:   []string{"favicon.png"},
		Expected: string(buildTag("link", "", AttrList{
			"rel":  "shortcut icon",
			"type": "image/png",
			"href": "/img/favicon.png",
		})),
	},
}
*/
