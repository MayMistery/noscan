package rules

import (
	"errors"
	"fmt"
	"github.com/MayMistery/noscan/cmd"
	"regexp"
	"strconv"
	"strings"
)

type Expression struct {
	//表达式数组
	paramSlice []*Param
	//表达式原文
	value string
	//表达式逻辑字符串
	//value = (body="test" || header="tt") && response="aaaa"
	//expr  = (${1} || ${2}) && ${3}
	expr string
}

// parseExpression parses a given expression string and returns an Expression struct.
// It first verifies the characters in the expression, then parses the parameters, and finally verifies the syntax.
func parseExpression(expr string) (*Expression, error) {
	e := &Expression{}
	e.value = expr
	//去除表达式尾部空格
	expr = strings.TrimSpace(expr)
	//修饰expr
	expr = strings.ReplaceAll(expr, `\"`, `[quota]`)
	//字符合法性校验
	if err := exprCharVerification(expr); err != nil {
		cmd.ErrLog("exprCharVerification error %v", err)
		return nil, err
	}
	//提取param数组
	var paramSlice []*Param
	paramRawSlice := paramRegx.FindAllStringSubmatch(expr, -1)
	//对param进行解析
	for index, value := range paramRawSlice {
		expr = strings.Replace(expr, value[0], "${"+strconv.Itoa(index+1)+"}", 1)
		param, err := parseParam(value[0])
		if err != nil {
			cmd.ErrLog("parseParam error %v", err)
			return nil, err
		}
		paramSlice = append(paramSlice, param)
	}
	//语义合法性校验
	if err := exprSyntaxVerification(expr); err != nil {
		cmd.ErrLog("exprSyntaxVerification error %v", err)
		return nil, err
	}
	e.expr = expr
	e.paramSlice = paramSlice
	return e, nil
}

// match checks if the given banner matches the expression.
// It generates a boolean expression using the banner and then parses the boolean expression.
func (e *Expression) match(banner *Banner) bool {
	expr := e.makeBoolExpression(banner)
	b, _ := parseBoolFromString(expr)
	return b
}

// makeBoolExpression generates a boolean expression from the given banner.
// It replaces each parameter in the expression with its boolean value.
func (e *Expression) makeBoolExpression(banner *Banner) string {
	var expr = e.expr
	for index, param := range e.paramSlice {
		b := param.match(banner)
		expr = strings.Replace(expr, "${"+strconv.Itoa(index+1)+"}", strconv.FormatBool(b), 1)
	}
	return expr
}

// Split splits the expression into a slice of subexpressions.
// It uses recursive splitting and then reduces each subexpression.
func (e *Expression) Split() []string {
	r := recursiveSplitExpression(e.expr)
	for i, v := range r {
		r[i] = e.Reduction(v)
	}
	return r
}

// Reduction replaces each parameter placeholder in a subexpression with its original string representation.
func (e *Expression) Reduction(s string) string {
	for i, v := range e.paramSlice {
		param := fmt.Sprintf("${%d}", i+1)
		s = strings.ReplaceAll(s, param, v.String())
	}
	return s
}

// exprCharVerification checks for any unknown characters in the expression.
// It replaces all parameters and logical characters with empty strings and then checks if any characters remain.
func exprCharVerification(expr string) error {
	//把所有param替换为空
	str := paramRegx.ReplaceAllString(expr, "")
	//把所有逻辑字符替换为空
	str = regexp.MustCompile(`[&| ()]`).ReplaceAllString(str, "")
	//检测是否存在其他字符
	if str != "" {
		str = strings.ReplaceAll(str, `[quota]`, `\"`)
		return errors.New(strconv.Quote(str) + " is unknown")
	}
	//检测语法合法性
	return nil
}

var regxSyntaxVerification = regexp.MustCompile(`\${\d+}`)

// exprSyntaxVerification verifies the syntax of the expression.
// It replaces all parameter placeholders with "true" and then tries to parse the boolean expression.
func exprSyntaxVerification(expr string) error {
	expr = regxSyntaxVerification.ReplaceAllString(expr, "true")
	_, err := parseBoolFromString(expr)
	if err != nil {
		return errors.New(expr + ":" + err.Error())
	}
	return nil
}

// trimParentheses trims the parentheses from the given string if it starts and ends with parentheses.
func trimParentheses(s string) string {
	if s[0:1] != "(" {
		return s
	}
	length := len(s)
	if length == 0 {
		return s
	}
	if s[length-1:length] != ")" {
		return s
	}
	index, err := findCoupleBracketIndex(s)
	if err != nil {
		return s
	}
	if index+1 == length {
		return s[1 : length-1]
	}
	return s
}

// splitExpression splits the given expression into a slice of strings.
// It separates the expression into segments based on parentheses, ANDs, ORs, and parameter placeholders.
func splitExpression(s string) []string {
	var data []string
	for i := 0; i < len(s); i++ {
		char := s[i : i+1]
		if char == "(" {
			length, err := findCoupleBracketIndex(s[i:])
			if err != nil {
				cmd.ErrLog("PANIC!!! findCoupleBracketIndex %v", err)
				panic(err)
				//TODO why panic
			}
			data = append(data, s[i+1:i+length])
			i += length
		}
		if char == "&" {
			data = append(data, "&&")
			i += 1
		}
		if char == "|" {
			data = append(data, "||")
			i += 1
		}
		if char == "$" {
			var j = 0
			for j = 0; j < len(s)-i; j++ {
				index := i + j
				if s[index:index+1] == "}" {
					break
				}
			}
			data = append(data, s[i:i+j+1])
			i += j
		}
	}
	var r []string
	for i := 0; i < len(data); i++ {
		var s = data[i]
		if s == "||" {
			r = append(r, data[i+1])
			i += 1
			continue
		}
		if s == "&&" {
			if len(r) == 0 {
				cmd.ErrLog("PANIC!!! expression invalid")
				//TODO why panic
				panic("expression invalid")
			}
			for j := 1; j < len(r); j++ {
				r[0] = r[0] + " || " + r[j]
			}
			r[0] = r[0] + " && " + data[i+1]
			r = r[:1]
			i += 1
			continue
		}
		r = append(r, data[i])
	}
	return r
}

// recursiveSplitExpression recursively splits the given expression into a slice of strings.
// It trims parentheses from the expression and then splits it into subexpressions.
func recursiveSplitExpression(s string) []string {
	var result []string
	s = trimParentheses(s)
	for _, str := range splitExpression(s) {
		sr := splitExpression(s)
		if len(sr) == 1 {
			result = append(result, sr[0])
			continue
		}
		result = append(result, recursiveSplitExpression(str)...)
	}
	return result
}
