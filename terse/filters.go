package terse

// Filters is map of available functions that take in the text from the
// template, then run transformations on it, then return it. The boolean
// indicates whether the transformed string may have delimited template
// code within it, so it should run it though a stdlib template parser.
// You are expected to do any escaping that may be needed for your
// text, and your filter will only receive the content once, when the
// template is compiled.
var Filters = map[string]FilterFunc{
	"plain": func(s string) (string, bool) {
		return s, true
	},
	"javascript": func(s string) (string, bool) {
		return `<script type="text/javascript">` + s + "</script>", true
	},

	"js": func(s string) (string, bool) {
		return `<script type="text/javascript">` + s + "</script>", true
	},
	"css": func(s string) (string, bool) {
		return "<style>" + s + "</style>", true
	},
}

type FilterFunc func(string) (string, bool)
