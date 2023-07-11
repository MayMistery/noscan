package scan

import (
	"fmt"
	"github.com/MayMistery/noscan/cmd"
	"sync"
	"testing"
)

func TestCheckLive(t *testing.T) {
	cmd.Config.InputFilepath = "../data/target"
	err := InitTarget()
	if err != nil {
		return
	}

	limiter := make(chan struct{}, 10)

	wg := sync.WaitGroup{}
	for host := cmd.IPPools(); host != ""; host = cmd.IPPools() {
		wg.Add(1)
		limiter <- struct{}{}
		go func(host string) {
			if CheckLive(host) {
				fmt.Println(host)
			}
			wg.Done()
			<-limiter
		}(host)

	}

	wg.Wait()
}
