package scan

import (
	"github.com/MayMistery/noscan/cmd"
	"github.com/MayMistery/noscan/storage/bolt"
	"github.com/MayMistery/noscan/utils"
	"net"
	"sync"
	"testing"
	"time"
)

func TestPortScan(t *testing.T) {
	cmd.Config.Threads = 1000
	cmd.Config.DBFilePath = "../data/database.db"
	cmd.Config.OutputFilepath = "../result/result.json"
	cmd.Config.DeepInspection = true
	cmd.Config.Timeout = 3 * time.Second
	bolt.InitAsyncDatabase()
	utils.InitResultMap()

	wg := &sync.WaitGroup{}
	wg.Add(2)
	PortScanner = PortScanPool()
	HttpScanner = HttpScanPool()

	go PortScanner.Run()
	go HttpScanner.Run()

	for i := 1; i < 65536; i++ {
		PortScanner.Push(Address{net.ParseIP("204.168.173.224"), i})
	}

	go func() {
		for {
			time.Sleep(time.Second)
			if PortScanner.RunningThreads() == 0 && PortScanner.Done == false {
				PortScanner.Stop()
				wg.Done()
			}
			if HttpScanner.RunningThreads() == 0 && HttpScanner.Done == false {
				HttpScanner.Stop()
				wg.Done()
			}
		}
	}()
	wg.Wait()
}
