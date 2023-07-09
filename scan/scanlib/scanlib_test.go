package scanlib

import (
	"fmt"
	"testing"
	"time"
)

func TestScan(t *testing.T) {
	var scanner = New()
	host := "45.126.125.13"
	port := 80
	status, response := scanner.ScanTimeout(host, port, time.Second*30)
	fmt.Println(status, response.FingerPrint.Service, host, ":", port)
}
