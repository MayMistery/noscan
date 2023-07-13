package scan

import (
	"fmt"
	"github.com/MayMistery/noscan/lib/simplehttp"
	"net/url"
	"testing"
)

func TestWappalyzer(t *testing.T) {
	testurl, _ := url.Parse("https://remix.ethereum.org")
	resp, _ := httpRequest(testurl, nil, simplehttp.NewClient())
	services := getServiceAppFromWappalyzer(resp.Response)
	for _, service := range services {
		fmt.Println(service)
	}
}
