package rsql

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/StevenCyb/go-mongo-tools/errs"
	"github.com/StevenCyb/go-mongo-tools/tokenizer"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Types that are used in this parser.
const (
	SkipType                             tokenizer.Type = "SKIP"
	AndCompositeType                     tokenizer.Type = ";"
	OrCompositeType                      tokenizer.Type = ","
	ContextStartType                     tokenizer.Type = "("
	ContextEndType                       tokenizer.Type = ")"
	ValueCompareOperatorType             tokenizer.Type = "VALUE_COMPARE_OPERATOR"
	QuotedStringValueCompareOperatorType tokenizer.Type = "QUOTED_STRING_VALUE_COMPARE_OPERATOR"
	NumericValueCompareOperatorType      tokenizer.Type = "NUMERIC_VALUE_COMPARE_OPERATOR"
	ArrayCompareOperatorType             tokenizer.Type = "ARRAY_COMPARE_OPERATOR"
	BoolLiteralType                      tokenizer.Type = "BOOL_LITERAL"
	QuotedStringLiteralType              tokenizer.Type = "QUOTED_STRING_LITERAL"
	OidLiteralType                       tokenizer.Type = "OID_LITERAL"
	FieldNameType                        tokenizer.Type = "FIELD_NAME"
	NumberLiteralType                    tokenizer.Type = "NUMERIC_LITERAL"

	intBase     = 10
	int64Size   = 64
	float64Size = 64
)

// specialEncode is the map for encoding
// a list of special characters.
//
//nolint:gochecknoglobals
var specialEncode = map[string]string{
	`,`: "%5C%2C",
	`;`: "%5C%3B",
	`=`: "%5C%3D",
	`"`: "%22",
	`'`: "%27",
	`$`: "%24",
	` `: "%20",
}

// NewParser creates a new parser.
func NewParser(policy *tokenizer.Policy) *Parser {
	return &Parser{
		policy: policy,
	}
}

// Parser provides the logic to parse
// rsql statements.
type Parser struct {
	tokenizer *tokenizer.Tokenizer
	lookahead *tokenizer.Token
	policy    *tokenizer.Policy
}

// eat return a token with expected type.
func (p *Parser) eat(tokenType tokenizer.Type) (*tokenizer.Token, error) {
	token := p.lookahead

	if token == nil {
		return nil, errs.NewErrUnexpectedInputEnd(tokenType.String())
	}

	if token.Type != tokenType {
		return nil, errs.NewErrUnexpectedTokenType(
			p.tokenizer.GetCursorPosition(),
			token.Type.String(),
			tokenType.String(),
		)
	}

	var err error
	p.lookahead, err = p.tokenizer.GetNextToken()

	return token, err //nolint:wrapcheck
}

// Parse a given query.
func (p *Parser) Parse(query string) (bson.D, error) {
	var err error

	if query == "" {
		return bson.D{}, nil
	}

	for dec, enc := range specialEncode {
		query = strings.ReplaceAll(query, enc, dec)
	}

	p.tokenizer = tokenizer.NewTokenizer(
		query,
		SkipType, FieldNameType,
		[]*tokenizer.Spec{
			tokenizer.NewSpec(`^\s+`, SkipType),
			tokenizer.NewSpec(`^\(`, ContextStartType),
			tokenizer.NewSpec(`^\)`, ContextEndType),
			tokenizer.NewSpec(`^;`, AndCompositeType),
			tokenizer.NewSpec(`^,`, OrCompositeType),
			tokenizer.NewSpec(`^(==|!=)`, ValueCompareOperatorType),
			tokenizer.NewSpec(`^(=sw=|=ew=)`, QuotedStringValueCompareOperatorType),
			tokenizer.NewSpec(`^(=gt=|=ge=|=lt=|=le=)`, NumericValueCompareOperatorType),
			tokenizer.NewSpec(`^(=in=|=out=)`, ArrayCompareOperatorType),
			tokenizer.NewSpec(`^\$oid\([0-9a-fA-F]+\)`, OidLiteralType),
			tokenizer.NewSpec(`(?i)^(true|false)`, BoolLiteralType),
			tokenizer.NewSpec(`^(-|\+)?\d+(\.\d+)?`, NumberLiteralType),
			tokenizer.NewSpec(`^("[^"]*"|'[^']*')`, QuotedStringLiteralType),
			tokenizer.NewSpec(`^[^!=]*`, FieldNameType),
		},
		p.policy,
	)

	p.lookahead, err = p.tokenizer.GetNextToken()
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	return p.expression()
}

