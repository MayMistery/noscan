package scan

import (
	"fmt"
	"github.com/MayMistery/noscan/cmd"
	"testing"
)

func TestInitTarget(t *testing.T) {
	var cfg cmd.Configs
	cfg.InputFilepath = "../data/target"
	InitTarget(cfg)
	for i := 0; i < 513; i++ {
		fmt.Println(cmd.IpPools())
	}
}
