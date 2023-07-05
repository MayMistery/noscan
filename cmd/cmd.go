package cmd

import (
	"flag"
	"github.com/MayMistery/noscan/.version"
	"runtime"
	"runtime/debug"
	"time"
)

func init() {
	go func() {
		for {
			GarbageCollection()
			time.Sleep(10 * time.Second)
		}
	}()
}

func GarbageCollection() {
	runtime.GC()
	debug.FreeOSMemory()
}

func Banner() {
	banner := `
_  _ ____ ____ ____ ____ _  _ 
|\ | |  | [__  |    |__| |\ | 
| \| |__| ___] |___ |  | | \| 
noscan version: ` + version.Version + `
`
	print(banner)
}

func Flag(Info *HostInfo) {
	Banner()
	flag.StringVar(&Proxy, "proxy", "", "set poc proxy, -proxy http://127.0.0.1:8080")
	flag.StringVar(&Socks5Proxy, "socks5", "", "set socks5 proxy, will be used in tcp connection, timeout setting will not work")
	flag.StringVar(&Cookie, "cookie", "", "set poc cookie,-cookie rememberMe=login")
	flag.Int64Var(&WebTimeout, "wt", 5, "Set web timeout")
	flag.BoolVar(&DnsLog, "dns", false, "using dnslog poc")
	flag.IntVar(&PocNum, "num", 20, "poc rate")
	flag.StringVar(&SC, "sc", "", "ms17 shellcode,as -sc add")
	flag.BoolVar(&IsWmi, "wmi", false, "start wmi")
	flag.StringVar(&Hash, "hash", "", "hash")
	flag.Parse()
}
