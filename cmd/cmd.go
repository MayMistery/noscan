package cmd

import (
	"github.com/MayMistery/noscan/scan"
	"github.com/MayMistery/noscan/utils"
	"runtime"
	"runtime/debug"
	"time"
)

func init() {
	go func() {
		for {
			GarbageCollection()
			time.Sleep(10 * time.Second)
		}
	}()
}

func GarbageCollection() {
	runtime.GC()
	debug.FreeOSMemory()
}

func Exec() {
	var cfg Configs

	Flag(cfg)
	scan.Scan(cfg)
	if cfg.jsonOutput == true {
		utils.OutputJsonResult(cfg.CIDRInfo)
	}
}
