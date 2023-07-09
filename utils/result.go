package utils

import (
	"encoding/json"
	"fmt"
	"github.com/MayMistery/noscan/cmd"
	"github.com/MayMistery/noscan/scan/gonmap"
	"github.com/MayMistery/noscan/storage"
	"github.com/MayMistery/noscan/storage/bolt"
	"log"
	"os"
	"time"
)

var (
	Result      = make(map[string]cmd.IpInfo)
	PortScanned = make(map[string]map[int]bool)
)

func InitResultMap() {
	// Retrieve data from the database
	var ips []storage.IpCache
	err := bolt.DB.Ipdb.All(&ips)
	if err != nil {
		log.Printf("Failed to retrieve data from the database: %v\n", err)
		return
	}
	fmt.Println(ips)
	// Populate the Result and PortScaned maps
	for _, ip := range ips {
		// Create a new IpInfo instance
		ipInfo := cmd.IpInfo{
			Services:   make([]*cmd.PortInfo, len(ip.Services)),
			DeviceInfo: ip.DeviceInfo,
			Honeypot:   ip.Honeypot,
			Timestamp:  ip.Timestamp,
		}

		// Populate the Services field of IpInfo
		for i, portInfoStore := range ip.Services {
			ipInfo.Services[i] = portInfoStore.PortInfo
		}

		// Add the IpInfo to the Result map
		Result[ip.Ip] = ipInfo

		// Add the scanned ports to the PortScaned map
		portScanedMap := make(map[int]bool)
		for _, portInfoStore := range ip.Services {
			portScanedMap[portInfoStore.Port] = true
		}
		PortScanned[ip.Ip] = portScanedMap
	}
}

func OutputResultMap() {
	ipInfos := Result
	filepath := cmd.Config.OutputFilepath
	file, err := os.Create(filepath)
	if err != nil {
		fmt.Println("Failed to create file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(ipInfos)
	if err != nil {
		fmt.Println("Failed to encode JSON:", err)
		return
	}

	fmt.Println("JSON file has been created.")
}

func AddPortInfo(host string, info *cmd.PortInfo, banner *gonmap.Response) {
	if _, ok := PortScanned[host]; !ok {
		PortScanned[host] = make(map[int]bool)
	}

	ipInfo, ok1 := Result[host]
	if ok1 {
		if _, ok2 := PortScanned[host][info.Port]; ok2 {
			updatePortInfo(host, info)
		} else {
			ipInfo.Services = append(ipInfo.Services, info)
			PortScanned[host][info.Port] = true
		}
	} else {
		ipInfo = cmd.IpInfo{}
		ipInfo.Services = []*cmd.PortInfo{info}
		PortScanned[host][info.Port] = true
	}

	//add timestamp
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	ipInfo.Timestamp = formattedTime

	Result[host] = ipInfo

	bolt.DB.UpdateCache(storage.IpCache{
		Ip: host,
		Services: []*storage.PortInfoStore{{
			PortInfo: info,
			Banner:   banner,
		}},
		DeviceInfo: Result[host].DeviceInfo,
		Honeypot:   Result[host].Honeypot,
		Timestamp:  Result[host].Timestamp,
	})
}

func updatePortInfo(host string, info *cmd.PortInfo) {
	for i := 0; i < len(Result[host].Services); i++ {
		if Result[host].Services[i].Port == info.Port {
			Result[host].Services[i] = info
			return
		}
	}
}
