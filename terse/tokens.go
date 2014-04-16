package terse

type tokenTree struct {
	roots []token
}

type token interface {
	Opening() []token
	Closing() []token
	Children() []token
	Type() tokenType
	String() string
}

type tokenType int

const (
	errorTokenType tokenType = iota
	textTokenType
	conditionalTokenType
	codeTokenType
	tagTokenType
	ifTokenType
	elseTokenType
)

type errorToken struct{}

func (*errorToken) Opening() []token {
	return []token{}
}

func (*errorToken) Closing() []token {
	return []token{}
}
func (*errorToken) Children() []token {
	return []token{}
}
func (*errorToken) Type() tokenType {
	return errorTokenType
}

func (*errorToken) String() string {
	return "PARSE ERROR"
}

type tokenDoctype struct {
	Text     string
	Comments string
}

func (td *tokenDoctype) Opening() []token {
	return []token{}
}
func (td *tokenDoctype) Children() []token {
	return []token{}
}
func (td *tokenDoctype) Closing() []token {
	return []token{}
}
func (td *tokenDoctype) Type() tokenType {
	return textTokenType
}

func (td *tokenDoctype) String() string {
}
