package scanlib

import (
	"fmt"
	"testing"
	"time"
)

func TestScan(t *testing.T) {
	var scanner = New()
	host := "137.184.166.61"
	port := 5672
	scanner.OpenDeepIdentify()
	status, response := scanner.ScanTimeout(host, port, time.Second*30)
	if response != nil {
		fmt.Println(status, response.FingerPrint.Service, host, ":", port, response.FingerPrint.ProductName)
	}
	//time.Sleep(100000)
}
