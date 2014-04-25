package terse

// Filters is map of available functions that take in the text from the
// template, then run transformations on it, then return it. The boolean
// indicates whether the transformed string may have delimited template
// code within it, so it should run it though a stdlib template parser.
var Filters = map[string]FilterFunc{
	"plain": func(s string) (string, bool) {
		return s, true
	},
}

type FilterFunc func(string) (string, bool)
