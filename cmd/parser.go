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
	case "common":
		Config.Ports = "1,7,9,13,19,21-23,25,37,42,49,53,69,79-81,85,105,109-111,113,123,135,137-139,143,161,179,222,264,384,389,402,407,443-446,465,500,502,512-515,523-524,540,548,554,587,617,623,689,705,771,783,873,888,902,910,912,921,993,995,998,1000,1024,1030,1035,1090,1098-1103,1128-1129,1158,1199,1211,1220,1234,1241,1300,1311,1352,1433-1435,1440,1494,1521,1530,1533,1581-1582,1604,1720,1723,1755,1811,1900,2000-2001,2049,2082,2083,2100,2103,2121,2199,2207,2222,2323,2362,2375,2380-2381,2525,2533,2598,2601,2604,2638,2809,2947,2967,3000,3037,3050,3057,3128,3200,3217,3273,3299,3306,3311,3312,3389,3460,3500,3628,3632,3690,3780,3790,3817,4000,4322,4433,4444-4445,4659,4679,4848,5000,5038,5040,5051,5060-5061,5093,5168,5247,5250,5351,5353,5355,5400,5405,5432-5433,5498,5520-5521,5554-5555,5560,5580,5601,5631-5632,5666,5800,5814,5900-5910,5920,5984-5986,6000,6050,6060,6070,6080,6082,6101,6106,6112,6262,6379,6405,6502-6504,6542,6660-6661,6667,6690,6905,6988,7001,7021,7071,7080,7144,7181,7210,7443,7510,7579-7580,7700,7770,7777-7778,7787,7800-7801,7879,7902,8000-8001,8008,8014,8020,8023,8028,8030,8080-8082,8087,8090,8095,8161,8180,8205,8222,8300,8303,8333,8400,8443-8444,8503,8800,8812,8834,8880,8888-8890,8899,8901-8903,9000,9002,9060,9080-9081,9084,9090,9099-9100,9111,9152,9200,9390-9391,9443,9495,9809-9815,9855,9999-10001,10008,10050-10051,10080,10098,10162,10202-10203,10443,10616,10628,11000,11099,11211,11234,11333,12174,12203,12221,12345,12397,12401,13364,13500,13838,14330,15200,16102,17185,17200,18881,19300,19810,20010,20031,20034,20101,20111,20171,20222,22222,23472,23791,23943,25000,25025,26000,26122,27000,27017,27888,28222,28784,30000,30718,31001,31099,32764,32913,34205,34443,37718,38080,38292,40007,41025,41080,41523-41524,44334,44818,45230,46823-46824,47001-47002,48899,49152,50000-50004,50013,50500-50504,52302,55553,57772,62078,62514,65535"
		return parsePortStr(Config.Ports)
	default:
		return parsePortStr(Config.Ports)
	}
}

func parsePortStr(inputPorts string) (scanPorts []int) {
	slices := strings.Split(inputPorts, ",")
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
