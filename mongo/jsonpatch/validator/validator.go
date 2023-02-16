package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/StevenCyb/go-mongo-tools/mongo/jsonpatch/forcecast"
	"github.com/StevenCyb/go-mongo-tools/mongo/jsonpatch/operation"
	"github.com/StevenCyb/go-mongo-tools/mongo/jsonpatch/rule"
)

const (
	prefix = "jp_"
)

// Validator interprets reference to validate JSON patch operations.
type Validator struct {
	knownForceCast map[string]forcecast.ForceCast
	forceCast      map[operation.Path]forcecast.ForceCast
	knownTagRules  map[string]rule.Rule
	generalRules   map[string]rule.Rule
	rules          map[operation.Path]map[string]rule.Rule
	wildcardRules  map[operation.Path]map[string]rule.Rule
}

// RegisterRule register addition rule for specific key (must have prefix `jp_`).
func (v *Validator) RegisterRule(key string, rule rule.Rule) error {
	if rule == nil {
		return ErrNilRule
	}

	if !strings.HasPrefix(key, prefix) {
		return ErrMissingPrefix
	}

	if _, exists := v.knownTagRules[key]; exists {
		return ErrDuplicateRuleTags
	}

	if _, exists := v.generalRules[key]; exists {
		return ErrDuplicateRuleTags
	}

	v.knownTagRules[key] = rule

	return nil
}

// Validate a given JSON patch operations again rules.
func (v Validator) Validate(operationSpec operation.Spec) error {
	if forceCast, match := v.forceCast[operationSpec.Path]; match {
		var err error

		operationSpec.Value, err = forceCast.Cast(operationSpec.Value)
		if err != nil {
			return fmt.Errorf("force cast failed: %w", err)
		}
	}

	if rules, match := v.rules[operationSpec.Path]; match {
		for _, rule := range rules {
			err := rule.Validate(operationSpec)
			if err != nil {
				return fmt.Errorf("operation no allowed: %w", err)
			}
		}

		return nil
	}

	for path, rules := range v.wildcardRules {
		if path.Equal(operationSpec.Path) {
			for _, rule := range rules {
				err := rule.Validate(operationSpec)
				if err != nil {
					return fmt.Errorf("operation no allowed: %w", err)
				}
			}

			return nil
		}
	}

	return UnknownPathError{path: string(operationSpec.Path)}
}

// UseReference interpret given reference to model rule set.
func (v *Validator) UseReference(referenceType reflect.Type) error {
	return v.parseReference(referenceType, "", map[string]rule.Rule{})
}

func (v *Validator) parseReference(objectType reflect.Type, path string, inheritedRules map[string]rule.Rule) error {
	var err error

	if objectType == nil {
		return ErrReferenceIsNil
	}

	if path != "" && objectType.Kind() != reflect.Ptr {
		path += "."
	}

	switch objectType.Kind() { //nolint:exhaustive
	case reflect.Ptr:
		err = v.parseReference(reflect.Zero(objectType.Elem()).Type(), path, inheritedRules)
	case reflect.Array, reflect.Map, reflect.Slice:
		err = v.parseReferenceIterable(objectType, path, inheritedRules)
	case reflect.Struct:
		err = v.parseReferenceStruct(objectType, path, inheritedRules)
	}

	return err
}

func (v *Validator) parseReferenceIterable(
	objectType reflect.Type, path string, inheritedRules map[string]rule.Rule,
) error {
	var (
		zeroValue reflect.Value
		valueType = objectType.Elem()
		kind      = valueType.Kind()
		err       error
	)

	abstractType := valueType.String()
	if fc, exists := v.knownForceCast[abstractType]; exists {
		v.forceCast[operation.Path(path+"*")] = fc
		zeroValue = reflect.ValueOf(fc.ZeroValue())
	} else {
		switch kind { //nolint:exhaustive
		case reflect.Ptr:
			zeroValue = reflect.Zero(valueType.Elem())
		case reflect.Array:
			zeroValue = reflect.New(reflect.ArrayOf(int(valueType.Size()), reflect.TypeOf(valueType))).Elem()
		case reflect.Slice:
			zeroValue = reflect.MakeSlice(valueType, 0, 0)
		case reflect.Map:
			zeroValue = reflect.MakeMap(valueType)
		default:
			zeroValue = reflect.Zero(valueType)
		}

		kind = zeroValue.Kind()
	}

	if err = v.addRule(path+"*", kind, "", zeroValue.Interface(), &inheritedRules); err != nil {
		return err
	}

	if err := v.parseReference(zeroValue.Type(), path+"*", inheritedRules); err != nil {
		return err
	}

	return nil
}

