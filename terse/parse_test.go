package terse

import (
	"bytes"
	"testing"

	"github.com/acsellers/multitemplate"
)

func TestParse(t *testing.T) {
	for _, test := range parseTests {
		tmpl := multitemplate.New("terse")
		var e error
		if len(test.Sources) == 0 {
			tmpl, e = tmpl.Parse("parse", test.Content, "terse")
			if e != nil {
				t.Logf("In test %s\n", test.Name)
				t.Error("Parse Error:", e)
				continue
			}
		} else {
			for tn, tc := range test.Sources {
				tmpl, e = tmpl.Parse(tn, tc, "terse")
				if e != nil {
					t.Logf("In test %s\n", test.Name)
					t.Error("Parse Error:", e)
					continue
				}
			}
		}
		b := &bytes.Buffer{}
		if test.Template == "" {
			e = tmpl.Execute(b, test.Data)
		} else {
			e = tmpl.ExecuteTemplate(b, test.Template, test.Data)
		}
		if e != nil {
			t.Logf("In test %s\n", test.Name)
			t.Error("Execute Error:", e)
		}
		if b.String() != test.Expected {
			t.Logf("In test %s\n", test.Name)
			t.Errorf("Result Error, Expected:`%s`\nReceived:`%s`", test.Expected, b.String())
		}
	}
}

var parseTests = []parseTest{
	parseTest{
		Name: "Blank Template",
	},
	parseTest{
		Name:     "Doctype Template Blank",
		Content:  "!!",
		Expected: "<!DOCTYPE html>",
	},
	parseTest{
		Name:     "Doctype Template Named",
		Content:  "!! Strict",
		Expected: `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">`,
	},
	parseTest{
		Name:     "Text in template",
		Content:  "blah blah",
		Expected: "blah blah",
	},
	parseTest{
		Name:     "Text in template",
		Content:  "blah blah\nnerr",
		Expected: "blah blah\nnerr",
	},
	parseTest{
		Name:     "Text in template",
		Content:  "bleh\n  wat",
		Expected: "bleh\n  wat",
	},
	parseTest{
		Name:     "Comment in template",
		Content:  "blah\n// don't show",
		Expected: "blah",
	},
	parseTest{
		Name:     "Nested Comment",
		Content:  "blah\n  // don't show",
		Expected: "blah",
	},
	parseTest{
		Name:     "Triple Nested Text",
		Content:  "first\n  second\n    third\n  fourth",
		Expected: "first\n  second\n    third\n  fourth",
	},
	parseTest{
		Name:     "Quadruple Nested Text",
		Content:  "First\n  Second\n    Third\n      Fourth",
		Expected: "First\n  Second\n    Third\n      Fourth",
	},
	parseTest{
		Name:     "If Statement",
		Content:  "?true\n  show",
		Expected: "show",
	},
	parseTest{
		Name:     "If/Else Statement (False)",
		Content:  "?false\n  no\n!?\n  yes",
		Expected: "yes",
	},
	parseTest{
		Name:     "If/Else Statement (True)",
		Content:  "?true\n  yes\n!?\n  no",
		Expected: "yes",
	},
	parseTest{
		Name:     "Range Statement (1 item)",
		Content:  "&.\n  wat",
		Expected: "\nwat",
		Data:     []string{"1"},
	},
	parseTest{
		Name:     "Range Statement (2 items)",
		Content:  "&.\n  wat",
		Expected: "\nwat\nwat",
		Data:     []string{"1", "2"},
	},
	parseTest{
		Name:     "Range Statement (0 items)",
		Content:  "&.\n  wat",
		Expected: "",
		Data:     []string{},
	},
	parseTest{
		Name:     "Range/Else Statement (0 items)",
		Content:  "&.\n  wat\n!&\n  no",
		Expected: "no",
		Data:     []string{},
	},
	parseTest{
		Name:     "Range/Else Statement (2 items)",
		Content:  "&.\n  wat\n!&\n  no",
		Expected: "\nwat\nwat",
		Data:     []string{"1", "2"},
	},
	parseTest{
		Name:     "Verbatim Statement",
		Content:  "/ $9@(#*$now",
		Expected: "$9@(#*$now",
	},
	parseTest{
		Name:     "Verbatim Statement with Nested lines",
		Content:  "/ $now\n  ?wat\n    do",
		Expected: "$now\n  ?wat\n    do",
	},
	parseTest{
		Name:     "Simple Exec",
		Content:  "= print 123",
		Expected: "123",
	},
	parseTest{
		Name:     "Continued Exec",
		Content:  "= print\n  /= 123\n  wat",
		Expected: "123\n  wat",
	},
}

type parseTest struct {
	Name     string
	Content  string
	Sources  map[string]string
	Expected string
	Template string
	Data     interface{}
}
