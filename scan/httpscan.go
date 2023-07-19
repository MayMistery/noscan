package scan

import (
	"errors"
	"fmt"
	"github.com/MayMistery/noscan/cmd"
	"github.com/MayMistery/noscan/lib/appfinger"
	"github.com/MayMistery/noscan/lib/simplehttp"
	"github.com/MayMistery/noscan/utils"
	wappalyzer "github.com/projectdiscovery/wappalyzergo"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type HttpTarget struct {
	URL    *url.URL
	req    *http.Request
	client *http.Client
}

const (
	NotSupportProtocol = "protocol is not support"
)

// HttpScanPool creates a new pool for HTTP scanning.
// It sets the function of the pool to be a function that scans a single HTTP target.
func HttpScanPool() *cmd.Pool {
	httpScanPool := cmd.NewPool(cmd.Config.Threads/2 + 1)
	httpScanPool.Function = func(in interface{}) {
		value := in.(HttpTarget)
		URL := value.URL
		req := value.req
		cli := value.client

		resp, err := httpRequest(URL, req, cli)
		var serviceApp []string
		if err == nil && resp != nil {
			serviceApp1 := getServiceAppFromAppFinger(URL, resp)
			serviceApp2 := getServiceAppFromWappalyzer(resp.Response)
			serviceApp = append(serviceApp1, serviceApp2...)
		} else {
			serviceApp = []string{}
		}
		if err != nil {
			HttpHandlerError(URL, err)
			return
		}

		if len(serviceApp) > 0 {
			HandleAppFingerprint(URL, serviceApp)
		}
	}

	return httpScanPool
}

// HttpHandlerError logs and prints an error that occurred during HTTP scanning.
func HttpHandlerError(url *url.URL, err error) {
	cmd.ErrLog("URLScanner Error: %s %v", url.String(), err)
	//TODO
	//fmt.Println("URLScanner Error: ", url.String(), err)
}

// httpRequest sends an HTTP request to the given URL and returns the response.
// If the request or client is not provided, it creates new ones.
func httpRequest(URL *url.URL, req *http.Request, cli *http.Client) (*simplehttp.Response, error) {
	if req == nil {
		req, _ = simplehttp.NewRequest(http.MethodGet, URL.String(), nil)
	}

	req.Header.Set("User-Agent", simplehttp.RandomUserAgent())

	if cli == nil {
		cli = simplehttp.NewClient()
		simplehttp.SetTimeout(cli, cmd.Config.Timeout)
	}

	resp, err := simplehttp.Do(cli, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// getServiceAppFromAppFinger gets the service application from the app finger and the given response.
func getServiceAppFromAppFinger(url *url.URL, resp *simplehttp.Response) []string {
	if appfinger.SupportCheck(url.Scheme) == false || resp == nil {
		cmd.ErrLog("%s %v", url.Host, errors.New(NotSupportProtocol))
		fmt.Println(url, errors.New(NotSupportProtocol))
		return []string{}
	}
	var banner *appfinger.Banner
	var finger *appfinger.FingerPrint
	banner, _ = appfinger.GetBannerWithResponse(url, resp.Raw.String())
	finger = appfinger.Search(url, banner)
	return finger.ProductName
}

// getServiceAppFromWappalyzer gets the service application from the wappalyzer and the given response.
func getServiceAppFromWappalyzer(resp *http.Response) []string {
	if resp == nil {
		return []string{}
	}

	data, _ := io.ReadAll(resp.Body) // Ignoring error for example

	var serviceApp []string
	wappalyzerClient, _ := wappalyzer.New()
	fingerprints := wappalyzerClient.Fingerprint(resp.Header, data)
	for key := range fingerprints {
		if strings.Contains(key, ":") {
			serviceApp = append(serviceApp, strings.Replace(key, ":", "/", 1))
		} else {
			serviceApp = append(serviceApp, key+"/N")
		}
	}
	return serviceApp
}

type deviceMapping struct {
	device       string
	fingerprints []string
}

// deviceMapping is a struct to hold information about a device mapping.
// It includes a device name and a slice of fingerprints.
var deviceMappings = []deviceMapping{
	{device: "webcam", fingerprints: []string{"摄像头", "webcam", "camera"}},
	{device: "router", fingerprints: []string{"路由器", "router"}},
	{device: "gateway", fingerprints: []string{"网关", "防火墙", "gateway", "pfsense"}},
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

// HandleAppFingerprint handles the app fingerprint of the given URL.
// It checks if the fingerprints match any known device types and updates the service, device, and honeypot information accordingly.
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
			honeyInfo := url.Port() + "/" + inputFinger[i]
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
	if len(deviceInfo) > 0 {
		utils.UpdateDeviceInfo(url.Hostname(), deviceInfo)
	}
	if len(honeyPot) > 0 {
		utils.UpdateHoneypot(url.Hostname(), honeyPot)
	}
}
