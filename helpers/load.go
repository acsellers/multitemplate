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
			loadFuncs(generalFuncs)
			loadFuncs(linkFuncs)
		case "forms":
			loadFuncs(formTagFunctions)
		case "general":
			loadFuncs(generalFuncs)
		case "link":
			loadFuncs(linkFuncs)
		}
	}
}

func loadFuncs(tf template.FuncMap) {
	for k, f := range tf {
		multitemplate.LoadedFuncs[k] = f
	}
}