func (v *Validator) parseReferenceStruct(
	objectType reflect.Type, path string, inheritedRules map[string]rule.Rule,
) error {
	for i := 0; i < objectType.NumField(); i++ {
		var (
			zeroValue reflect.Value
			field     = objectType.Field(i)
			bsonName  = field.Tag.Get("bson")
			kind      = field.Type.Kind()
			err       error
		)

		if bsonName == "" || bsonName == "-" {
			continue
		}

		abstractType := ""
		if kind == reflect.Ptr {
			abstractType = field.Type.Elem().String()
		} else {
			abstractType = field.Type.String()
		}

		if fc, exists := v.knownForceCast[abstractType]; exists {
			v.forceCast[operation.Path(path+bsonName)] = fc
			zeroValue = reflect.ValueOf(fc.ZeroValue())
		} else {
			switch kind { //nolint:exhaustive
			case reflect.Ptr:
				zeroValue = reflect.Zero(field.Type.Elem())
			case reflect.Array:
				zeroValue = reflect.New(reflect.ArrayOf(int(field.Type.Size()), reflect.TypeOf(field.Type))).Elem()
			case reflect.Slice:
				zeroValue = reflect.MakeSlice(field.Type, 0, 0)
			case reflect.Map:
				zeroValue = reflect.MakeMap(field.Type)
			default:
				zeroValue = reflect.Zero(field.Type)
			}

			kind = zeroValue.Kind()
		}

		if err = v.addRule(path+bsonName, kind, field.Tag, zeroValue.Interface(), &inheritedRules); err != nil {
			return err
		}

		if kind == reflect.Struct ||
			kind == reflect.Slice ||
			kind == reflect.Array ||
			kind == reflect.Ptr ||
			kind == reflect.Map {
			if err := v.parseReference(zeroValue.Type(), path+bsonName, inheritedRules); err != nil {
				return err
			}
		}
	}

	return nil
}

func (v *Validator) addRule(
	path string, kind reflect.Kind, tag reflect.StructTag, reference interface{}, inheritedRules *map[string]rule.Rule,
) error {
	if kind == reflect.Invalid {
		return InvalidTypeError{path: path}
	}

	var (
		err   error
		rules = map[string]rule.Rule{}
	)

	for name, rule := range *inheritedRules {
		rules[name], err = rule.NewInheritInstance(path, kind, reference)
		if err != nil {
			return fmt.Errorf("could not inherit rule on '%s': %w", path, err)
		}
	}

	for name, rule := range v.generalRules {
		rules[name], err = rule.NewInstance(path, kind, reference, "")
		if err != nil {
			return fmt.Errorf("could not instantiate rule on '%s': %w", path, err)
		}
	}

	if tag != "" {
		err := v.composeRulesFromTags(path, kind, reference, tag, &rules, inheritedRules)
		if err != nil {
			return err
		}
	}

	if strings.Contains(path, "*") {
		v.wildcardRules[operation.Path(path)] = rules
	} else {
		v.rules[operation.Path(path)] = rules
	}

	return nil
}

func (v *Validator) composeRulesFromTags(
	path string, kind reflect.Kind, reference interface{},
	tag reflect.StructTag, rules, inheritedRules *map[string]rule.Rule,
) error {
	tagsToInherit := []string{}

	if value := tag.Get("jp_inherit"); value != "" {
		split := strings.Split(value, ",")

		for _, tagToInherit := range split {
			if tag.Get(tagToInherit) == "" {
				return InheritNonExistingTagError{name: tagToInherit}
			}

			if _, exist := v.knownTagRules[tagToInherit]; !exist {
				return UnknownRuleError{name: tagToInherit}
			}

			tagsToInherit = append(tagsToInherit, tagToInherit)
		}
	}

	for name, rule := range v.knownTagRules {
		if name == "jp_inherit" {
			continue
		}

		value, use := tag.Lookup(name)
		if use {
			newRule, err := rule.NewInstance(path, kind, reference, value)
			if err != nil {
				return fmt.Errorf("could not instantiate rule on '%s': %w", path, err)
			}

			(*rules)[name] = newRule

			for _, tagToInherit := range tagsToInherit {
				if tagToInherit == name {
					(*inheritedRules)[name] = newRule

					break
				}
			}
		}
	}

	return nil
}

// NewValidator create a new instance of validator using given reference.
func NewValidator(referenceType reflect.Type) (*Validator, error) {
	validator := Validator{
		knownForceCast: map[string]forcecast.ForceCast{
			"primitive.ObjectID":   forcecast.ObjectIDCast{},
			"[]primitive.ObjectID": forcecast.ObjectIDArrayCast{},
		},
		forceCast: map[operation.Path]forcecast.ForceCast{},
		generalRules: map[string]rule.Rule{
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{},
			"jp_general_matching_kind":              &rule.MatchingKindRule{},
		},
		knownTagRules: map[string]rule.Rule{
			"jp_disallow":      &rule.DisallowRule{},
			"jp_min":           &rule.MinRule{},
			"jp_max":           &rule.MaxRule{},
			"jp_expression":    &rule.ExpressionRule{},
			"jp_op_allowed":    &rule.AllowedOperationsRule{},
			"jp_op_disallowed": &rule.DisallowedOperationsRule{},
		},
		rules:         map[operation.Path]map[string]rule.Rule{},
		wildcardRules: map[operation.Path]map[string]rule.Rule{},
	}

	err := validator.UseReference(referenceType)
	if err != nil {
		return nil, fmt.Errorf("initializing validator failed: %w", err)
	}

	return &validator, nil
}
