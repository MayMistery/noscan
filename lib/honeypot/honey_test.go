package honeypot

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"
	"testing"
)

func TestFindAll(t *testing.T) {
	_, ex, _, _ := runtime.Caller(0)
	fmt.Println(ex)
	exPath := filepath.Dir(filepath.Dir(filepath.Dir(ex)))
	fmt.Println(exPath)
	inputJson := path.Join(exPath, "result/final3.json")
	fmt.Println(inputJson)
	FindAll(inputJson)
}
