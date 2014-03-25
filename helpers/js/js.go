//
package js

import (
	"html/template"
)

var (
	Functions = map[string]interface{}{
		"button_to_function": ButtonToFunction,
		"javascript_tag":     JavascriptTag,
		"link_to_function":   LinkToFunction,
	}
)

func ButtonToFunction(buttonTitle, jsCode string, options ...string) template.HTML {

}

func JavascriptTag(jsCode string, options ...string) template.HTML {

}

func LinkToFunction(linkTitle, jsCode string, options ...string) template.HTML {

}
