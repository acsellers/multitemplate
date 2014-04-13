package helpers

import (
	"html/template"
	"net/url"
	"testing"
	. "github.com/acsellers/assert"
)

func TestAsset(t *testing.T) {
	Within(t, func(test *Test) {
		u, _ := url.Parse("http://localhost/")
		AppInfo = AssetInfo{
			RootURL:            u,
			ImageRelative:      "img",
			JavascriptRelative: "js",
			StylesheetRelative: "css",
		}
		for _, assetTest := range assetTests {
			test.Section("Helper: " + assetTest.Helper)
			if f, ok := assetFuncs[assetTest.Helper]; ok {
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

type helperTest struct {
	Helper   string
	Args     []string
	Attrs    []AttrList
	Expected string
}

var assetTests = []helperTest{
	helperTest{
		Helper:   "atom_link",
		Args:     []string{"/api/atom"},
		Expected: `<link rel="alternate" type="application/atom+xml" title="ATOM" href="http://localhost/api/atom" />`,
	},
	helperTest{
		Helper:   "favicon_link",
		Args:     []string{"favicon.png"},
		Expected: `<link rel="shortcut icon" type="image/png" href="/img/favicon.png" />`,
	},
	helperTest{
		Helper:   "favicon_link",
		Args:     []string{"favicon.ico"},
		Expected: `<link rel="shortcut icon" type="image/vnd.microsoft.icon" href="/img/favicon.ico" />`,
	},
	helperTest{
		Helper:   "image_tag",
		Args:     []string{"window.png"},
		Expected: `<img src="/img/window.png" />`,
	},
	helperTest{
		Helper:   "image_tag",
		Args:     []string{"http://placekitten.com/200/300"},
		Expected: `<img src="http://placekitten.com/200/300" />`,
	},
	helperTest{
		Helper:   "image_tag",
		Args:     []string{"subdir/first.png"},
		Expected: `<img src="/img/subdir/first.png" />`,
	},
	helperTest{
		Helper:   "image_tag",
		Args:     []string{"window.png"},
		Attrs:    []AttrList{AttrList{"alt": "Blah"}},
		Expected: `<img alt="Blah" src="/img/window.png" />`,
	},
	helperTest{
		Helper:   "image_tag",
		Args:     []string{"window.png"},
		Attrs:    []AttrList{AttrList{"width": 16, "height": 16}},
		Expected: `<img width="16" height="16" src="/img/window.png" />`,
	},
	helperTest{
		Helper:   "javascript_link",
		Args:     []string{"jquery"},
		Expected: `<script src="/js/jquery.js"> </script>`,
	},
	helperTest{
		Helper:   "javascript_link",
		Args:     []string{"bootstrap.js", "jquery.min"},
		Expected: `<script src="/js/bootstrap.js"> </script><script src="/js/jquery.min.js"> </script>`,
	},
	helperTest{
		Helper:   "javascript_link",
		Args:     []string{"http://ajaxcdn.net/jquery/1.2.3/jquery.min.js"},
		Expected: `<script src="http://ajaxcdn.net/jquery/1.2.3/jquery.min.js"> </script>`,
	},
	helperTest{
		Helper:   "rss_link",
		Args:     []string{"/api/rss"},
		Expected: `<link rel="alternate" type="application/rss+xml" title="RSS" href="http://localhost/api/rss" />`,
	},
	helperTest{
		Helper:   "rss_link",
		Args:     []string{"http://rssforeveryone.com/349852.xml"},
		Expected: `<link rel="alternate" type="application/rss+xml" title="RSS" href="http://rssforeveryone.com/349852.xml" />`,
	},
	helperTest{
		Helper:   "stylesheet_link",
		Args:     []string{"bootstrap"},
		Expected: `<link media="screen" rel="stylesheet" href="/css/bootstrap.css" />`,
	},
}
