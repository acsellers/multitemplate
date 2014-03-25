package mustache

import (
	"fmt"
	htmlTemplate "html/template"
	"reflect"
	textTemplate "text/template"
	"text/template/parse"
)

type MustacheParser string

func (mp MustacheParser) ParseFile(name, content string) (map[string]*parse.Tree, error) {
	return Parse(name, content)
}
func (mp MustacheParser) RequiredHtmlFuncs() htmlTemplate.FuncMap {
	return htmlTemplate.FuncMap{
		"mustacheIsCollection": func(i interface{}) bool {
			it := reflect.TypeOf(i)
			switch it.Kind() {
			case reflect.Array, reflect.Slice:
				return true
			default:
				return false
			}
		},
		"mustacheUnescape": func(i ...interface{}) htmlTemplate.HTML {
			if len(i) == 1 {
				return htmlTemplate.HTML(fmt.Sprint(i[0]))
			}
			return htmlTemplate.HTML("")
		},
		"mustacheUpscope":   upscope,
		"mustacheDownscope": downscope,
	}
}
func (mp MustacheParser) RequiredTextFuncs() textTemplate.FuncMap {
	return textTemplate.FuncMap{
		"mustacheIsCollection": func(i interface{}) bool {
			it := reflect.TypeOf(i)
			switch it.Kind() {
			case reflect.Array, reflect.Slice:
				return true
			default:
				return false
			}
		},
		"mustacheUnescape": func(i ...interface{}) string {
			if len(i) == 1 {
				return fmt.Sprint(i[0])
			}
			return ""
		},
		"mustacheUpscope":   upscope,
		"mustacheDownscope": downscope,
	}
}
func (mp MustacheParser) ApplicableExtensions() []string {
	return []string{"mustache", "ms", "mche"}
}
