package cmd

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	inputFilepath := "../data/target"
	cidrIPs, err := ReadIPAddressesFromFile(inputFilepath)
	if err != nil {
		fmt.Println("[-]Read target ip fail")
	}

	var ipPool IPPool
	var ipPoolsFuncList []func() string
	for _, cidrIp := range cidrIPs {
		err := ipPool.SetPool(cidrIp)
		if err != nil {
		}
		IPPoolsSize += ipPool.GetPoolSize()
		IPNetPools = append(IPNetPools, ipPool)
		ipPoolsFuncList = append(ipPoolsFuncList, ipPool.GetPool())
	}
	IPPools = GetPools(ipPoolsFuncList)

	for host := IPPools(); host != ""; host = IPPools() {
		fmt.Println(host)
	}
}
