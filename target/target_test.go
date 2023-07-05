package target

import "testing"

func TestReadIPAddressesFromFile(t *testing.T) {
	ipAddresses, err := ReadIPAddressesFromFile("target")
	if err != nil {
		t.Errorf("无法读取文件：%v", err)
		return
	}

	for _, ip := range ipAddresses {
		t.Log(ip)
	}
}
