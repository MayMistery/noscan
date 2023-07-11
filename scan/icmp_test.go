package scan

import (
	"fmt"
	"github.com/MayMistery/noscan/cmd"
	"sync"
	"testing"
)

func TestCheckLive(t *testing.T) {
	cmd.Config.InputFilepath = "../data/target"
	//cmd.Config.Ping = true
	err := InitTarget()
	if err != nil {
		return
	}

	limiter := make(chan struct{}, 10000)

	wg := sync.WaitGroup{}
	for host := cmd.IPPools(); host != ""; host = cmd.IPPools() {
		wg.Add(1)
		limiter <- struct{}{} // limit concurrency

		go func(h string) {
			defer wg.Done()              // ensure Done is called even if CheckLive panics
			defer func() { <-limiter }() // ensure limiter is emptied even if CheckLive panics

			if CheckLive(h) {
				fmt.Println(h)
			}
		}(host) // pass a copy of host to the goroutine
	}

	wg.Wait() // wait for all goroutines to finish
}
