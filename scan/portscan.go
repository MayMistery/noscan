package scan

import (
	"fmt"
	"github.com/MayMistery/noscan/cmd"
	"github.com/MayMistery/noscan/scan/scanlib"
	"github.com/MayMistery/noscan/utils"
	"github.com/lcvvvv/appfinger"
	"github.com/lcvvvv/simplehttp"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func PortScanPool() *cmd.Pool {
	portScanPool := cmd.NewPool(cmd.Config.Threads)
	portScanPool.Function = func(in interface{}) {
		nmap := scanlib.New()
		nmap.SetTimeout(cmd.Config.Timeout)
		if cmd.Config.DeepInspection == true {
			nmap.OpenDeepIdentify()
		}
		value := in.(Address)
		status, response := nmap.ScanTimeout(value.IP.String(), value.Port, 100*cmd.Config.Timeout)
		switch status {
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

func PortHandlerOpen(value Address) {
	protocol := scanlib.GuessProtocol(value.Port)
	portInfo := &cmd.PortInfo{
		Port:       value.Port,
		Protocol:   protocol,
		ServiceApp: nil,
	}
	utils.AddPortInfo(value.IP.String(), portInfo, nil)
}

func PortHandlerNotMatched(value Address, response *scanlib.Response) {
	portInfo := &cmd.PortInfo{
		Port:     value.Port,
		Protocol: "unknow",
	}
	utils.AddPortInfo(value.IP.String(), portInfo, response)
}

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

	utils.AddPortInfo(value.IP.String(), portInfo, response)
	URLRaw := fmt.Sprintf("%s://%s:%d", protocol, value.IP.String(), value.Port)
	URL, _ := url.Parse(URLRaw)
	if appfinger.SupportCheck(URL.Scheme) == true {
		pushURLTarget(URL, response)
		return
	}
}

func pushURLTarget(URL *url.URL, response *scanlib.Response) {
	var cli *http.Client
	//判断是否初始化client
	if cmd.Config.Proxy != "" || cmd.Config.Timeout != 3*time.Second {
		cli = simplehttp.NewClient()
	}
	//判断是否需要设置代理
	if cmd.Config.Proxy != "" {
		simplehttp.SetProxy(cli, cmd.Config.Proxy)
	}
	//判断是否需要设置超时参数
	if cmd.Config.Timeout != 3*time.Second {
		simplehttp.SetTimeout(cli, cmd.Config.Timeout)
	}

	HttpScanner.Push(HttpTarget{URL, response, nil, cli})
}

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
