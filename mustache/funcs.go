package mustache

import (
	"fmt"
	"html/template"
	"reflect"
)

var testFuncs = template.FuncMap{
	"mustacheIsCollection": func(i interface{}) bool {
		it := reflect.TypeOf(i)
		switch it.Kind() {
		case reflect.Array, reflect.Slice:
			return true
		default:
			return false
		}
	},
	"mustacheUnescape": func(i ...interface{}) template.HTML {
		if len(i) == 1 {
			return template.HTML(fmt.Sprint(i[0]))
		}
		return template.HTML("")
	},
	"mustacheUpscope":   upscope,
	"mustacheDownscope": downscope,
}
