package errs

import chain "github.com/g8rswimmer/error-chain"

// Chain holds multiple errors in chain.
type Chain []error

// AddIf adds error to chain if not nil.
func (c *Chain) AddIf(err error) {
	if err != nil {
		*c = append(*c, err)
	}
}

// GetError returns error as `g8rswimmer/error-chain` if contains any error.
func (c Chain) GetError() error {
	if len(c) == 0 {
		return nil
	}

	errorChain := chain.New()

	for _, err := range c {
		errorChain.Add(err)
	}

	return errorChain
}
