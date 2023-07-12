package scan

import (
	"errors"
	"fmt"
	"github.com/MayMistery/noscan/cmd"
	appfinger2 "github.com/MayMistery/noscan/lib/appfinger"
	"github.com/MayMistery/noscan/lib/scanlib"
	"github.com/MayMistery/noscan/utils"
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

func HttpScanPool() *cmd.Pool {
	httpScanPool := cmd.NewPool(cmd.Config.Threads)
	httpScanPool.Function = func(in interface{}) {
		value := in.(HttpTarget)
		URL := value.URL
		response := value.response
		req := value.req
		cli := value.client
		if appfinger2.SupportCheck(URL.Scheme) == false {
			cmd.ErrLog("%s %v", URL.Host, errors.New(NotSupportProtocol))
			fmt.Println(URL, errors.New(NotSupportProtocol))
			return
		}
		var banner *appfinger2.Banner
		var finger *appfinger2.FingerPrint
		var err error
		if response == nil || req != nil || cli != nil {
			banner, err = appfinger2.GetBannerWithURL(URL, req, cli)
			if err != nil {
				HttpHandlerError(URL, err)
				return
			}
			finger = appfinger2.Search(URL, banner)
		} else {
			banner, err = appfinger2.GetBannerWithURL(URL, req, cli)
			if err != nil {
				HttpHandlerError(URL, err)
				return
			}
			finger = appfinger2.Search(URL, banner)
		}
		if len(finger.ProductName) > 0 {
			HandleAppFingerprint(URL, finger.ProductName)
		}
	}

	return httpScanPool
}

func HttpHandlerError(url *url.URL, err error) {
	cmd.ErrLog("URLScanner Error: %s %v", url.String(), err)
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

type deviceMapping struct {
	device       string
	fingerprints []string
}

var deviceMappings = []deviceMapping{
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

func HandleAppFingerprint(url *url.URL, inputFinger []string) {
	//遍历ProductName []string
	var honeyPot, service []string
	var deviceInfo string
	for i := 0; i < len(inputFinger); i++ {
		//strip /t
		inputFinger[i] = strings.TrimRight(inputFinger[i], "\t")

		// 匹配 device 类型
		matchedDevice := "other"
		for _, mapping := range deviceMappings {
			for _, fingerprint := range mapping.fingerprints {
				matched, err := regexp.MatchString("(?i)"+fingerprint, inputFinger[i])
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
			honeyInfo := url.Port() + "/N"
			//TODO strip honeypot name
			honeyPot = append(honeyPot, honeyInfo)
		} else if matchedDevice == "other" {
			service = append(service, inputFinger[i])
		} else {
			//TODO strip device name
			deviceInfo = matchedDevice + "/" + inputFinger[i]
			//device = append(device, deviceInfo)
		}
	}

	port, _ := strconv.Atoi(url.Port())
	utils.UpdateServiceInfo(url.Hostname(), port, service)
	utils.UpdateDeviceInfo(url.Hostname(), deviceInfo)
	utils.UpdateHoneypot(url.Hostname(), honeyPot)
}
