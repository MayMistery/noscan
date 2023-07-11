package scan

import (
	"fmt"
	"github.com/MayMistery/noscan/cmd"
	"testing"
)

func TestInitTarget(t *testing.T) {
	cmd.Config.InputFilepath = "../data/target"
	InitTarget()
	for i := 0; i < 513; i++ {
		fmt.Println(cmd.IPPools())
	}
}