/*
 * <expression>
 *   : <context>
 *   | <context> <composite_operator> <expression>
 *   | <comparison>
 *   | <comparison> <composite_operator> <expression>
 * .
 */
func (p *Parser) expression() ([]bson.E, error) { //nolint:funlen
	var (
		left           bson.E
		sortStatements = []bson.E{}
	)

	if p.lookahead == nil {
		return nil, errs.NewErrUnexpectedInputEnd(FieldNameType.String())
	}

	if p.lookahead.Type == ContextStartType {
		tmp, err := p.context()
		if err != nil {
			return nil, err
		}

		left = tmp[0]
	} else {
		tmp, err := p.comparison()
		if err != nil {
			return nil, err
		}
		left = *tmp
	}

	//nolint:nestif
	if p.lookahead != nil && p.lookahead.Type != ContextEndType {
		logicalOperation, err := p.compositeOperation()
		if err != nil {
			return nil, err
		}

		right, err := p.expression()
		if err != nil {
			return nil, err
		}

		if logicalOperation.Type == AndCompositeType {
			if right[0].Key == "$and" {
				newA, ok := right[0].Value.(bson.A)
				if !ok {
					return nil, errs.NewErrUnexpectedInput(right[0].Value)
				}

				newA = append(bson.A{bson.D{left}}, newA...)
				left = bson.E{Key: "$and", Value: newA}
			} else {
				left = bson.E{Key: "$and", Value: bson.A{bson.D{left}, bson.D{right[0]}}}
			}
		} else {
			if right[0].Key == "$or" {
				newA, ok := right[0].Value.(bson.A)
				if !ok {
					return nil, errs.NewErrUnexpectedInput(right[0].Value)
				}

				newA = append(bson.A{bson.D{left}}, newA...)
				left = bson.E{Key: "$or", Value: newA}
			} else {
				left = bson.E{Key: "$or", Value: bson.A{bson.D{left}, bson.D{right[0]}}}
			}
		}
	}

	sortStatements = append(sortStatements, left)

	return sortStatements, nil
}

/*
 * <context>
 *   : "(" <expression> ")"
 * .
 */
func (p *Parser) context() ([]bson.E, error) {
	_, err := p.eat(ContextStartType)
	if err != nil {
		return nil, err
	}

	context, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.eat(ContextEndType)
	if err != nil {
		return nil, err
	}

	return context, nil
}

/*
 * <composite_operator>
 *   : ";"
 *   | ","
 * .
 */
func (p *Parser) compositeOperation() (*tokenizer.Token, error) {
	if p.lookahead.Type == AndCompositeType {
		return p.eat(AndCompositeType)
	} else if p.lookahead.Type == OrCompositeType {
		return p.eat(OrCompositeType)
	}

	return nil, errs.NewErrUnexpectedTokenType(
		p.tokenizer.GetCursorPosition()-len(p.lookahead.Value),
		p.lookahead.Type.String(),
		";/,")
}

/*
 * <array_comparison>
 *   | <plural_operator> "(" <literal_list> ")"
 * .
 */
