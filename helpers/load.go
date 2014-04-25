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

func GetHelpers(modules ...string) template.FuncMap {
	tf := template.FuncMap{}
	getFuncs(tf, coreFuncs)
	for _, module := range modules {
		switch module {
		case "all":
			getFuncs(tf, formTagFuncs)
			getFuncs(tf, selectTagFuncs)
			getFuncs(tf, generalFuncs)
			getFuncs(tf, linkFuncs)
			getFuncs(tf, assetFuncs)
		case "forms":
			getFuncs(tf, formTagFuncs)
			getFuncs(tf, selectTagFuncs)
		case "general":
			getFuncs(tf, generalFuncs)
		case "link":
			getFuncs(tf, linkFuncs)
		case "asset":
			getFuncs(tf, assetFuncs)
		}
	}
	return tf
}

func loadFuncs(tf template.FuncMap) {
	for k, f := range tf {
		multitemplate.LoadedFuncs[k] = f
	}
}

func getFuncs(host, source template.FuncMap) {
	for k, f := range source {
		host[k] = f
	}
}
