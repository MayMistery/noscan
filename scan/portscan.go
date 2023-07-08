package scan

import (
	"github.com/MayMistery/noscan/cmd"
	"github.com/MayMistery/noscan/scan/gonmap"
	"github.com/MayMistery/noscan/storage"
	"github.com/MayMistery/noscan/utils"
	"time"
)

func PortScanPool() *utils.Pool {
	portScanPool := utils.NewPool(cmd.Config.Threads)
	portScanPool.Function = func(in interface{}) {
		nmap := gonmap.New()
		nmap.SetTimeout(cmd.Config.Timeout)
		if cmd.Config.DeepInspection == true {
			nmap.OpenDeepIdentify()
		}
		value := in.(Address)
		status, response := nmap.ScanTimeout(value.IP.String(), value.Port, 100*cmd.Config.Timeout)
		switch status {
		case gonmap.Open:
			HandlerOpen(value)
		case gonmap.NotMatched:
			HandlerNotMatched(value, response.Raw)
		case gonmap.Matched:
			HandlerMatched(value, response)
		}
	}

	return portScanPool
}

func HandlerOpen(value Address) {
	protocol := gonmap.GuessProtocol(value.Port)
	portInfo := cmd.PortInfo{
		Port:       value.Port,
		Protocol:   protocol,
		ServiceApp: nil,
	}
	ipCache := storage.IpCache{
		Ip: value.IP.String(),
		IpInfo: cmd.IpInfo{
			Services:   []cmd.PortInfo{portInfo},
			DeviceInfo: "",
			Honeypot:   nil,
			Timestamp:  time.Now().String(),
		},
	}
	cmd.DB.SaveIpCache(ipCache)
}

func HandlerNotMatched(value Address, response string) {

}

func HandlerMatched(value Address, response *gonmap.Response) {
	protocol := response.FingerPrint.Service
	ipInfo := cmd.PortInfo{
		Port:     value.Port,
		Protocol: protocol,
	}
	//TODO Further application probing
}
