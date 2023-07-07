package cmd

import (
	"bufio"
	"fmt"
	"math/big"
	"net"
	"os"
	"strings"
)

type IPPool struct {
	ipNet *net.IPNet
}

func ReadIPAddressesFromFile(config Configs) ([]string, error) {
	filepath := config.InputFilepath
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	scanner := bufio.NewScanner(file)
	var ipAddresses []string
	for scanner.Scan() {
		ip := strings.TrimSpace(scanner.Text())
		if ip != "" {
			ipAddresses = append(ipAddresses, ip)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return ipAddresses, nil
}

func IPRange(ipNet *net.IPNet) (net.IP, net.IP) {
	start := ipNet.IP
	mask := ipNet.Mask
	bcst := make(net.IP, len(ipNet.IP))
	copy(bcst, ipNet.IP)
	for i := 0; i < len(mask); i++ {
		ipIdx := len(bcst) - i - 1
		bcst[ipIdx] = ipNet.IP[ipIdx] | ^mask[len(mask)-i-1]
	}
	return start, bcst
}

func InTarget(ip string) bool {
	ipParsed := net.ParseIP(ip)
	for _, ipPool := range IPNetPools {
		if ipPool.ipNet.Contains(ipParsed) {
			return true
		}
	}
	return false
}

func GetPools(ipPoolFuncList []func() string) func() string {
	return func() string {
		for _, ipPoolFunc := range ipPoolFuncList {
			//TODO maybe have bug, but look well
			for ip := ipPoolFunc(); ip != ""; ip = ipPoolFunc() {
				return ip
			}
		}
		return ""
	}
}

func (ipPool *IPPool) SetPool(cidrIp string) error {
	_, ipNet, err := net.ParseCIDR(cidrIp)
	if err != nil {
		fmt.Println("[-]Read ip fail", err)
		return err
	}
	ipPool.ipNet = ipNet
	return nil
}

func (ipPool IPPool) GetPool() func() string {
	var counter int64 = 0
	start, end := IPRange(ipPool.ipNet)
	return func() string {
		startIpInt := big.NewInt(0)
		startIpInt.SetBytes(start.To4())
		endIpInt := big.NewInt(0)
		endIpInt.SetBytes(end.To4())

		nowIpInt := startIpInt.Int64() + counter
		counter++
		if nowIpInt > endIpInt.Int64() {
			return ""
		}
		return fmt.Sprintf("%d.%d.%d.%d", byte(nowIpInt>>24), byte(nowIpInt>>16), byte(nowIpInt>>8), byte(nowIpInt))
	}
}

func (ipPool IPPool) GetPoolSize() int64 {
	start, end := IPRange(ipPool.ipNet)
	startIpInt := big.NewInt(0)
	startIpInt.SetBytes(start.To4())
	endIpInt := big.NewInt(0)
	endIpInt.SetBytes(end.To4())
	return endIpInt.Int64() - startIpInt.Int64()
}
