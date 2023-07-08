package rules

import (
	"errors"
	"io"
	"strings"
)

func InitDatabaseFS(fs io.Reader) (n int, lastErr error) {
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
		err := add(r[0], r[1])
		if err != nil {
			lastErr = err
		}
	}
	return len(source), lastErr
}
