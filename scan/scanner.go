package scan

import (
	"fmt"
	"github.com/MayMistery/noscan/cmd"
	"github.com/MayMistery/noscan/utils"
	"net"
)

type Address struct {
	IP   net.IP
	Port int
}

var (
	Scanner     *utils.Pool
	PortScanner *utils.Pool
	HttpScanner *utils.Pool
)

func InitTarget(cfg cmd.Configs) error {
	cidrIPs, err := cmd.ReadIPAddressesFromFile(cfg)
	if err != nil {
		fmt.Println("[-]Read target ip fail")
		return err
	}

	var ipPool cmd.IPPool
	var ipPoolsFuncList []func() string
	for _, cidrIp := range cidrIPs {
		err := ipPool.SetPool(cidrIp)
		if err != nil {
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
	Scanner = ScanPool()
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
		return err
	}

	for host := cmd.IPPools(); host != ""; host = cmd.IPPools() {
		Scanner.Push(host)
	}

	return nil
}

func ScanPool() *utils.Pool {
	scanPool := utils.NewPool(cmd.Config.Threads/4 + 1)
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
