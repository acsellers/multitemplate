package terse

func tokenize(rt rawTree) tokenTree {
	tt := tokenTree{}
	var current token
	for _, root := range rt.Children {
		tt.roots = append(tt.roots, codeTokenizer(root.Code)(root))
	}

	return tt
}

func codeTokenizer(code string) func(*rawNode) (token, error) {
	switch {
	case doctypeCode(code):
		return doctypeToken
	case execCode(code):
		return execToken
	case commentCode(code):
		return commentToken
	case verbatimCode(code):
		return verbatimToken
	case tagCode(code):
		return tagToken
	case filterCode(code):
		return filterToken
	case blockCode(code):
		return blockToken
	case defineBlockCode(code):
		return defineBlockToken
	case execBlockCode(code):
		return execBlockToken
	case yieldCode(code):
		return yieldToken
	case ifCode(code):
		return ifToken
	case elseCode(code):
		return elseToken
	case rangeCode(code):
		return rangeToken
	case rangeElseCode(code):
		return rangeElseToken
	case idClassCode(code):
		return idClassToken
	default:
		return textToken
	}
}
