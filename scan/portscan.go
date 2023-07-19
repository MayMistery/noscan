package scan

import (
	"fmt"
	"github.com/MayMistery/noscan/cmd"
	"github.com/MayMistery/noscan/lib/appfinger"
	"github.com/MayMistery/noscan/lib/scanlib"
	"github.com/MayMistery/noscan/lib/simplehttp"
	"github.com/MayMistery/noscan/storage"
	"github.com/MayMistery/noscan/storage/bolt"
	"github.com/MayMistery/noscan/utils"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

var (
	secondChance          = make(map[string]map[int]struct{})
	secondChanceMapMutex  = sync.RWMutex{}
	secondScanAddressChan = make(chan Address)
)

func init() {
	go func() {
		for addr := range secondScanAddressChan {
			PortScanner.Push(addr)
		}
	}()
}

// PortScanPool creates a new pool for port scanning.
// It sets the function of the pool to be a function that scans a single port.
func PortScanPool() *cmd.Pool {
	portScanPool := cmd.NewPool(cmd.Config.Threads)
	portScanPool.Function = func(in interface{}) {
		nmap := scanlib.New()
		nmap.SetTimeout(cmd.Config.Timeout)
		if cmd.Config.DeepInspection == true {
			nmap.OpenDeepIdentify()
		}
		value := in.(Address)
		//fmt.Println("[+]Scanning", value.IP.String(), value.Port)
		status, response := nmap.ScanTimeout(value.IP.String(), value.Port, cmd.Config.Timeout)
		switch status {
		case scanlib.Closed:
			// If closed scan twice
			secondChanceMapMutex.RLock()
			_, isSet := secondChance[value.IP.String()]
			secondChanceMapMutex.RUnlock()
			if !isSet {
				secondChanceMapMutex.Lock()
				secondChance[value.IP.String()] = make(map[int]struct{})
				secondChanceMapMutex.Unlock()
			}
			secondChanceMapMutex.RLock()
			_, isScan := secondChance[value.IP.String()][value.Port]
			secondChanceMapMutex.RUnlock()
			if !isScan {
				//fmt.Println("[*]Second scan", value.IP.String(), value.Port)
				secondScanAddressChan <- value
				secondChanceMapMutex.Lock()
				secondChance[value.IP.String()][value.Port] = struct{}{}
				secondChanceMapMutex.Unlock()
			}
		case scanlib.Open:
			PortHandlerOpen(value)
		case scanlib.NotMatched:
			PortHandlerNotMatched(value, response)
		case scanlib.Matched:
			PortHandlerMatched(value, response)
		}
	}

	return portScanPool
}

// PortHandlerOpen handles the case when a port is open during port scanning.
// It guesses the protocol based on the port number and updates the port information.
func PortHandlerOpen(value Address) {
	protocol := scanlib.GuessProtocol(value.Port)
	portInfo := &cmd.PortInfo{
		Port:       value.Port,
		Protocol:   protocol,
		ServiceApp: nil,
	}
	utils.AddPortInfo(value.IP.String(), portInfo)
}

// PortHandlerNotMatched handles the case when a port does not match any fingerprints during port scanning.
// It sets the protocol to "unknown" and updates the port information.
func PortHandlerNotMatched(value Address, response *scanlib.Response) {
	portInfo := &cmd.PortInfo{
		Port:     value.Port,
		Protocol: "unknown",
	}
	utils.AddPortInfo(value.IP.String(), portInfo)
}

// PortHandlerMatched handles the case when a port matches a fingerprint during port scanning.
// It gets the protocol from the fingerprint, gets the product version, and updates the port information and banner cache.
// It also pushes the URL to the HTTP scanner if the URL scheme is supported by the app finger.
func PortHandlerMatched(value Address, response *scanlib.Response) {
	protocol := response.FingerPrint.Service
	var services []string
	if product := getProductVersionFromNmap(response); product != "" {
		services = []string{product}
	}
	portInfo := &cmd.PortInfo{
		Port:       value.Port,
		Protocol:   protocol,
		ServiceApp: services,
	}
	fmt.Println("[+]Get service:", value.IP.String(), value.Port, response.FingerPrint.Service, ":", response.FingerPrint.ProductName)
	utils.AddPortInfo(value.IP.String(), portInfo)
	bolt.UpdateBannerCacheAsync(&storage.BannerCache{Ip: value.IP.String(), Port: value.Port, Banner: response.Raw})
	URLRaw := fmt.Sprintf("%s://%s:%d", protocol, value.IP.String(), value.Port)
	URL, _ := url.Parse(URLRaw)
	if appfinger.SupportCheck(URL.Scheme) == true {
		pushURLTarget(URL, response)
		return
	}
}

// pushURLTarget pushes the given URL to the HTTP scanner.
// It creates a new client if necessary, sets the proxy and timeout if specified, and pushes the HTTP target to the HTTP scanner.
func pushURLTarget(URL *url.URL, response *scanlib.Response) {
	var cli *http.Client
	//判断是否初始化client
	if cmd.Config.Proxy != "" || cmd.Config.Timeout != 3*time.Second {
		cli = simplehttp.NewClient()
	}
	//判断是否需要设置代理
	if cmd.Config.Proxy != "" {
		err := simplehttp.SetProxy(cli, cmd.Config.Proxy)
		if err != nil {
			cmd.ErrLog("SetProxy error %v", err)
			return
		}
	}
	//判断是否需要设置超时参数
	if cmd.Config.Timeout != 3*time.Second {
		simplehttp.SetTimeout(cli, cmd.Config.Timeout)
	}

	HttpScanner.Push(HttpTarget{URL, nil, cli})
}

// getProductVersionFromNmap gets the product version from a scanlib.Response.
// It returns a string that combines the product name and version, separated by a "/".
func getProductVersionFromNmap(response *scanlib.Response) string {
	var (
		version string
		product string
	)
	if response.FingerPrint.ProductName == "" {
		return ""
	} else {
		product = response.FingerPrint.ProductName
	}

	if response.FingerPrint.Version != "" {
		version = response.FingerPrint.Version
	} else {
		version = "N"
	}

	return strings.Join([]string{product, version}, "/")
}
