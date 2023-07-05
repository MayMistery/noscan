package scan

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

/*
Copyright 2023 noname

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under tgo he License.
*/

package main

import (
"encoding/json"
"flag"
"fmt"
"net"
"os"
"strconv"
"strings"
"sync"
"time"
)
type Result struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
}

func checkPort(protocol, hostname string, port int) bool {
	address := net.JoinHostPort(hostname, strconv.Itoa(port))
	//fmt.Println("hello", port, hostname)
	conn, err := net.DialTimeout(protocol, address, 2*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func main() {
	// Get the ports from command line flags
	port1 := flag.Int("p1", 80, "port number 1")
	port2 := flag.Int("p2", 443, "port number 2")
	port3 := flag.Int("p3", 8080, "port number 3")
	flag.Parse()

	ports := []int{*port1, *port2, *port3}
	ipLists := make([][]string, len(ports))
	var results []Result

	var wg sync.WaitGroup
	for ip := 0; ip < 256; ip++ {
		for i, port := range ports {
			wg.Add(1)
			go func(ip, port, i int) {
				defer wg.Done()
				ipStr := "39.156.66." + strconv.Itoa(ip)
				if checkPort("tcp", ipStr, port) {
					ipLists[i] = append(ipLists[i], "\""+ipStr+"\"")
					results = append(results, Result{IP: ipStr, Port: strconv.Itoa(port)})
				}
			}(ip, port, i)
		}
	}

	wg.Wait()

	for i, port := range ports {
		file, err := os.Create(strconv.Itoa(port) + "_ip.txt")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer file.Close()
		file.WriteString("[" + strings.Join(ipLists[i], ",") + "]")
	}

	// Create JSON result file
	file, _ := json.MarshalIndent(results, "", " ")
	_ = os.WriteFile("results.json", file, 0644)
}
