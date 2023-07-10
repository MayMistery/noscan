package scan

import (
	"fmt"
	"github.com/MayMistery/noscan/cmd"
	"net"
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

func InitTarget(cfg cmd.Configs) error {
	cidrIPs, err := cmd.ReadIPAddressesFromFile(cfg)
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

func InitScanner() {
	Scanner = ScannerPool()
	PortScanner = PortScanPool()
	HttpScanner = HttpScanPool()

	go Scanner.Run()
	go PortScanner.Run()
	go HttpScanner.Run()
}

func Scan() error {
	err := InitTarget(cmd.Config)
	InitScanner()
	if err != nil {
		cmd.ErrLog("InitTarget error %v", err)
		return err
	}

	for host := cmd.IPPools(); host != ""; host = cmd.IPPools() {
		Scanner.Push(host)
	}

	return nil
}

func ScannerPool() *cmd.Pool {
	scanPool := cmd.NewPool(cmd.Config.Threads/4 + 1)
	scanPool.Function = func(input interface{}) {
		host := input.(string)
		if CheckLive(host) {
			for _, port := range cmd.Ports {
				PortScanner.Push(Address{net.ParseIP(host), port})
			}
		}
	}
	return scanPool
}
