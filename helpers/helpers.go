package helpers

import (
	"fmt"
	"html/template"
	"path"
	"strings"

	"github.com/robfig/revel"

	"github.com/acsellers/helpers/assets"
	"github.com/acsellers/helpers/forms"
	"github.com/acsellers/helpers/inputs"
	"github.com/acsellers/helpers/js"
)

var (
	ALL = []string{
		"js_include",
	}
	JavascriptDirectory = "/public/js/"
	PrettyOutput        = true

	fMap = map[string]interface{}{
		"js_include": JavascriptIncludeTag,
	}
)

func Load(functions []string) {
	LoadInto(functions, revel.TemplateFuncs)
}

func LoadInto(functions []string, funcMap map[string]interface{}) {
	for _, f := range functions {
		if fun, ok := fMap[f]; ok {
			funcMap[f] = fun
		}
	}
}

func JavascriptIncludeTag(files ...string) template.HTML {
	output := make([]string, len(files))
	for i, file := range files {
		if !strings.HasSuffix(file, ".js") {
			file = file + ".js"
		}

		output[i] = fmt.Sprintf(
			`<script src="%s" type="text/javascript" charset="utf-8"></script>`,
			template.HTMLEscapeString(path.Join(JavascriptDirectory, file)),
		)
	}
	if PrettyOutput {
		return template.HTML(strings.Join(output, "\n"))
	} else {
		return template.HTML(strings.Join(output, ""))
	}
}
