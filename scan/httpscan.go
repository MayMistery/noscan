package scan

import (
	"errors"
	"fmt"
	"github.com/MayMistery/noscan/cmd"
	"github.com/MayMistery/noscan/scan/scanlib"
	"github.com/MayMistery/noscan/utils"
	"github.com/lcvvvv/appfinger"
	"net/http"
	"net/url"
	"strconv"
)

type HttpTarget struct {
	URL      *url.URL
	response *scanlib.Response
	req      *http.Request
	client   *http.Client
}

const (
	NotSupportProtocol = "protocol is not support"
)

func HttpScanPool() *utils.Pool {
	httpScanPool := utils.NewPool(cmd.Config.Threads)
	httpScanPool.Function = func(in interface{}) {
		value := in.(HttpTarget)
		URL := value.URL
		response := value.response
		req := value.req
		cli := value.client
		if appfinger.SupportCheck(URL.Scheme) == false {
			fmt.Println(URL, errors.New(NotSupportProtocol))
			return
		}
		var banner *appfinger.Banner
		var finger *appfinger.FingerPrint
		var err error
		if response == nil || req != nil || cli != nil {
			banner, err = appfinger.GetBannerWithURL(URL, req, cli)
			if err != nil {
				HttpHandlerError(URL, err)
				return
			}
			finger = appfinger.Search(URL, banner)
		} else {
			banner, err = appfinger.GetBannerWithURL(URL, req, cli)
			if err != nil {
				HttpHandlerError(URL, err)
				return
			}
			finger = appfinger.Search(URL, banner)
		}
		if len(finger.ProductName) > 0 {
			port, _ := strconv.Atoi(URL.Port())
			utils.UpdateServiceInfo(URL.Hostname(), port, finger.ProductName)
		}
	}

	return httpScanPool
}

func HttpHandlerError(url *url.URL, err error) {
	fmt.Println("URLScanner Error: ", url.String(), err)
}
