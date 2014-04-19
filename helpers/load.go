package helpers

import (
	"html/template"

	"github.com/acsellers/multitemplate"
)

// Available helper modules "forms", "general", "link", "asset"
func LoadHelpers(modules ...string) {
	loadFuncs(coreFuncs)
	for _, module := range modules {
		switch module {
		case "all":
			loadFuncs(formTagFuncs)
			loadFuncs(selectTagFuncs)
			loadFuncs(generalFuncs)
			loadFuncs(linkFuncs)
			loadFuncs(assetFuncs)
		case "forms":
			loadFuncs(formTagFuncs)
			loadFuncs(selectTagFuncs)
		case "general":
			loadFuncs(generalFuncs)
		case "link":
			loadFuncs(linkFuncs)
		case "asset":
			loadFuncs(assetFuncs)
		}
	}
}

func loadFuncs(tf template.FuncMap) {
	for k, f := range tf {
		multitemplate.LoadedFuncs[k] = f
	}
}
