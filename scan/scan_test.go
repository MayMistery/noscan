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
	for _, ipPoolFunc := range cmd.IpPools {
		fmt.Println(ipPoolFunc())
		fmt.Println(ipPoolFunc())
		fmt.Println(ipPoolFunc())
	}
}
