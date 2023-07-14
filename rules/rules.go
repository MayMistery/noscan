package rules

import (
	"errors"
	"github.com/MayMistery/noscan/cmd"
	"io"
	"os"
	"strings"
)

type FingerPrint struct {
	//指纹适配产品
	ProductName string
	//指纹识别规则
	Rule *Expression
}

var FingerPrints []*FingerPrint

func addFingerPrint(productName string, ruleExpr string) error {
	expression, err := parseExpression(ruleExpr)
	if err != nil {
		cmd.ErrLog("Parse expression error %s %s", productName, ruleExpr)
		return err
	}

	httpFinger := &FingerPrint{
		ProductName: productName,
		Rule:        expression,
	}
	FingerPrints = append(FingerPrints, httpFinger)
	return nil
}

func InitFingerPrints(path string) (n int, lastErr error) {
	fs, err := os.Open(path)
	if err != nil {
		cmd.ErrLog("Open fingerprints error %s %v", path, err)
		return 0, err
	}
	return InitFingerPrintsFS(fs)
}

func InitFingerPrintsFS(fs io.Reader) (n int, lastErr error) {
	sourceBuf, err := io.ReadAll(fs)
	if err != nil {
		cmd.ErrLog("Open fingerprintsFs error %v", err)
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
			cmd.ErrLog("add fingerprintsFs error %v", err)
			lastErr = err
		}
	}
	return len(source), lastErr
}

func Clear() {
	FingerPrints = []*FingerPrint{}
}
