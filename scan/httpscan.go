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
	"regexp"
	"strconv"
	"strings"
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

type deviceMapping struct {
	device       string
	fingerprints []string
}

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
			//遍历ProductName []string

			var device, honeyPot []string
			port, _ := strconv.Atoi(URL.Port())

			deviceMappings := []deviceMapping{
				{device: "webcam", fingerprints: []string{"摄像头", "webcam", "camera"}},
				{device: "router", fingerprints: []string{"路由器", "router"}},
				{device: "gateway", fingerprints: []string{"网关", "防火墙", "gateway"}},
				{device: "vpn", fingerprints: []string{"虚拟专用网络", "vpn"}},
				{device: "storage", fingerprints: []string{"存储设备", "storage"}},
				{device: "switch", fingerprints: []string{"交换机", "switch"}},
				{device: "printers", fingerprints: []string{"打印机设备", "printers"}},
				{device: "proxy server", fingerprints: []string{"代理服务器", "proxy"}},
				{device: "kvm", fingerprints: []string{"虚拟化平台", "kvm"}},
				{device: "cdn", fingerprints: []string{"内容分发平台", "CloudFlare", "cdn"}},
				{device: "phone", fingerprints: []string{"移动通信", "phone"}},
				{device: "bridge", fingerprints: []string{"虚拟网络设备", "bridge"}},
				{device: "security", fingerprints: []string{"安全防护设备", "security"}},
				{device: "honeypot", fingerprints: []string{"蜜罐", "honeypot"}},
			}

			for i := 0; i < len(finger.ProductName); i++ {
				//strip /t
				finger.ProductName[i] = strings.TrimRight(finger.ProductName[i], "\t")

				// 匹配 device 类型
				matchedDevice := "other"
				for _, mapping := range deviceMappings {
					for _, fingerprint := range mapping.fingerprints {
						matched, err := regexp.MatchString("(?i)"+fingerprint, finger.ProductName[i])
						if err != nil {
							continue
						}
						if matched {
							matchedDevice = mapping.device
							break
						}
					}
					if matchedDevice != "other" {
						break
					}
				}

				if matchedDevice == "honeypot" {
					honeyInfo := string(port) + "/N"
					//TODO strip honeypot name
					honeyPot = append(honeyPot, honeyInfo)
				} else {
					//TODO strip device name
					deviceInfo := matchedDevice + "/" + finger.ProductName[i]
					device = append(device, deviceInfo)
				}
			}

			//identify device

			//identify honeypot

			utils.UpdateServiceInfo(URL.Hostname(), port, finger.ProductName)
			utils.UpdateServiceInfo(URL.Hostname(), port, device)

			//TODO updateIpInfo
			//utils.UpdateIpInfo(URL.Hostname(), honeyPot)
			//utils.UpdateDevice
			//utils.UpdateHoneypot
		}
	}

	return httpScanPool
}

func HttpHandlerError(url *url.URL, err error) {
	fmt.Println("URLScanner Error: ", url.String(), err)
}

//webcam : "摄像头" "webcam" "camera"
//router : "路由器" "router"
//gateway : "网关" "防火墙" "gateway"
//vpn : "虚拟专用网络" "vpn"
//storage : "存储设备" "storage"
//switch : "交换机" "switch"
//printers : "打印机设备" "printers"
//proxy server : "代理服务器" "proxy"
//kvm : "虚拟化平台" "kvm"
//cdn : "内容分发平台" "CloudFlare" "cdn"
//phone : "移动通信" "phone"
//bridge : "虚拟网络设备" "bridge"
//security : "安全防护设备" "security"
//honeypot : "蜜罐" "honeypot"
//other :
