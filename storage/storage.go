package storage

import "github.com/MayMistery/noscan/cmd"

//type PortInfoStore struct {
//	*cmd.PortInfo
//	Banner *lib.Response
//}

type IpCache struct {
	Ip         string `storm:"id,increment"`
	DeviceInfo string
	Honeypot   []string
	Services   []*cmd.PortInfo
	Timestamp  string
}

type BannerCache struct {
	Ip     string `storm:"id,increment"`
	Port   int
	Banner string
}
