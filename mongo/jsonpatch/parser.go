package jsonpatch

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"

	"github.com/StevenCyb/go-mongo-tools/errs"
	"github.com/StevenCyb/go-mongo-tools/mongo/jsonpatch/operation"
	"github.com/StevenCyb/go-mongo-tools/mongo/jsonpatch/validator"
	"go.mongodb.org/mongo-driver/bson"
)

var ErrNoOperationToPerform = errors.New("no operation to perform")

// NewParser creates a new parser that uses optional policies.
func NewParser(policies ...Policy) *Parser {
	return &Parser{
		policies: policies,
	}
}

// NewSmartParser creates a new parser that create rules based on given type.
func NewSmartParser(reference reflect.Type) (*Parser, error) {
	if reference == nil {
		return nil, validator.ErrReferenceIsNil
	}

	validator, err := validator.NewValidator(reference)
	if err != nil {
		return nil, fmt.Errorf("failed to create rules from reference: %w", err)
	}

	return &Parser{
		validator: validator,
	}, nil
}

// Parser that can parse patch operation to generate mongo queries.
type Parser struct {
	validator *validator.Validator
	policies  []Policy
}

// Parse given operation spec to generate mongo queries if not violating policies.
func (p Parser) Parse(operationSpecs ...operation.Spec) (bson.A, error) {
	if len(operationSpecs) == 0 {
		return nil, ErrNoOperationToPerform
	}

	for _, policy := range p.policies {
		for _, operationSpec := range operationSpecs {
			if !operationSpec.Valid() {
				return nil, errs.NewErrUnexpectedInput(operationSpec)
			}

			if !policy.Test(operationSpec) {
				return nil, errs.NewErrPolicyViolation(policy.GetDetails())
			}
		}
	}

	if p.validator != nil {
		for _, operationSpec := range operationSpecs {
			err := p.validator.Validate(operationSpec)
			if err != nil {
				return nil, fmt.Errorf("operation '%+v' invalid: %w", operationSpec, err)
			}
		}
	}

	return p.generateMongoQuery(operationSpecs...)
}

// generateMongoQuery generates the mongo query out of operation spec.
//
//nolint:funlen
func (p Parser) generateMongoQuery(operationSpecs ...operation.Spec) (bson.A, error) {
	var (
		element  bson.M
		query    = bson.A{}
		noSuffix = regexp.MustCompile(`\.[0-9]+$`)
	)

	for _, operationSpec := range operationSpecs {
		if operationSpec.Value != nil &&
			fmt.Sprintf("%v", operationSpec.Value) == "map[]" ||
			fmt.Sprintf("%v", operationSpec.Value) == "{}" {
			operationSpec.Value = nil
		}

		switch operationSpec.Operation {
		case operation.RemoveOperation:
			if noSuffix.Match([]byte(operationSpec.Path)) {
				extract := regexp.MustCompile(`^(?P<path>.*)\.(?P<index>[0-9]+)$`)
				match := extract.FindStringSubmatch(string(operationSpec.Path))
				paramsMap := make(map[string]string)

				for i, name := range extract.SubexpNames() {
					if i > 0 && i <= len(match) {
						paramsMap[name] = match[i]
					}
				}

				path := paramsMap["path"]
				index, _ := strconv.ParseInt(paramsMap["index"], 10, 64)
				element = bson.M{
					"$set": bson.M{
						path: bson.M{
							"$concatArrays": bson.A{
								bson.M{
									"$slice": bson.A{
										"$" + path,
										index,
									},
								},
								bson.M{
									"$slice": bson.A{
										"$" + path,
										bson.M{
											"$add": bson.A{1, index},
										},
										bson.M{"$size": "$" + path},
									},
								},
							},
						},
					},
				}
			} else {
				element = bson.M{
					"$unset": string(operationSpec.Path),
				}
			}
		case operation.AddOperation:
			if reflect.TypeOf(operationSpec.Value).Kind() != reflect.Slice {
				operationSpec.Value = []interface{}{operationSpec.Value}
			}

			element = bson.M{
				"$set": bson.M{
					string(operationSpec.Path): bson.M{
						"$concatArrays": bson.A{
							"$" + string(operationSpec.Path),
							operationSpec.Value,
						},
					},
				},
			}
		case operation.ReplaceOperation:
			query = append(query, bson.M{
				"$unset": string(operationSpec.Path),
			})
			element = bson.M{
				"$set": bson.M{
					string(operationSpec.Path): operationSpec.Value,
				},
			}
		case operation.MoveOperation:
			query = append(query, bson.M{
				"$set": bson.M{
					string(operationSpec.Path): "$" + string(operationSpec.From),
				},
			})
			element = bson.M{
				"$unset": string(operationSpec.From),
			}
		case operation.CopyOperation:
			element = bson.M{
				"$set": bson.M{
					string(operationSpec.Path): "$" + string(operationSpec.From),
				},
			}
		}

		query = append(query, element)
	}

	return query, nil
}
