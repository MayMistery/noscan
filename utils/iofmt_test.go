package utils

import (
	"github.com/MayMistery/noscan/cmd"
	"testing"
)

func TestName(t *testing.T) {
	var test cmd.IpInfo
	test.DeviceInfo = "vvv"
	//test.Services = "ccc"

	test1 := make(map[string]cmd.IpInfo)
	test1["1.1.1.1"] = test

	OutputJsonResult(test1)
}
