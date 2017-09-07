package iscdhcp

type token struct {
	data []byte
	typ  int
}

type tokenStream struct {
	tokens []token
	index  int
}

func (ts *tokenStream) next() (token, error) {
	ts.index++
	// Correct for possibly having a negative index (from repeated calls to
	// ts.undo())
	if ts.index <= 0 {
		ts.index = 1
	}
	// Return an error if we've reached the end of the tokens
	if ts.index > len(ts.tokens) {
		return token{}, contextError("no more tokens")
	}
	return ts.tokens[ts.index-1], nil
}

func (ts *tokenStream) undo() {
	ts.index--
}

const (
	tokenTypeIdentifier = iota
	tokenTypeString
	tokenTypeWhiteSpace
	tokenTypeComment
	tokenTypeBlockStart
	tokenTypeBlockEnd
	tokenTypeStatementEnd
)

type tokenizer struct{}

func (t tokenizer) Tokenize(data []byte) ([]token, error) {
	var tokens []token
	var currentToken token

	scanner := &scanner{}
	scanner.init()

	for i := range data {
		b := data[i]
		code, err := scanner.step(b)
		if err != nil {
			return nil, err
		}

		switch code {
		case codeContinue:
			currentToken.data = append(currentToken.data, b)
		case codeIdentifierBegin:
			if len(currentToken.data) != 0 {
				tokens = append(tokens, currentToken)
			}
			currentToken = token{
				typ:  tokenTypeIdentifier,
				data: []byte{b},
			}
		case codeIdentifierEnd:
			fallthrough
		case codeWhitespace:
			if currentToken.typ != tokenTypeWhiteSpace {
				if len(currentToken.data) != 0 {
					tokens = append(tokens, currentToken)
					currentToken = token{
						typ:  tokenTypeWhiteSpace,
						data: []byte{b},
					}
				} else {
					currentToken = token{
						typ:  tokenTypeWhiteSpace,
						data: []byte{b},
					}
				}
			} else {
				currentToken.data = append(currentToken.data, b)
			}
		case codeBlockBegin:
			if len(currentToken.data) != 0 {
				tokens = append(tokens, currentToken)
			}
			currentToken = token{
				typ:  tokenTypeBlockStart,
				data: []byte{b},
			}
		case codeBlockEnd:
			if len(currentToken.data) != 0 {
				tokens = append(tokens, currentToken)
			}
			currentToken = token{
				typ:  tokenTypeBlockEnd,
				data: []byte{b},
			}
		case codeStatementEnd:
			if len(currentToken.data) != 0 {
				tokens = append(tokens, currentToken)
			}
			currentToken = token{
				typ:  tokenTypeStatementEnd,
				data: []byte{b},
			}
		case codeCommentBegin:
			if len(currentToken.data) != 0 {
				tokens = append(tokens, currentToken)
			}
			currentToken = token{
				typ:  tokenTypeComment,
				data: []byte{b},
			}
		case codeCommentEnd:
			tokens = append(tokens, currentToken)
			currentToken = token{
				typ:  tokenTypeWhiteSpace,
				data: []byte{b},
			}
		case codeStringBegin:
			if len(currentToken.data) != 0 {
				tokens = append(tokens, currentToken)
			}
			currentToken = token{
				typ:  tokenTypeString,
				data: []byte{b},
			}
		case codeStringEnd:
			currentToken.data = append(currentToken.data, b)
		}
	}

	if len(currentToken.data) != 0 {
		tokens = append(tokens, currentToken)
	}

	return tokens, nil
}
