package scan

import (
	"fmt"
	"github.com/MayMistery/noscan/cmd"
	"github.com/MayMistery/noscan/lib/simplenet"
	"net"
	"sync"
	"time"
)

type Address struct {
	IP   net.IP
	Port int
}

var (
	Scanner     *cmd.Pool
	PortScanner *cmd.Pool
	HttpScanner *cmd.Pool
)

func InitTarget() error {
	cidrIPs, err := cmd.ReadIPAddressesFromFile(cmd.Config.InputFilepath)
	if err != nil {
		cmd.ErrLog("Read target ip fail %v", err)
		fmt.Println("[-]Read target ip fail")
		return err
	}

	var ipPool cmd.IPPool
	var ipPoolsFuncList []func() string
	for _, cidrIp := range cidrIPs {
		err := ipPool.SetPool(cidrIp)
		if err != nil {
			cmd.ErrLog("SetPool fail %v", err)
			return err
		}
		cmd.IPPoolsSize += ipPool.GetPoolSize()
		cmd.IPNetPools = append(cmd.IPNetPools, ipPool)
		ipPoolsFuncList = append(ipPoolsFuncList, ipPool.GetPool())
	}

	cmd.IPPools = cmd.GetPools(ipPoolsFuncList)

	//Init ports array
	cmd.Ports = cmd.ParsePort()
	return nil
}

func InitScanner() *sync.WaitGroup {
	wg := &sync.WaitGroup{}
	conn, err := net.DialTimeout("ip4:icmp", "127.0.0.1", cmd.Config.Timeout)
	if err != nil {
		cmd.Config.Ping = false
	}
	defer func() {
		if conn != nil {
			err := conn.Close()
			if err != nil {
				cmd.ErrLog("connect close error %v", err)
				return
			}
		}
	}()

	Scanner = ScannerPool()
	PortScanner = PortScanPool()
	HttpScanner = HttpScanPool()

	wg.Add(3)
	go Scanner.Run()
	go PortScanner.Run()
	go HttpScanner.Run()
	return wg
}

func StopScanner(wg *sync.WaitGroup) {
	for {
		time.Sleep(cmd.Config.Timeout * 3)
		if Scanner.RunningThreads() == 0 && Scanner.Done == false {
			Scanner.Stop()
			wg.Done()
		}
		if PortScanner.RunningThreads() == 0 && PortScanner.Done == false {
			PortScanner.Stop()
			wg.Done()
		}
		if HttpScanner.RunningThreads() == 0 && HttpScanner.Done == false {
			HttpScanner.Stop()
			wg.Done()
		}
	}
}

func Scan() error {
	err := InitTarget()
	wg := InitScanner()
	if err != nil {
		cmd.ErrLog("InitTarget error %v", err)
		return err
	}

	for host := cmd.IPPools(); host != ""; host = cmd.IPPools() {
		Scanner.Push(host)
	}

	go StopScanner(wg)
	wg.Wait()
	return nil
}

func ScannerPool() *cmd.Pool {
	scanPool := cmd.NewPool(cmd.Config.Threads/20 + 1)
	scanPool.Function = func(input interface{}) {
		host := input.(string)
		if CheckIcmpLive(host) {
			//cmd.ResultLog("[+]%s is alive", host)
			HandleAliveHost(host)
		} else if CheckCommonPort(host) {
			HandleAliveHost(host)
		}

	}
	return scanPool
}

func HandleAliveHost(host string) {
	for _, port := range cmd.Ports {
		PortScanner.Push(Address{net.ParseIP(host), port})
	}
}

func CheckIcmpLive(host string) bool {
	if cmd.Config.Ping == true {
		//Use system ping command
		return ExecCommandPing(host)
	} else {
		return IcmpAlive(host)
	}
}

func CheckCommonPort(host string) bool {
	var commonPorts = []int{22, 23, 80, 139, 512, 443, 445, 3389}
	for _, port := range commonPorts {
		addr := fmt.Sprintf("%s:%d", host, port)
		_, err := simplenet.Send("tcp", false, addr, "\r\n", cmd.Config.Timeout, 128)
		if err == nil {
			return true
		}
	}
	return false
}
