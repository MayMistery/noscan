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
			HandlerNotMatched(value)
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
	utils.AddPortInfo(value.IP.String(), portInfo)
}

func HandlerNotMatched(value Address) {
	portInfo := cmd.PortInfo{
		Port:     value.Port,
		Protocol: "unknow",
	}
	utils.AddPortInfo(value.IP.String(), portInfo)
}

func HandlerMatched(value Address, response *gonmap.Response) {
	protocol := response.FingerPrint.Service
	portInfo := cmd.PortInfo{
		Port:     value.Port,
		Protocol: protocol,
	}
	utils.AddPortInfo(value.IP.String(), portInfo)
	//TODO Further application probing
}
