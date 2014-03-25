//
package inputs

import (
	"html/template"
)

var (
	Functions = map[string]interface{}{
		"button_tag":          ButtonTag,
		"check_box_tag":       CheckBoxTag,
		"email_field_tag":     EmailFieldTag,
		"form_tag":            CreateFormTag,
		"open_form_tag":       OpenFormTag,
		"close_form_tag":      CloseFormTag,
		"hidden_field_tag":    HiddenFieldTag,
		"image_submit_tag":    ImageSubmitTag,
		"label_tag":           LabelTag,
		"number_field_tag":    NumberFieldTag,
		"password_field_tag":  PasswordFieldTag,
		"phone_field_tag":     PhoneFieldTag,
		"radio_button_tag":    RadioButtonTag,
		"range_field_tag":     RangeFieldTag,
		"search_field_tag":    SearchFieldTag,
		"select_tag":          SelectTag,
		"submit_tag":          SubmitTag,
		"telephone_field_tag": TelephoneFieldTag,
		"text_area_tag":       TextAreaTag,
		"text_field_tag":      TextFieldTag,
		"csrf_token_tag":      CSRFTokenTag,
		"url_field_tag":       UrlFieldTag,
		"utf8_tag":            Utf8Tag,
	}
)
