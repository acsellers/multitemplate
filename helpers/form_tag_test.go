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
				case func() template.HTML:
					r = string(af())
				case func(string, interface{}, ...AttrList) template.HTML:
					r = string(af(formTagTest.Args[0], formTagTest.Args[1], formTagTest.Attrs...))
				case func(string, string, ...AttrList) template.HTML:
					r = string(af(formTagTest.Args[0], formTagTest.Args[1], formTagTest.Attrs...))
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
	helperTest{
		Helper:   "email_field_tag",
		Args:     []string{"user_email", ""},
		Expected: `<input name="user_email" id="user_email" type="email" />`,
	},
	helperTest{
		Helper:   "email_field_tag",
		Args:     []string{"user_email", "you@example.com"},
		Expected: `<input name="user_email" id="user_email" value="you@example.com" type="email" />`,
	},
	helperTest{
		Helper:   "fieldset_tag",
		Args:     []string{"Options"},
		Expected: "<fieldset><legend>Options</legend>",
	},
	helperTest{
		Helper:   "file_field_tag",
		Args:     []string{"picture"},
		Expected: `<input name="picture" id="picture" type="file" />`,
	},
	helperTest{
		Helper:   "form_tag",
		Args:     []string{"/users/create"},
		Expected: `<form action="/users/create" method="post">`,
	},
	helperTest{
		Helper:   "hidden_field_tag",
		Args:     []string{"user_id", "123456"},
		Expected: `<input name="user_id" id="user_id" value="123456" type="hidden" />`,
	},
	helperTest{
		Helper:   "label_tag",
		Args:     []string{"user_name", "What's ya name?"},
		Expected: `<label for="user_name">What&#39;s ya name?</label>`,
	},
	helperTest{
		Helper:   "number_field_tag",
		Args:     []string{"open_slots", "12"},
		Expected: `<input name="open_slots" id="open_slots" value="12" type="number" />`,
	},
	helperTest{
		Helper:   "password_field_tag",
		Args:     []string{"user_pass"},
		Expected: `<input name="user_pass" id="user_pass" type="password" />`,
	},
	helperTest{
		Helper:   "phone_field_tag",
		Args:     []string{"telephone", "555-5555"},
		Expected: `<input name="telephone" id="telephone" value="555-5555" type="tel" />`,
	},
	/*
		helperTest{
			Helper:   "radio_button_tag",
			Args:     []string{"wake", "yes"},
			Expected: ``,
		},
	*/
	helperTest{
		Helper:   "range_field_tag",
		Args:     []string{"date_ages", ""},
		Expected: `<input name="date_ages" id="date_ages" type="range" />`,
	},
	helperTest{
		Helper:   "search_field_tag",
		Args:     []string{"search"},
		Expected: `<input name="search" id="search" type="search" />`,
	},
	helperTest{
		Helper:   "submit_tag",
		Args:     []string{"create"},
		Expected: `<input name="create" id="create" type="submit" />`,
	},
	helperTest{
		Helper:   "text_area_tag",
		Args:     []string{"story", "blah blah blah"},
		Expected: `<textarea name="story" id="story">blah blah blah</textarea>`,
	},
	helperTest{
		Helper:   "text_field_tag",
		Args:     []string{"user_name", "acsellers"},
		Expected: `<input name="user_name" id="user_name" value="acsellers" type="text" />`,
	},
	helperTest{
		Helper:   "url_field_tag",
		Args:     []string{"website", "google.com"},
		Expected: `<input name="website" id="website" value="google.com" type="url" />`,
	},
	helperTest{
		Helper:   "utf8_tag",
		Expected: `<input type="hidden" name="utf8" value="&#x269b" />`,
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
