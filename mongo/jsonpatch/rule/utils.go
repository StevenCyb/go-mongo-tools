package rule

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/StevenCyb/go-mongo-tools/mongo/jsonpatch/operation"
)

var (
	ErrInvalidKind   = errors.New("is invalid kind")
	ErrInvalidBool   = errors.New("is invalid bool")
	ErrInvalidNumber = errors.New("is invalid number")
)

// getBoolIfNotEmpty parse string to bool or throw error.
func getBoolIfNotEmpty(valueString, path, key string) (*bool, error) {
	value, err := strconv.ParseBool(valueString)
	if err != nil {
		return nil, fmt.Errorf("rule `%s` at '%s': %w", key, path, ErrInvalidBool)
	}

	return &value, nil
}

// getFloat64IfNotEmpty parse string to float64 or throw error.
func getFloat64IfNotEmpty(valueString, path, key string) (*float64, error) {
	value, err := strconv.ParseFloat(valueString, 64)
	if err != nil {
		return nil, fmt.Errorf("rule `%s` at '%s': %w", key, path, ErrInvalidNumber)
	}

	return &value, nil
}

// getRegexpIfNotEmpty parse string to regex expression or throw error.
func getRegexpIfNotEmpty(valueString, path, key string) (*regexp.Regexp, error) {
	value, err := regexp.Compile(valueString)
	if err != nil {
		return nil, fmt.Errorf("rule `%s` at '%s': %w", key, path, err)
	}

	return value, nil
}

// getOperationsIfNotEmpty parse string to patch operations or throw error.
func getOperationsIfNotEmpty(valueString, path, key string) (*[]operation.Operation, error) {
	operations := []operation.Operation{}

	values := strings.Split(valueString, ",")

	for _, value := range values {
		operation, err := operation.FromString(value)
		if err != nil {
			return nil, fmt.Errorf("rule `%s` at '%s': %w", key, path, err)
		}

		operations = append(operations, *operation)
	}

	return &operations, nil
}
