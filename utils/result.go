package utils

import (
	"encoding/json"
	"fmt"
	"github.com/MayMistery/noscan/cmd"
	"github.com/MayMistery/noscan/storage"
	"github.com/MayMistery/noscan/storage/bolt"
	"os"
	"sync"
	"time"
)

var mu sync.RWMutex
var icu sync.RWMutex

var (
	Result      = make(map[string]cmd.IpInfo)
	PortScanned = make(map[string]map[int]bool)
)

func InitResultMap() {
	// Retrieve data from the database
	var ips []storage.IpCache
	err := bolt.DB.Ipdb.All(&ips)
	if err != nil {
		cmd.ErrLog("Failed to retrieve data from the database: %v", err)
		return
	}
	//fmt.Println(ips)
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
		ipInfo.Services = ip.Services

		// Add the IpInfo to the Result map
		mu.Lock()
		Result[ip.Ip] = ipInfo
		mu.Unlock()

		// Add the scanned ports to the PortScaned map
		portScanedMap := make(map[int]bool)
		for _, portInfoStore := range ip.Services {
			portScanedMap[portInfoStore.Port] = true
		}
		icu.Lock()
		PortScanned[ip.Ip] = portScanedMap
		icu.Unlock()
	}
}

func OutputResultMap() {
	mu.RLock()
	ipInfos := Result
	mu.RUnlock()

	filepath := cmd.Config.OutputFilepath
	file, err := os.Create(filepath)
	if err != nil {
		cmd.ErrLog("Failed to create file: %v", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			cmd.ErrLog("Failed to close file: %v", err)
		}
	}(file)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(ipInfos)
	if err != nil {
		cmd.ErrLog("Failed to encode JSON: %v", err)
		return
	}

	fmt.Println("JSON file has been created.")
}

func AddPortInfo(host string, info *cmd.PortInfo) {
	icu.RLock()
	_, ok := PortScanned[host]
	icu.RUnlock()
	if !ok {
		icu.Lock()
		PortScanned[host] = make(map[int]bool)
		icu.Unlock()
	}
	mu.RLock()
	ipInfo, ok1 := Result[host]
	mu.RUnlock()

	if ok1 {
		icu.RLock()
		_, ok2 := PortScanned[host][info.Port]
		icu.RUnlock()

		if ok2 {
			updatePortInfo(host, info)
		} else {
			ipInfo.Services = append(ipInfo.Services, info)
			icu.Lock()
			PortScanned[host][info.Port] = true
			icu.Unlock()
		}
	} else {
		ipInfo = cmd.IpInfo{}
		ipInfo.Services = []*cmd.PortInfo{info}
		icu.Lock()
		PortScanned[host][info.Port] = true
		icu.Unlock()
	}

	//add timestamp
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	ipInfo.Timestamp = formattedTime

	mu.Lock()
	Result[host] = ipInfo
	mu.Unlock()
	mu.RLock()
	device := Result[host].DeviceInfo
	honey := Result[host].Honeypot
	timeStamp := Result[host].Timestamp
	mu.RUnlock()

	//TODO unlock or not
	bolt.UpdateCacheAsync(&storage.IpCache{
		Ip:         host,
		Services:   ipInfo.Services,
		DeviceInfo: device,
		Honeypot:   honey,
		Timestamp:  timeStamp,
	})
}

func updatePortInfo(host string, info *cmd.PortInfo) {
	mu.RLock()
	length := len(Result[host].Services)
	mu.RUnlock()

	for i := 0; i < length; i++ {
		mu.RLock()
		tmp := Result[host].Services[i]
		mu.RUnlock()
		if tmp.Port == info.Port {
			info.ServiceApp = RemoveDuplicateStringArr(append(info.ServiceApp, tmp.ServiceApp...))
			mu.Lock()
			Result[host].Services[i] = info
			mu.Unlock()
			return
		}
	}
}

func UpdateServiceInfo(host string, port int, serviceInfo []string) {
	mu.RLock()
	length := len(Result[host].Services)
	mu.RUnlock()
	for i := 0; i < length; i++ {
		mu.RLock()
		tmp := Result[host].Services[i]
		mu.RUnlock()
		if tmp.Port == port {
			tmp.ServiceApp = RemoveDuplicateStringArr(append(tmp.ServiceApp, serviceInfo...))
			bolt.UpdateServiceInfoAsync(host, port, tmp)
			mu.Lock()
			Result[host].Services[i] = tmp
			mu.Unlock()
			return
		}
	}
}

func UpdateDeviceInfo(host string, deviceInfo string) {
	mu.RLock()
	ipInfo, ok1 := Result[host]
	mu.RUnlock()
	if ok1 {
		ipInfo.DeviceInfo = deviceInfo
	} else {
		ipInfo = cmd.IpInfo{}
		ipInfo.DeviceInfo = deviceInfo
	}
	mu.Lock()
	Result[host] = ipInfo
	mu.Unlock()

	bolt.UpdateDeviceInfoAsync(host, deviceInfo)
}

func UpdateHoneypot(host string, honeypot []string) {
	mu.RLock()
	ipInfo, ok1 := Result[host]
	mu.RUnlock()
	if ok1 {
		ipInfo.Honeypot = RemoveDuplicateStringArr(append(ipInfo.Honeypot, honeypot...))
	} else {
		ipInfo = cmd.IpInfo{}
		ipInfo.Honeypot = honeypot
	}
	mu.Lock()
	Result[host] = ipInfo
	mu.Unlock()

	bolt.UpdateHoneypotAsync(host, ipInfo.Honeypot)
}
