package utils

import (
	"encoding/json"
	"fmt"
	"github.com/MayMistery/noscan/cmd"
	"os"
)

func OutputJsonResult() {
	ipInfos := []cmd.IpInfo{
		{
			Host: "192.108.10.11",
			Services: []cmd.PortInfo{
				{
					Port:       22,
					Protocol:   "ssh",
					ServiceApp: []string{"SSH/N"},
				},
				{
					Port:       23,
					Protocol:   "telnet",
					ServiceApp: []string{"Telnet/N"},
				},
				{
					Port:       80,
					Protocol:   "http",
					ServiceApp: []string{"Nginx/0.8.53", "WordPress/4.7"},
				},
			},
			DeviceInfo: "route/fritz",
			Honeypot:   []string{"22/kippo", "80/glastopf"},
			Timestamp:  "2023-05-06 20:21:22",
		},
		{
			Host: "192.118.30.33",
			Services: []cmd.PortInfo{
				{
					Port:       22,
					Protocol:   "ssh",
					ServiceApp: []string{"SSH/N"},
				},
				{
					Port:       8888,
					Protocol:   "http",
					ServiceApp: []string{"Apache/2.2.15", "Joomla/N"},
				},
				{
					Port:       3306,
					Protocol:   "mysql",
					ServiceApp: []string{"mysql/N"},
				},
				{
					Port:       3389,
					Protocol:   "rdp",
					ServiceApp: []string{"rdp/N"},
				},
				{
					Port:       9999,
					Protocol:   "",
					ServiceApp: nil,
				},
			},
			DeviceInfo: "",
			Honeypot:   []string{"22/kippo", "80/glastopf", "3306/N"},
			Timestamp:  "2023-06-07 15:16:17",
		},
	}

	cidrInfo := cmd.CIDRInfo{ipInfos}

	file, err := os.Create("res.json")
	if err != nil {
		fmt.Println("Failed to create file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(cidrInfo)
	if err != nil {
		fmt.Println("Failed to encode JSON:", err)
		return
	}

	fmt.Println("JSON file has been created.")
}
