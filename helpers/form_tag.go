package helpers

import (
	"fmt"
	"html/template"
	"strings"
)

var formTagFuncs = template.FuncMap{
	"button_tag": func(text string, options ...AttrList) template.HTML {
		al := combine("", "", options)
		tc := template.HTML(template.HTMLEscapeString(text))
		defaultName := strings.Map(nameFilter, text)
		if _, ok := al["name"]; !ok {
			al["name"] = defaultName
		}
		return buildTag("button", tc, al)
	},
	"check_box_tag": func(name string, value interface{}, options ...AttrList) template.HTML {
		al := combine(name, fmt.Sprint(value), options)
		al["type"] = "checkbox"
		return buildTag("input", "", al)
	},
	"email_field_tag": func(name, value string, options ...AttrList) template.HTML {
		al := combine(name, value, options)
		al["type"] = "email"
		return buildTag("input", "", al)
	},
	"fieldset_tag": func(name string, options ...AttrList) template.HTML {
		al := combine("", "", options)
		al["MT_skip_close"] = true
		return buildTag(
			"fieldset",
			template.HTML("<template"+template.HTMLEscapeString(name)+"</template>"),
			al,
		)
	},
	"end_fieldset_tag": func() template.HTML {
		return "</fieldset>"
	},
	"file_field_tag": func(name string, options ...AttrList) template.HTML {
		al := combine(name, "", options)
		al["type"] = "file"
		return buildTag("input", "", al)
	},
	"form_tag": func(target string, options ...AttrList) template.HTML {
		al := combine("", "", options)
		al["action"] = target
		if _, ok := al["method"]; !ok {
			al["method"] = "post"
		}
		al["MT_skip_close"] = true
		return buildTag("form", "", al)
	},
	"end_form_tag": func() template.HTML {
		return "</form>"
	},
	"hidden_field_tag": func(name string, value interface{}, options ...AttrList) template.HTML {
		al := combine(name, fmt.Sprint(value), options)
		al["type"] = "hidden"
		return buildTag("input", "", al)
	},
	"label_tag": func(target, text string, options ...AttrList) template.HTML {
		al := combine("", "", options)
		al["for"] = strings.Map(nameFilter, target)
		return buildTag("input", template.HTML(template.HTMLEscapeString(text)), al)
	},
	"number_field_tag": func(name string, value interface{}, options ...AttrList) template.HTML {
		al := combine(name, fmt.Sprint(value), options)
		al["type"] = "number"
		return buildTag("input", "", al)
	},
	"password_field_tag": func(name string, value interface{}, options ...AttrList) template.HTML {
		al := combine(name, fmt.Sprint(value), options)
		al["type"] = "password"
		return buildTag("input", "", al)
	},
	"phone_field_tag": func(name string, value interface{}, options ...AttrList) template.HTML {
		al := combine(name, fmt.Sprint(value), options)
		al["type"] = "tel"
		return buildTag("input", "", al)
	},
	"radio_button_tag": func(name string, value interface{}, active bool, options ...AttrList) template.HTML {
		al := combine(name, fmt.Sprint(value), options)
		if active {
			al["checked"] = "checked"
		}
		al["type"] = "radio"
		return buildTag("input", "", al)
	},
	"range_field_tag": func(name string, value interface{}, options ...AttrList) template.HTML {
		al := combine(name, fmt.Sprint(value), options)
		al["type"] = "range"
		return buildTag("input", "", al)
	},
	"search_field_tag": func(name string, value interface{}, options ...AttrList) template.HTML {
		al := combine(name, fmt.Sprint(value), options)
		al["type"] = "search"
		return buildTag("input", "", al)
	},
	"submit_tag": func(name string, options ...AttrList) template.HTML {
		al := combine(name, "", options)
		al["type"] = "submit"
		return buildTag("input", "", al)
	},
	"text_area_tag": func(name string, value string, options ...AttrList) template.HTML {
		al := combine(name, "", options)
		if b, ok := al["disable_escape"].(bool); ok && b {
			return buildTag("textarea", template.HTML(value), al)
		}
		return buildTag("textarea", template.HTML(template.HTMLEscapeString(value)), al)
	},
	"text_field_tag": func(name string, value string, options ...AttrList) template.HTML {
		al := combine(name, value, options)
		al["type"] = "text"
		return buildTag("input", "", al)
	},
	"url_field_tag": func(name string, value string, options ...AttrList) template.HTML {
		al := combine(name, value, options)
		al["type"] = "url"
		return buildTag("input", "", al)
	},
	"utf8_tag": func() template.HTML {
		al := AttrList{}
		al["type"] = "hidden"
		al["name"] = "utf8"
		al["value"] = "&#x269b"
		return buildTag("input", "", al)
	},
}

func combine(name, value string, opts []AttrList) AttrList {
	al := AttrList{}
	for _, attrs := range opts {
		for k, v := range attrs {
			if _, ok := al[k]; !ok {
				al[k] = v
			}
		}
	}
	if name != "" {
		al["name"] = strings.Map(nameFilter, name)
		if _, ok := al["id"]; !ok {
			al["id"] = strings.Map(nameFilter, name)
		}
	}
	if value != "" {
		al["value"] = value
	}

	return al
}

/*
  Ignored tags
  * color_field_tag
  * date_field_tag
  * datetime_field_tag
  * datetime_local_field_tag
  * field_set_tag
  * image_submit_tag
  * month_tag
  * time_field_tag
  * week_field_tag
*/
