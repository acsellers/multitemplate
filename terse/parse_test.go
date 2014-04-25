package terse

import (
	"bytes"
	"html/template"
	"testing"

	"github.com/acsellers/multitemplate"
)

func TestParse(t *testing.T) {
	for _, test := range parseTests {
		tmpl := multitemplate.New("terse").Funcs(test.Funcs)
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
		Name:     "Range Statement",
		Content:  "&.:$element:$index\n  li= $element",
		Expected: "\n<li>a\n</li>\n<li>i\n</li>\n<li>b\n</li>",
		Data:     []string{"a", "i", "b"},
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
	parseTest{
		Name:     "Auto-Close Exec",
		Content:  "= div \n  blah",
		Expected: "<div>\n  blah\n</div>",
		Funcs: template.FuncMap{
			"div": func() template.HTML {
				return "<div>"
			},
			"end_div": func() template.HTML {
				return "</div>"
			},
		},
	},
	parseTest{
		Name: "Auto-Close Exec w/ Args",
		Content: `= tag "div"
  blah`,
		Expected: "<div>\n  blah\n</div>",
		Funcs: template.FuncMap{
			"tag": func(s string) template.HTML {
				return "<" + template.HTML(s) + ">"
			},
			"end_tag": func(s string) template.HTML {
				return "</" + template.HTML(s) + ">"
			},
		},
	},
	parseTest{
		Name: "Child Template",
		Sources: map[string]string{
			"main":  "[block]\n  1234",
			"child": "child\n  4321",
		},
		Template: "child",
		Expected: "child\n  4321",
	},
	parseTest{
		Name: "Block",
		Sources: map[string]string{
			"main":  "[block]\n  12345",
			"child": "child\n  54321",
		},
		Template: "main",
		Expected: "12345",
	},
	parseTest{
		Name: "Inherited Block",
		Sources: map[string]string{
			"main":  "[block]\n  12345",
			"child": "@@main\n[block]\n  54321",
		},
		Template: "child",
		Expected: "\n54321",
	},
	parseTest{
		Name: "Define Block",
		Sources: map[string]string{
			"main":  "[block]\n  12345\n678",
			"child": "@@main\n^block]\n  54321",
		},
		Template: "child",
		Expected: "\n54321\n678",
	},
	parseTest{
		Name: "Yield Block",
		Sources: map[string]string{
			"main":  "@block\n123",
			"child": "@@main\n^block]\n  54321",
		},
		Template: "child",
		Expected: "\n54321\n123",
	},
	parseTest{
		Name:     "Fake Filter Block",
		Content:  ":wat",
		Expected: ":wat",
	},
	parseTest{
		Name:     "Fake Nested Filter Block",
		Content:  ":wat\n  two",
		Expected: ":wat\n  two",
	},
	parseTest{
		Name:     "Nested Filter Block",
		Content:  ":plain\n  two",
		Expected: "two",
	},
	parseTest{
		Name:     "Nested Filter Block",
		Content:  ":plain\n  two",
		Expected: "two",
	},
	parseTest{
		Name:     "Empty tag",
		Content:  "img",
		Expected: "<img />",
	},
	parseTest{
		Name:     "Empty tag",
		Content:  "span\n  name",
		Expected: "<span>name</span>",
	},
	parseTest{
		Name:     "Class div",
		Content:  ".wat\n  here",
		Expected: "<div class=\"wat\">here</div>",
	},
	parseTest{
		Name:     "Passthrough HTML",
		Content:  "h1\n  <span>no</span>",
		Expected: "<h1><span>no</span></h1>",
	},
	parseTest{
		Name:     "Tag with content",
		Content:  "h1 Here",
		Expected: "<h1>Here</h1>",
	},
	parseTest{
		Name:     "Interpolated content",
		Content:  "Welcome {{ . }}",
		Expected: "Welcome Gopher",
		Data:     "Gopher",
	},
	parseTest{
		Name:     "With line",
		Content:  ">.Username\n  = .",
		Expected: "Gopher",
		Data:     map[string]string{"Username": "Gopher"},
	},
	parseTest{
		Name:     "With line using variable",
		Content:  ">.Username:$name\n  = $name",
		Expected: "Andrew",
		Data:     map[string]string{"Username": "Andrew"},
	},
	parseTest{
		Name:     "With line",
		Content:  ">.NotExist\n  = .\n!>\n  Other animal",
		Expected: "Other animal",
		Data:     map[string]string{"Username": "Gopher"},
	},
	parseTest{
		Name:     "Tag with attribute",
		Content:  `input type="checkbox"`,
		Expected: `<input type="checkbox" />`,
	},
	parseTest{
		Name:     "Tag with attribute and single quotes",
		Content:  `input type='checkbox'`,
		Expected: `<input type="checkbox" />`,
	},
	parseTest{
		Name:     "Tag with dot attribute",
		Content:  `input type=.Type`,
		Expected: `<input type="checkbox" />`,
		Data:     map[string]string{"Type": "checkbox"},
	},
	parseTest{
		Name: "Tag with variable attribute",
		Content: `= $t := "checkbox"
input type=$t`,
		Expected: `
<input type="checkbox" />`,
	},
	parseTest{
		Name:     "Tag with parentheses attribute",
		Content:  `input type=(print "checkbox")`,
		Expected: `<input type="checkbox" />`,
	},
	parseTest{
		Name:     "Tag with function attribute",
		Content:  `input type=checker`,
		Expected: `<input type="checkbox" />`,
		Funcs: template.FuncMap{
			"checker": func() string {
				return "checkbox"
			},
		},
	},
	parseTest{
		Name:     "If with trailing content",
		Content:  "?true\n  Yes\ntrailer",
		Expected: "Yes\ntrailer",
	},
	parseTest{
		Name:     "Collapsed Tags",
		Content:  "table > tr > %td\n  I'm different",
		Expected: "<table><tr><td>I'm different</td></tr></table>",
	},
	parseTest{
		Name:     "Totem Pole Tags",
		Content:  "table > tr > td",
		Expected: "<table><tr><td></td></tr></table>",
	},
	parseTest{
		Name:     "Collapsed tags with content on line",
		Content:  "table > tr > td Hello",
		Expected: "<table><tr><td>Hello</td></tr></table>",
	},
	parseTest{
		Name:     "Nested If/Else Statements",
		Content:  "?true\n  ?false\n    1\n  !?\n    2\n!?\n  3",
		Expected: "2",
	},
	parseTest{
		Name: "Nested Else things",
		Content: `?true
  >func1
    = .
    news
  !>
    bad
    wherever
!?
  unknown
  iops`,
		Expected: "asdf\nnews",
		Funcs: template.FuncMap{
			"func1": func() string {
				return "asdf"
			},
		},
	},
	parseTest{
		Name:     "Percentage Sign",
		Content:  "%",
		Expected: "%",
	},
	parseTest{
		Name:     "Blank Yield",
		Content:  "@",
		Expected: "",
	},
	parseTest{
		Name:     "Equals sign",
		Content:  "=",
		Expected: "=",
	},
	parseTest{
		Name:     "Single slash",
		Content:  "/",
		Expected: "/",
	},
	parseTest{
		Name:     "Blank comment",
		Content:  "//",
		Expected: "",
	},
	parseTest{
		Name:     "Cont, but actually verbatim",
		Content:  "/=",
		Expected: "=",
	},
	parseTest{
		Name:     "Malformed Block 1",
		Content:  "[name",
		Expected: "[name",
	},
	parseTest{
		Name:     "Malformed Block 2",
		Content:  "[",
		Expected: "[",
	},
	parseTest{
		Name:     "Malformed Block 3",
		Content:  "[]",
		Expected: "[]",
	},
	parseTest{
		Name:     "Malformed Exec Block 1",
		Content:  "$name",
		Expected: "$name",
	},
	parseTest{
		Name:     "Malformed Exec Block 2",
		Content:  "$",
		Expected: "$",
	},
	parseTest{
		Name:     "Malformed Exec Block 3",
		Content:  "$]",
		Expected: "$]",
	},
	parseTest{
		Name:     "Malformed Define Block 1",
		Content:  "^name",
		Expected: "^name",
	},
	parseTest{
		Name:     "Malformed Define Block 2",
		Content:  "^",
		Expected: "^",
	},
	parseTest{
		Name:     "Malformed Define Block 3",
		Content:  "^]",
		Expected: "^]",
	},
	parseTest{
		Name:     "Multi-line attributes",
		Content:  "textarea(\n  data-template=\"checkbox\"\n  name=\"blah\"\n  )\n  blah\n",
		Expected: "<textarea data-template=\"checkbox\" name=\"blah\">blah</textarea>",
	},
	parseTest{
		Name:     "Interpolated attribute with quote marks",
		Content:  `span name="first_{{ printf "%d_name" 4 }}" wat`,
		Expected: `<span name="first_4_name">wat</span>`,
	},
	parseTest{
		Name:     "Interpolated attribute with single quote marks",
		Content:  `span name='first_{{ printf "%v_name" 'a' }}' wat`,
		Expected: `<span name="first_97_name">wat</span>`,
	},
	parseTest{
		Name:     "JS Filter",
		Content:  ":js\n  var name=\"{{.}}\";\n  alert('here');",
		Data:     "Andrew\";alert('XSS');\"",
		Expected: `<script type="text/javascript">var name="Andrew\x22;alert(\x27XSS\x27);\x22";alert('here');</script>`,
	},
}

type parseTest struct {
	Name     string
	Content  string
	Sources  map[string]string
	Expected string
	Template string
	Funcs    template.FuncMap
	Data     interface{}
}
