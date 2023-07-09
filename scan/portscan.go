package scan

import (
	"github.com/MayMistery/noscan/cmd"
	"github.com/MayMistery/noscan/scan/gonmap"
	"github.com/MayMistery/noscan/utils"
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
			HandlerNotMatched(value, response)
		case gonmap.Matched:
			HandlerMatched(value, response)
		}
	}

	return portScanPool
}

func HandlerOpen(value Address) {
	protocol := gonmap.GuessProtocol(value.Port)
	portInfo := &cmd.PortInfo{
		Port:       value.Port,
		Protocol:   protocol,
		ServiceApp: nil,
	}
	utils.AddPortInfo(value.IP.String(), portInfo, nil)
}

func HandlerNotMatched(value Address, response *gonmap.Response) {
	portInfo := &cmd.PortInfo{
		Port:     value.Port,
		Protocol: "unknow",
	}
	utils.AddPortInfo(value.IP.String(), portInfo, response)
}

func HandlerMatched(value Address, response *gonmap.Response) {
	protocol := response.FingerPrint.Service
	portInfo := &cmd.PortInfo{
		Port:     value.Port,
		Protocol: protocol,
	}
	//TODO Further application probing
	if protocol == "http" || protocol == "https" {

	}

	utils.AddPortInfo(value.IP.String(), portInfo, response)
}
