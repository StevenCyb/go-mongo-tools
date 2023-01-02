package tokenizer

// TokenType is the type of the token.
type Type string

// String return as string.
func (t Type) String() string {
	return string(t)
}

// Token represents a single token.
type Token struct {
	Type  Type
	Value string
}

// NewToken creates a new token with given arguments.
func NewToken(tokenType Type, value string) *Token {
	return &Token{
		Type:  tokenType,
		Value: value,
	}
}
