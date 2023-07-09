package storage

import (
	"github.com/MayMistery/noscan/cmd"
	"github.com/MayMistery/noscan/scan/gonmap"
)

type PortInfoStore struct {
	*cmd.PortInfo
	Banner *gonmap.Response
}

type IpCache struct {
	Ip         string   `storm:"id,increment"`
	DeviceInfo string   `json:"deviceinfo"`
	Honeypot   []string `json:"honeypot"`
	Services   []*PortInfoStore
}
