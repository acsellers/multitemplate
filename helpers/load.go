package helpers

import (
	"html/template"

	"github.com/acsellers/multitemplate"
)

// LoadHelpers loads helper functions into the
// multitemplate function map. All modules
// depend on a "core" module that will always be loaded. The modules may
// be all be loaded by asking for the "all" module, or they can be loaded
// by their names, which are "form", "general", "link" and "asset".
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
		case "form":
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

// GetHelpers loads helper functions into a html/template FuncMap
// then returns that FuncMap. Since helpers does not depend on
// any special functions from multitemplate, this would allow
// you to use these helpers in any Go template library that allows
// you to add helper functions.
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
