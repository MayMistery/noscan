package honeypot

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
)

type Service struct {
	Port       int      `json:"port"`
	Protocol   string   `json:"protocol"`
	ServiceApp []string `json:"service_app"`
}

type IPAddress struct {
	Services   []Service   `json:"services"`
	DeviceInfo interface{} `json:"deviceinfo"`
	Honeypot   []string    `json:"honeypot"`
	Timestamp  string      `json:"timestamp"`
}

func FindAll(inputfile string) {
	data, err := os.ReadFile(inputfile)
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}

	mu := sync.Mutex{}

	var ipAddresses map[string]IPAddress
	err = json.Unmarshal(data, &ipAddresses)
	if err != nil {
		fmt.Println("Error parsing JSON", err)
		return
	}

	var wg sync.WaitGroup
	for ip, ipAddress := range ipAddresses {
		for _, service := range ipAddress.Services {
			if service.Protocol == "ssh" {
				wg.Add(1)
				go func(ip string, service Service, ipAddress IPAddress) {
					defer wg.Done()
					conn, isKippo := isKippoHoneypot(ip, service.Port)
					if isKippo {
						fmt.Println(ip, service.Port, "kippo")
						ipAddress.Honeypot = append(ipAddress.Honeypot, strconv.Itoa(service.Port)+"/kippo")
					}
					if isHfishHoneypot(conn) {
						ipAddress.Honeypot = append(ipAddress.Honeypot, strconv.Itoa(service.Port)+"/Hfish")
					}
					mu.Lock()
					ipAddresses[ip] = ipAddress
					mu.Unlock()
				}(ip, service, ipAddress)
			}
		}
	}
	wg.Wait()

	output, err := json.MarshalIndent(ipAddresses, "", "  ")
	if err != nil {
		fmt.Println("Error generating JSON", err)
		return
	}

	err = os.WriteFile("output.json", output, 0644)
	if err != nil {
		fmt.Println("Error writing to file", err)
		return
	}

	fmt.Println("Successfully updated JSON with kippo honeypots!")
}
