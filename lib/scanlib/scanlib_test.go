package scanlib

import (
	"fmt"
	"testing"
	"time"
)

func TestScan(t *testing.T) {
	var scanner = New()
	host := "204.168.173.224"
	port := 443
	scanner.OpenDeepIdentify()
	status, response := scanner.ScanTimeout(host, port, time.Second*30)
	if response != nil {
		fmt.Println(status, response.FingerPrint.Service, host, ":", port, response.FingerPrint.ProductName)
	}
	//time.Sleep(100000)
}
