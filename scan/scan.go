package scan

import (
	"fmt"
	"github.com/MayMistery/noscan/cmd"
)

func tcp(CIDRInfo map[string]cmd.IpInfo) {
	//TODO scan implementation details.
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
		ipPool.SetPool(cidrIp)
		ipPoolsFuncList = append(ipPoolsFuncList, ipPool.GetPool())
	}

	cmd.IpPools = cmd.GetPools(ipPoolsFuncList)
	return nil
}

func Scan(config cmd.Configs) {

	ipList := config.CIDRInfo

	switch config.ScanType {
	case "syn":
		break
	//TODO add scan type
	default:
		tcp(ipList)
	}
}

//func CheckAlive()

//type Result struct {
//	IP   string `json:"ip"`
//	Port string `json:"port"`
//}
//
//func checkPort(protocol, hostname string, port int) bool {
//	address := net.JoinHostPort(hostname, strconv.Itoa(port))
//	//fmt.Println("hello", port, hostname)
//	conn, err := net.DialTimeout(protocol, address, 2*time.Second)
//	if err != nil {
//		return false
//	}
//	defer conn.Close()
//	return true
//}

//func Scan(config cmd.Configs) {
//	// Get the ports from command line flags
//	port1 := flag.Int("p1", 80, "port number 1")
//	port2 := flag.Int("p2", 443, "port number 2")
//	port3 := flag.Int("p3", 8080, "port number 3")
//	flag.Parse()
//
//	ports := []int{*port1, *port2, *port3}
//	ipLists := make([][]string, len(ports))
//	var results []Result
//
//	var wg sync.WaitGroup
//	for ip := 0; ip < 256; ip++ {
//		for i, port := range ports {
//			wg.Add(1)
//			go func(ip, port, i int) {
//				defer wg.Done()
//				ipStr := "39.156.66." + strconv.Itoa(ip)
//				if checkPort("tcp", ipStr, port) {
//					ipLists[i] = append(ipLists[i], "\""+ipStr+"\"")
//					results = append(results, Result{IP: ipStr, Port: strconv.Itoa(port)})
//				}
//			}(ip, port, i)
//		}
//	}
//
//	wg.Wait()
//
//	for i, port := range ports {
//		file, err := os.Create(strconv.Itoa(port) + "_ip.txt")
//		if err != nil {
//			fmt.Println(err)
//			os.Exit(1)
//		}
//		defer file.Close()
//		file.WriteString("[" + strings.Join(ipLists[i], ",") + "]")
//	}
//
//	// Create JSON result file
//	file, _ := json.MarshalIndent(results, "", " ")
//	_ = os.WriteFile("results.json", file, 0644)
//}