func (p *Parser) arrayComparison(key string) (*bson.E, error) {
	operator, err := p.eat(ArrayCompareOperatorType)
	if err != nil {
		return nil, err
	}

	_, err = p.eat(ContextStartType)
	if err != nil {
		return nil, err
	}

	literalList, err := p.literalList()
	if err != nil {
		return nil, err
	}

	_, err = p.eat(ContextEndType)
	if err != nil {
		return nil, err
	}

	switch operator.Value {
	case "=in=":
		return &bson.E{Key: key, Value: bson.E{Key: "$in", Value: literalList}}, nil
	case "=out=":
		return &bson.E{Key: key, Value: bson.E{Key: "$nin", Value: literalList}}, nil
	default:
		return nil, errs.NewErrUnexpectedTokenType(
			p.tokenizer.GetCursorPosition()-len(operator.Value),
			operator.Type.String(),
			ArrayCompareOperatorType.String())
	}
}

/*
 * <numeric_value_comparison>
 *   | <singular_string_operator> <quoted_string_literal>
 * .
 */
func (p *Parser) numericValueComparison(key string) (*bson.E, error) {
	operator, err := p.eat(NumericValueCompareOperatorType)
	if err != nil {
		return nil, err
	}

	literal, err := p.numericLiteral()
	if err != nil {
		return nil, err
	}

	switch operator.Value {
	case "=gt=":
		return &bson.E{Key: key, Value: bson.D{bson.E{Key: "$gt", Value: literal}}}, nil
	case "=ge=":
		return &bson.E{Key: key, Value: bson.D{bson.E{Key: "$gte", Value: literal}}}, nil
	case "=lt=":
		return &bson.E{Key: key, Value: bson.D{bson.E{Key: "$lt", Value: literal}}}, nil
	case "=le=":
		return &bson.E{Key: key, Value: bson.D{bson.E{Key: "$lte", Value: literal}}}, nil
	default:
		return nil, errs.NewErrUnexpectedTokenType(
			p.tokenizer.GetCursorPosition()-len(operator.Value),
			operator.Type.String(),
			ValueCompareOperatorType.String())
	}
}

/*
 * <quoted_string_comparison>
 *   | <singular_string_operator> <quoted_string_literal>
 * .
 */
func (p *Parser) quotedStringComparison(key string) (*bson.E, error) {
	operator, err := p.eat(QuotedStringValueCompareOperatorType)
	if err != nil {
		return nil, err
	}

	literal, err := p.stringLiteral()
	if err != nil {
		return nil, err
	}

	switch operator.Value {
	case "=sw=":
		wildcard, err := regexp.Compile("^" + fmt.Sprintf("%v", literal))

		return &bson.E{Key: key, Value: *wildcard}, errors.Wrap(err, "failed to create wildcard expression")
	case "=ew=":
		wildcard, err := regexp.Compile(fmt.Sprintf("%v", literal) + "$")

		return &bson.E{Key: key, Value: *wildcard}, errors.Wrap(err, "failed to create wildcard expression")
	default:
		return nil, errs.NewErrUnexpectedTokenType(
			p.tokenizer.GetCursorPosition()-len(operator.Value),
			operator.Type.String(),
			ValueCompareOperatorType.String())
	}
}

/*
 * <literal_comparison>
 *   : <singular_operator> <literal>
 * .
 */
func (p *Parser) literalComparison(key string) (*bson.E, error) {
	operator, err := p.eat(ValueCompareOperatorType)
	if err != nil {
		return nil, err
	}

	var literal interface{}
	//nolint:nestif
	if p.lookahead.Type == ContextStartType {
		_, err = p.eat(ContextStartType)
		if err != nil {
			return nil, err
		}

		literal, err = p.literalList()
		if err != nil {
			return nil, err
		}

		_, err = p.eat(ContextEndType)
		if err != nil {
			return nil, err
		}
	} else {
		literal, err = p.literal()
		if err != nil {
			return nil, err
		}
	}

	switch operator.Value {
	case "==":
		return &bson.E{Key: key, Value: literal}, nil
	case "!=":
		return &bson.E{Key: key, Value: bson.D{bson.E{Key: "$ne", Value: literal}}}, nil
	default:
		return nil, errs.NewErrUnexpectedTokenType(
			p.tokenizer.GetCursorPosition()-len(operator.Value),
			operator.Type.String(),
			ValueCompareOperatorType.String())
	}
}

