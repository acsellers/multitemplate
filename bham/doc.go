// bham or the "block-based html abstraction markup"
// is an attempt to take what is good about
// haml, and port it to Go, but it isn't a direct port.
// bham takes advantage of go's existing template library
// and will use that template library syntax, not haml's ruby
// syntax.
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
