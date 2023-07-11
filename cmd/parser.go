package cmd

import (
	"bufio"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"strconv"
	"strings"
)

type IPPool struct {
	ipNet *net.IPNet
}

func ReadIPAddressesFromFile(filepath string) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Printf("Fail to open the target file %v", err)
		ErrLog("Fail to open the target file %v", err)
		return nil, err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Fail to close the target file %v", err)
			ErrLog("Fail to close the target file %v", err)
		}
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
		log.Printf("Fail to read the target file %v", err)
		ErrLog("Fail to read the target file %v", err)
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
			if ip := ipPoolFunc(); ip != "" {
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

func ParsePort() (scanPorts []int) {
	if Config.Ports == "" {
		return
	}
	switch Config.Ports {
	case "all":
		for i := 1; i < 65536; i++ {
			scanPorts = append(scanPorts, i)
		}
		return scanPorts
		//case "top1000":
		//TODO set top 1000 ports
		//case "top10000":
		//TODO set top 10000 ports
	default:
		slices := strings.Split(Config.Ports, ",")
		for _, port := range slices {
			port = strings.TrimSpace(port)
			if port == "" {
				continue
			}
			upper := port
			if strings.Contains(port, "-") {
				ranges := strings.Split(port, "-")
				if len(ranges) < 2 {
					continue
				}

				startPort, _ := strconv.Atoi(ranges[0])
				endPort, _ := strconv.Atoi(ranges[1])
				if startPort < endPort {
					port = ranges[0]
					upper = ranges[1]
				} else {
					port = ranges[1]
					upper = ranges[0]
				}
			}
			start, _ := strconv.Atoi(port)
			end, _ := strconv.Atoi(upper)
			for i := start; i <= end; i++ {
				scanPorts = append(scanPorts, i)
			}
		}
		scanPorts = removeDuplicate(scanPorts)
	}
	return scanPorts
}

func removeDuplicate(old []int) []int {
	var result []int
	temp := map[int]struct{}{}
	for _, item := range old {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
