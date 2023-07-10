package utils

import (
	"fmt"
	"github.com/MayMistery/noscan/cmd"
	"github.com/MayMistery/noscan/storage/bolt"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	var test cmd.IpInfo
	test.DeviceInfo = "vvv"
	//test.Services = "ccc"

	test1 := make(map[string]cmd.IpInfo)
	test1["1.1.1.1"] = test

	//OutputJsonResult()
}

func TestInitResultMap(t *testing.T) {
	cmd.Config.DBFilePath = "../data/database.db"
	bolt.InitAsyncDatabase()
	cmd.Config.OutputFilepath = "../result/result.json"
	InitResultMap()
	OutputResultMap()
}

func TestTimestamp(t *testing.T) {
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	fmt.Println(formattedTime)
}
