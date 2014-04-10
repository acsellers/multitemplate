package helpers

import (
	"html/template"

	"github.com/acsellers/multitemplate"
)

func LoadHelpers(modules ...string) {
	loadFuncs(coreFunctions)
	for _, module := range modules {
		switch module {
		case "all":
			loadFuncs(formTagFunctions)
			loadFuncs(jsTagFunctions)
		case "forms":
			loadFuncs(formTagFunctions)
		}
	}
}

func loadFuncs(tf template.FuncMap) {
	for k, f := range tf {
		multitemplate.LoadedFuncs[k] = f
	}
}