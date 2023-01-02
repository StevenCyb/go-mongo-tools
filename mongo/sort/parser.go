package sort

import (
	"strings"

	"github.com/StevenCyb/go-mongo-tools/errs"
	"github.com/StevenCyb/go-mongo-tools/tokenizer"
	"go.mongodb.org/mongo-driver/bson"
)

// Types that are used in this parser.
const (
	SkipType          tokenizer.Type = "SKIP"
	AndType           tokenizer.Type = ","
	SetType           tokenizer.Type = "="
	SortConditionType tokenizer.Type = "SORT_CRITERIA"
	FieldNameType     tokenizer.Type = "FIELD_NAME"
)

// specialEncode is the map for encoding
// a list of special characters.
//
//nolint:gochecknoglobals
var specialEncode = map[string]string{
	`,`: "%5C%2C",
	`=`: "%5C%3D",
	` `: "%20",
}

// NewParser creates a new parser.
func NewParser(policy *tokenizer.Policy) *Parser {
	return &Parser{
		policy: policy,
	}
}

// Parser provides the logic to parse rsql statements.
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
			tokenizer.NewSpec(`^,`, AndType),
			tokenizer.NewSpec(`^(=)`, SetType),
			tokenizer.NewSpec(`^(asc|desc|1|-1)`, SortConditionType),
			tokenizer.NewSpec(`^[^=]*`, FieldNameType),
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
 *   | <sort_statement>
 *   | <sort_statement> "," <sort_statement>
 * .
 */
func (p *Parser) expression() ([]bson.E, error) {
	sortStatements := []bson.E{}

	if p.lookahead == nil {
		return nil, errs.NewErrUnexpectedInputEnd(FieldNameType.String())
	}

	sortStatement, err := p.sortStatement()
	if err != nil {
		return nil, err
	}

	sortStatements = append(sortStatements, *sortStatement)

	if p.lookahead != nil {
		_, err := p.eat(AndType)
		if err != nil {
			return nil, err
		}

		nextSortStatements, err := p.expression()
		if err != nil {
			return nil, err
		}

		sortStatements = append(sortStatements, nextSortStatements...)
	}

	return sortStatements, nil
}

/*
 * <sort_statement>
 *   : <key> "=" <sort_condition>
 * .
 */
func (p *Parser) sortStatement() (*bson.E, error) {
	keyToken, err := p.eat(FieldNameType)
	if err != nil {
		return nil, err
	}

	_, err = p.eat(SetType)
	if err != nil {
		return nil, err
	}

	sortConditionToken, err := p.eat(SortConditionType)
	if err != nil {
		return nil, err
	}

	var sort int
	if sortConditionToken.Value == "asc" || sortConditionToken.Value == "1" {
		sort = 1
	} else {
		sort = -1
	}

	return &bson.E{Key: keyToken.Value, Value: sort}, nil
}
