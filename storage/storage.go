package storage

import (
	"github.com/MayMistery/noscan/cmd"
	"github.com/MayMistery/noscan/storage/bolt"
)

type IpCache struct {
	Ip     string `storm:"id,increment"`
	IpInfo cmd.IpInfo
}

var DB *bolt.Storage
