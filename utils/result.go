package utils

import (
	"encoding/json"
	"fmt"
	"github.com/MayMistery/noscan/cmd"
	"github.com/MayMistery/noscan/storage"
	"github.com/MayMistery/noscan/storage/bolt"
	"os"
)

var (
	Result     map[string]cmd.IpInfo
	PortScaned map[string]map[int]bool
)

func OutputJsonResult() {
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

func AddPortInfo(host string, info cmd.PortInfo) {
	if ipInfo, ok1 := Result[host]; ok1 {
		if _, ok2 := PortScaned[host][info.Port]; ok2 {
			UpdatePortInfo(host, info)
		} else {
			ipInfo.Services = append(ipInfo.Services, info)
			PortScaned[host][info.Port] = true
		}
	} else {
		ipInfo = cmd.IpInfo{}
		ipInfo.Services = []cmd.PortInfo{info}
		PortScaned[host][info.Port] = true
	}
	bolt.DB.UpdateCache(storage.IpCache{
		Ip:     host,
		IpInfo: Result[host],
	})
}

func UpdatePortInfo(host string, info cmd.PortInfo) {
	for i := 0; i < len(Result[host].Services); i++ {
		if Result[host].Services[i].Port == info.Port {
			Result[host].Services[i] = info
			return
		}
	}
}