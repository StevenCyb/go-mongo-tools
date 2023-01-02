package tokenizer

import (
	"regexp"
)

// spec for the tokenizer.
type Spec struct {
	expression *regexp.Regexp
	tokenType  Type
}

// NewSpec creates a new spec with given arguments.
func NewSpec(expression string, tokenType Type) *Spec {
	return &Spec{
		expression: regexp.MustCompile(expression),
		tokenType:  tokenType,
	}
}
