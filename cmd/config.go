package cmd

import (
	"time"
)

type PortInfo struct {
	Port       int      `json:"port"`
	Protocol   string   `json:"protocol"`
	ServiceApp []string `json:"service_app"`
}

type IpInfo struct {
	Services   []*PortInfo `json:"services"`
	DeviceInfo string      `json:"deviceinfo"`
	Honeypot   []string    `json:"honeypot"`
	Timestamp  string      `json:"timestamp"`
}

type Configs struct {
	Ports          string
	JsonOutput     bool
	Ciscn          bool
	InputFilepath  string
	OutputFilepath string
	DBFilePath     string
	RulesFilePath  string
	ScanType       string
	Ping           bool
	Proxy          string
	Threads        int
	Timeout        time.Duration
	DeepInspection bool
	Debug          bool
	CIDR           string
	help           bool
}

var portList = map[string]int{
	"ftp":     21,
	"ssh":     22,
	"findnet": 135,
	"netbios": 139,
	"smb":     445,
	"mssql":   1433,
	"oracle":  1521,
	"mysql":   3306,
	"rdp":     3389,
	"psql":    5432,
	"redis":   6379,
	"fcgi":    9000,
	"mem":     11211,
	"mgo":     27017,
}

var (
	Config      Configs
	IPPools     func() string
	IPPoolsSize int64 = 0
	IPNetPools  []IPPool
	Ports       []int
)
