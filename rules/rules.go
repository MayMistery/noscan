package rules

import (
	"errors"
	"fmt"
	"github.com/MayMistery/noscan/cmd"
	"io"
	"os"
	"strings"
)

type FingerPrint struct {
	ProductName string
	Rule        *Expression
}

type Expression struct {
	//TODO 根据你的需求解析成你想读取的形式
}

var FingerPrints []*FingerPrint

func InitRulesFromFile() (n int, lastErr error) {
	fs, err := os.Open(cmd.Config.RulesFilePath)
	if err != nil {
		return 0, err
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(fs)
	sourceBuf, err := io.ReadAll(fs)
	if err != nil {
		return 0, err
	}
	source := strings.Split(string(sourceBuf), "\n")
	for _, line := range source {
		line = strings.TrimSpace(line)
		r := strings.SplitAfterN(line, "\t", 2)
		if len(r) != 2 {
			lastErr = errors.New(line + "invalid")
			continue
		}
		err := addFingerPrint(r[0], r[1])
		if err != nil {
			lastErr = err
		}
	}
	return len(source), lastErr
}

func addFingerPrint(product string, ruleExpr string) error {
	rule, err := parseRuleExpr(ruleExpr)
	if err != nil {
		fmt.Println("Rule Expression parse failed")
		return err
	}
	FingerPrints = append(FingerPrints, &FingerPrint{product, rule})
	return nil
}

func parseRuleExpr(ruleExpr string) (*Expression, error) {
	var rule *Expression
	//TODO 将规则串解析成你要的Expression
	return rule, nil
}
