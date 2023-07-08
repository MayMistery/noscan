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

func Scan() error {
	err := InitTarget(cmd.Config)
	if err != nil {
		return err
	}

	scanPool := ScanPool()
	go scanPool.Run()
	for host := cmd.IPPools(); host != ""; host = cmd.IPPools() {
		scanPool.Push(host)
	}

	return nil
}

func ScanPool() *utils.Pool {
	scanPool := utils.NewPool(cmd.Config.Threads/4 + 1)
	portScanPool := PortScanPool()
	go portScanPool.Run()
	scanPool.Function = func(input interface{}) {
		host := input.(string)
		if CheckLive(host) {
			//TODO port parser
			for _, port := range cmd.Ports {
				portScanPool.Push(Address{net.ParseIP(host), port})
			}
		}
	}
	return scanPool
}
