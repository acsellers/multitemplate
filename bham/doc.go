// bham or the "blocky hypertext abstraction markup"
// is an attempt to take what is good about languages
// like haml, jade, slim, etc. and port it to Go, but
// not blindly. It will take into account the capabilities
// of Go's template libraries to parse directly into the
// internal template structures that the stdlib template
// libraries use to provide both speed and interoperability
// with standard Go templates.
package bham

var (
	// Strict determines whether only tabs will be considered
	// as indentation operators (Strict == true) or whether
	// two spaces can be counted as an indentation operator
	// (Strict == false), this is included for haml
	// semi-comapibility
	Strict bool

	// To add multiple id declarations, the outputter puts them together
	// with a join string, by default this is an underscore
	IdJoin = "_"

	// Like the template library, you need to be able to set code delimeters
	LeftDelim  = "{{"
	RightDelim = "}}"
)