/*
 * <comparison>
 *   : TEXT <singular_operator> <literal>
 *   | TEXT <quoted_string_comparison>
 *   | TEXT <numeric_value_comparison>
 *   | TEXT <array_comparison>
 * .
 */
func (p *Parser) comparison() (*bson.E, error) {
	keyToken, err := p.eat(FieldNameType)
	if err != nil {
		return nil, err
	}

	if p.lookahead == nil {
		return nil, errs.NewErrUnexpectedInputEnd(FieldNameType.String())
	}

	key := keyToken.Value

	switch p.lookahead.Type {
	case ValueCompareOperatorType:
		return p.literalComparison(key)
	case QuotedStringValueCompareOperatorType:
		return p.quotedStringComparison(key)
	case NumericValueCompareOperatorType:
		return p.numericValueComparison(key)
	case ArrayCompareOperatorType:
		return p.arrayComparison(key)
	}

	return nil, errs.NewErrUnexpectedToken(
		p.tokenizer.GetCursorPosition()-len(p.lookahead.Value),
		p.lookahead.Value)
}

/*
 * <literal>
 * : <oid_literal>
 * : <bool_literal>
 * | <quoted_string_literal>
 * | <numeric_literal>
 * .
 */
func (p *Parser) literal() (interface{}, error) {
	switch p.lookahead.Type {
	case OidLiteralType:
		token, err := p.eat(OidLiteralType)
		if err != nil {
			return nil, err
		}

		oidHex := strings.TrimPrefix(token.Value, "$oid(")
		oidHex = strings.TrimSuffix(oidHex, ")")

		oid, err := primitive.ObjectIDFromHex(oidHex)
		if err != nil {
			return nil, fmt.Errorf("could not parse $oid '%s': %w", oidHex, err)
		}

		return oid, nil
	case BoolLiteralType:
		token, err := p.eat(BoolLiteralType)
		if err != nil {
			return nil, err
		}

		return strings.ToLower(token.Value) == "true", nil
	case QuotedStringLiteralType:
		return p.stringLiteral()
	case NumberLiteralType:
		return p.numericLiteral()
	}

	return nil, errs.NewErrUnexpectedTokenType(
		p.tokenizer.GetCursorPosition()-len(p.lookahead.Value),
		p.lookahead.Type.String(),
		"LITERAL")
}

/*
 * <quoted_string_literal>
 * : "'" <TEXT> "'"
 * | """ <TEXT> """
 * .
 */
func (p *Parser) stringLiteral() (interface{}, error) {
	token, err := p.eat(QuotedStringLiteralType)
	if err != nil {
		return nil, err
	}

	replacer := strings.NewReplacer(`"`, "", "'", "")

	return replacer.Replace(token.Value), nil
}

/*
 * <numeric_literal>
 * : <INT>
 * | <FLOAT>
 * .
 */
func (p *Parser) numericLiteral() (interface{}, error) {
	token, err := p.eat(NumberLiteralType)
	if err != nil {
		return nil, err
	}

	if strings.Contains(token.Value, ".") {
		var value float64
		value, err = strconv.ParseFloat(token.Value, float64Size)

		return value, errors.Wrap(err, "failed to parse float value")
	}

	var value int64
	value, err = strconv.ParseInt(token.Value, intBase, int64Size)

	return value, errors.Wrap(err, "failed to parse int value")
}

/*
 * <literal_list>
 * : <quoted_string_literal> "," <literal_list>
 * | <numeric_literal> "," <literal_list>
 * .
 */
func (p *Parser) literalList() (bson.A, error) {
	items := bson.A{}

	body, err := p.literal()
	if err != nil {
		return nil, err
	}

	items = append(items, body)

	for p.lookahead.Type == OrCompositeType {
		_, err := p.eat(OrCompositeType)
		if err != nil {
			return nil, err
		}

		body, err = p.literal()
		if err != nil {
			return nil, err
		}

		items = append(items, body)
	}

	return items, nil
}
