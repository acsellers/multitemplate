//
package forms

import (
	"html/template"
)

var (
	Functions = map[string]interface{}{
		"form_for":        FormFor,
		"check_box":       CheckBox,
		"email_field":     EmailField,
		"file_field":      FileField,
		"hidden_field":    HiddenField,
		"label":           LabelField,
		"number_field":    NumberField,
		"password_field":  PasswordField,
		"phone_field":     PhoneField,
		"radio_button":    RadioButton,
		"range_field":     RangeField,
		"search_field":    SearchField,
		"select_field":    SelectField,
		"telephone_field": TelephoneField,
		"text_area":       TextArea,
		"text_field":      TextField,
		"url_field":       UrlField,
		"submit_button":   SubmitButton,
	}
)
