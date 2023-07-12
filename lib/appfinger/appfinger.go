package appfinger

import (
	"errors"
	"github.com/MayMistery/noscan/lib/appfinger/gorpc"
	httpfinger2 "github.com/MayMistery/noscan/lib/appfinger/httpfinger"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var supportProtocols = []string{
	"http",
	"https",
	"rpc",
}

var supportProtocolRegx = regexp.MustCompile("^" + strings.Join(supportProtocols, "|") + "$")

func Search(URL *url.URL, banner *Banner) *FingerPrint {
	var products, hostnames []string
	switch URL.Scheme {
	case "http":
		products = search((*httpfinger2.Banner)(convHttpBanner(URL, banner)))
		return &FingerPrint{products, "", "", ""}
	case "https":
		products := search((*httpfinger2.Banner)(convHttpBanner(URL, banner)))
		return &FingerPrint{products, "", "", ""}
	case "rpc":
		hostnames, _ = gorpc.GetHostname(URL.Hostname())
		return &FingerPrint{emptyProductName, strings.Join(hostnames, ";"), "", ""}
	}
	return nil
}

func SupportCheck(protocol string) bool {
	return supportProtocolRegx.MatchString(protocol)
}

func GetBannerWithResponse(URL *url.URL, response string, req *http.Request, cli *http.Client) (*Banner, error) {
	switch URL.Scheme {
	case "http":
		httpBanner, err := httpfinger2.GetBannerWithResponse(URL, response, req, cli)
		return convBanner(httpBanner), err
	case "https":
		httpBanner, err := httpfinger2.GetBannerWithResponse(URL, response, req, cli)
		return convBanner(httpBanner), err
	default:
		return convBannerWithRaw(response), nil
	}
}

func GetBannerWithURL(URL *url.URL, req *http.Request, cli *http.Client) (*Banner, error) {
	switch URL.Scheme {
	case "http":
		httpBanner, err := httpfinger2.GetBannerWithURL(URL, req, cli)
		return convBanner(httpBanner), err
	case "https":
		httpBanner, err := httpfinger2.GetBannerWithURL(URL, req, cli)
		return convBanner(httpBanner), err
	}
	return nil, errors.New("unsupported protocol")
}

func convHttpBanner(URL *url.URL, banner *Banner) *httpfinger2.Banner {
	return &httpfinger2.Banner{
		Protocol: URL.Scheme,
		Port:     URL.Port(),
		Header:   banner.Header,
		Body:     banner.Body,
		Response: banner.Response,
		Cert:     banner.Cert,
		Title:    banner.Title,
		Hash:     banner.Hash,
		Icon:     banner.Icon,
	}
}
