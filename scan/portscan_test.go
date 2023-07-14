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
	cmd.Config.DBFilePath = "../data/database3.db"
	cmd.Config.OutputFilepath = "../result/result3.json"
	cmd.Config.DeepInspection = true
	cmd.Config.Timeout = 10 * time.Second
	bolt.InitAsyncDatabase()
	utils.InitResultMap()

	wg := &sync.WaitGroup{}
	PortScanner = PortScanPool()
	HttpScanner = HttpScanPool()

	wg.Add(2)
	go PortScanner.Run()
	go HttpScanner.Run()

	for i := 1; i < 65536; i++ {
		PortScanner.Push(Address{net.ParseIP("104.248.48.130"), i})
	}
	//PortScanner.Push(Address{net.ParseIP("204.168.173.224"), 22})
	//PortScanner.Push(Address{net.ParseIP("204.168.173.224"), 443})
	//PortScanner.Push(Address{net.ParseIP("204.168.173.224"), 80})

	go func() {
		for {
			time.Sleep(time.Second * 40)
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
