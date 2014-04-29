package helpers

import (
	"html/template"
	"testing"
	. "github.com/acsellers/assert"
)

func TestSelectTag(t *testing.T) {
	Within(t, func(test *Test) {
		f := selectTagFuncs["select_tag"].(func(string, interface{}, ...AttrList) template.HTML)
		r := string(f("wat", OptionList{Option{"a", "a"}}))
		test.AreEqual(r, buildTag("select", `<option value="a">a</option>`, AttrList{"name": "wat", "id": "wat"}))
	})
}

func TestOptionGroup(t *testing.T) {
	Within(t, func(test *Test) {
		f := selectTagFuncs["select_tag"].(func(string, interface{}, ...AttrList) template.HTML)
		r := string(f("wat", OptionList{OptionGroup{"blah", []Option{Option{"a", "a"}}}}))
		test.AreEqual(r, `<select name="wat" id="wat"><optgroup label="blah"><option value="a">a</option></optgroup></select>`)
	})
}

func TestMultipleOptions(t *testing.T) {
	Within(t, func(test *Test) {
		f := selectTagFuncs["select_tag"].(func(string, interface{}, ...AttrList) template.HTML)
		r := string(f("wat", OptionList{Option{"a", "a"}, Option{"b", "b"}}))
		test.AreEqual(`<select name="wat" id="wat"><option value="a">a</option><option value="b">b</option></select>`, r)
	})
}

func TestOption(t *testing.T) {
	Within(t, func(test *Test) {
		f := selectTagFuncs["option"].(func(interface{}, interface{}) Option)
		r := f("one", "two")
		test.AreEqual(r, Option{"one", "two"})
		test.AreEqual(r.ToHTML(), `<option value="two">one</option>`)
	})
}

func TestOptionsWithValues(t *testing.T) {
	Within(t, func(test *Test) {
		f := selectTagFuncs["options_with_values"].(func(...interface{}) OptionList)
		r := f("one", "two", "three", "four")
		test.AreEqual(OptionList{Option{"one", "two"}, Option{"three", "four"}}, r)
		test.AreEqual(r.ToHTML(), `<option value="two">one</option><option value="four">three</option>`)
	})
}

func TestOptions(t *testing.T) {
	Within(t, func(test *Test) {
		f := selectTagFuncs["options"].(func(...interface{}) OptionList)
		r := f("a", "b")
		test.AreEqual(r, OptionList{Option{"a", "a"}, Option{"b", "b"}})
	})
}

func TestGroupOptions(t *testing.T) {
	Within(t, func(test *Test) {
		f := selectTagFuncs["group_options"].(func(string, ...OptionLike) OptionGroup)
		r := f("wat", OptionList{Option{"a", "b"}})
		test.AreEqual(r.ToHTML(), `<optgroup label="wat"><option value="b">a</option></optgroup>`)
	})
}
