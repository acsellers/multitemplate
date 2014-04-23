package terse

func tokenize(rt rawTree) tokenTree {
	tt := tokenTree{}
	var current *token
	for _, root := range rt.Children {
		current, tt.err = codeTokenizer(root.Code)(root)
		if tt.err != nil {
			return tt
		}
		tt.roots = append(tt.roots, current)
	}

	return tt
}

func childTokenize(node *rawNode) ([]*token, error) {
	if len(node.Children) == 0 {
		return []*token{}, nil
	}

	tokens := []*token{}
	for _, child := range node.Children {
		current, err := codeTokenizer(child.Code)(child)
		if err != nil {
			return []*token{}, err
		}
		tokens = append(tokens, current)
	}

	return tokens, nil
}

func codeTokenizer(code string) func(*rawNode) (*token, error) {
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
	case extendCode(code):
		return extendToken
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
	case withCode(code):
		return withToken
	case withElseCode(code):
		return withElseToken
	case idClassCode(code):
		return tagToken
	default:
		return textToken
	}
}
