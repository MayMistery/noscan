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
	Start net.IP
	End   net.IP
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

// Todo trans target CIDR to list
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

func (ipPool *IPPool) SetPool(cidrIp string) {
	_, ipNet, err := net.ParseCIDR(cidrIp)
	if err != nil {
		fmt.Println("[-]Read ip fail", err)
		return
	}
	start, end := IPRange(ipNet)
	ipPool.Start = start
	ipPool.End = end
}

func (ipPool IPPool) GetPool() func() string {
	var counter int64 = 0
	return func() string {
		startIpInt := big.NewInt(0)
		startIpInt.SetBytes(ipPool.Start.To4())
		endIpInt := big.NewInt(0)
		endIpInt.SetBytes(ipPool.End.To4())

		nowIpInt := startIpInt.Int64() + counter
		counter++
		if nowIpInt > endIpInt.Int64() {
			return ""
		}
		return fmt.Sprintf("%d.%d.%d.%d", byte(nowIpInt>>24), byte(nowIpInt>>16), byte(nowIpInt>>8), byte(nowIpInt))
	}
}
