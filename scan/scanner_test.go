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

func TestCIDR(t *testing.T) {
	CheckCIDR("127.0.0.1/24,196.1.168.1,196.1.168.1,196.1.168.1,196.1.168.1")
}
