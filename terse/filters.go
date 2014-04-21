package terse

var Filters = map[string]FilterFunc{
	"plain": func(s string) (string, bool) {
		return s, true
	},
}

type FilterFunc func(string) (string, bool)
