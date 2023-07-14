package rules

import (
	"errors"
	"fmt"
	"github.com/MayMistery/noscan/cmd"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Param struct {
	keyword  string
	value    string
	operator Operator
}
type Operator string

const (
	unequal    Operator = "!=" // !=
	equal               = "="  // =
	regxEqual           = "~=" // ~=
	superEqual          = "==" // ==
)

var keywordSlice = []string{
	"Title",
	"Header",
	"Body",
	"Response",
	"Protocol",
	"Cert",
	"Port",
	"Hash",
	"Icon",
}

var paramRegx = regexp.MustCompile(`([a-zA-Z0-9]+) *(!=|=|~=|==) *"([^"\n]+)"`)
var keywordRegx = regexp.MustCompile("^" + strings.Join(keywordSlice, "|") + "$")

// parseParam parses a given parameter string and returns a Param struct.
// It verifies the keyword, converts the operator, and processes the value.
func parseParam(expr string) (*Param, error) {
	p := paramRegx.FindStringSubmatch(expr)

	keyword := p[1]
	valueRaw := p[3]

	keyword = strings.ToUpper(keyword[:1]) + keyword[1:]

	if keywordRegx.MatchString(keyword) == false {
		return nil, errors.New(keyword + " keyword is unknown")
	}

	operator := convOperator(p[2])

	valueRaw = strings.ReplaceAll(valueRaw, `[quota]`, `\"`)
	value, err := strconv.Unquote("\"" + valueRaw + "\"")
	if err != nil {
		cmd.ErrLog("parseParam strconv.Unquote error")
		return nil, err
	}
	if operator == regxEqual {
		_, err = regexp.Compile(value)
		if err != nil {
			cmd.ErrLog("parseParam regexp.Compile error")
			return nil, err
		}
	}

	return &Param{
		keyword:  keyword,
		value:    value,
		operator: operator,
	}, nil
}

// match checks if the given banner matches the parameter.
// It gets the value of the field in the banner that matches the keyword and then checks if the value matches the parameter based on the operator.
func (p *Param) match(banner *Banner) bool {
	subStr := p.value
	keyword := p.keyword

	v := reflect.ValueOf(*banner)
	str := v.FieldByName(keyword).String()

	switch p.operator {
	case unequal:
		return !strings.Contains(str, subStr)
	case equal:
		return strings.Contains(str, subStr)
	case regxEqual:
		return regexp.MustCompile(subStr).MatchString(str)
	case superEqual:
		return str == subStr
	default:
		return false
	}
}

func (p *Param) String() string {
	return fmt.Sprintf("%s%s%s", p.keyword, p.operator, strconv.Quote(p.value))
}

const (
	and = iota // &&
	or         // ||
)

// convOperator converts a given operator string to an Operator.
// It returns the corresponding Operator for the given string, or panics if the string is not recognized.
func convOperator(expr string) Operator {
	switch expr {
	case "!=":
		return unequal
	case "=":
		return equal
	case "~=":
		return regxEqual
	case "==":
		return superEqual
	default:
		panic(expr)
	}
}

// parseBoolFromString parses a given boolean expression string and returns a boolean.
// It first verifies the characters in the expression, then parses the expression.
func parseBoolFromString(expr string) (bool, error) {
	//去除空格
	expr = strings.ReplaceAll(expr, " ", "")
	//如果存在其他异常字符，则报错
	s := regexp.MustCompile(`true|false|&|\||\(|\)`).ReplaceAllString(expr, "")
	if s != "" {
		return false, errors.New(s + "is unknown")
	}
	return stringParse(expr)
}

// stringParse parses a given boolean expression string and returns a boolean.
// It goes through the expression character by character, processing true, false, AND, OR, and parentheses.
func stringParse(expr string) (bool, error) {
	first := true
	operator := and
	if expr == "true" {
		return true, nil
	}
	if expr == "false" {
		return false, nil
	}

	for i := 0; i < len(expr); i++ {
		char := expr[i : i+1]
		if char == "t" {
			first = parseCoupleBool(first, true, operator)
			i += 3
		}
		if char == "f" {
			first = parseCoupleBool(first, false, operator)
			i += 4
		}
		if char == "&" {
			operator = and
			i += 1

		}
		if char == "|" {
			operator = or
			i += 1
		}
		if char == "(" {
			length, err := findCoupleBracketIndex(expr[i:])
			if err != nil {
				cmd.ErrLog("findCoupleBracketIndex error")
				return false, err
			}
			next, err := stringParse(expr[i+1 : i+length])
			if err != nil {
				cmd.ErrLog("stringParse error")
				return false, err
			}
			first = parseCoupleBool(first, next, operator)

			i += length
		}

	}
	return first, nil
}

// parseCoupleBool combines two booleans based on the given operator.
// If the operator is OR, it returns the logical OR of the booleans. If the operator is AND, it returns the logical AND of the booleans.
func parseCoupleBool(first bool, next bool, operator int) bool {
	if operator == or {
		return first || next
	}
	if operator == and {
		return first && next
	}
	return false
}

// findCoupleBracketIndex finds the index of the closing parenthesis that matches the first opening parenthesis in the given expression.
// It returns an error if the parentheses are not balanced.
func findCoupleBracketIndex(expr string) (int, error) {
	var leftIndex []int
	var rightIndex []int

	for index, value := range expr {
		if value == '(' {
			leftIndex = append(leftIndex, index)
		}
		if value == ')' {
			rightIndex = append(rightIndex, index)
		}
	}

	if len(leftIndex) != len(rightIndex) {
		return 0, errors.New("bracket is not couple")
	}
	for i, index := range rightIndex {
		countLeft := strings.Count(expr[:index], "(")
		if countLeft == i+1 {
			return index, nil
		}

	}
	return 0, errors.New("bracket is not couple")
}
