package storage

import (
	"github.com/MayMistery/noscan/cmd"
)

type IpCache struct {
	Ip     string `storm:"id,increment"`
	IpInfo cmd.IpInfo
}
