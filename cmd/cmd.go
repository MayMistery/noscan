package cmd

import (
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
	Flag()
	//Parse()
	//scan.Scan(CIDRInfo)
	//utils.
}
