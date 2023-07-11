package storage

import "github.com/MayMistery/noscan/cmd"

//type PortInfoStore struct {
//	*cmd.PortInfo
//	Banner *scanlib.Response
//}

type IpCache struct {
	Ip         string `storm:"id,increment"`
	DeviceInfo string
	Honeypot   []string
	Services   []*cmd.PortInfo
	Timestamp  string
}
