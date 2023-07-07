package scan

import (
	"github.com/MayMistery/noscan/cmd"
	"net"
	"sync"
	"time"
)

func TCPScan(hostslist []string, ports []int, timeout int64) []net.TCPAddr {
	var AliveAddress []net.TCPAddr
	workers := cmd.Config.Threads
	addrs := make(chan net.TCPAddr, len(hostslist)*len(ports))
	results := make(chan net.TCPAddr, len(hostslist)*len(ports))
	var wg sync.WaitGroup

	//Handle results
	go func() {
		for found := range results {
			AliveAddress = append(AliveAddress, found)
			//TODO store the result into database
			wg.Done()
		}
	}()

	//MultiThreaded scanning
	for i := 0; i < workers; i++ {
		go func() {
			for addr := range addrs {
				TCPConnect(addr, results, timeout, &wg)
				wg.Done()
			}
		}()
	}

	//Add the address through channel
	for _, port := range ports {
		for _, host := range hostslist {
			wg.Add(1)
			hostIP := net.ParseIP(host)
			addrs <- net.TCPAddr{IP: hostIP, Port: port}
		}
	}
	wg.Wait()
	close(addrs)
	close(results)
	return AliveAddress
}

func TCPConnect(addr net.TCPAddr, respondingHosts chan<- net.TCPAddr, adjustedTimeout int64, wg *sync.WaitGroup) {
	conn, err := cmd.WrapperTcpWithTimeout("tcp4", addr.String(), time.Duration(adjustedTimeout)*time.Second)
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()
	if err == nil {
		wg.Add(1)
		respondingHosts <- addr
	}
}
