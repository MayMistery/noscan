package cmd

import (
	"flag"
	version "github.com/MayMistery/noscan/.version"
)

func Banner() {
	banner := `
_  _ ____ ____ ____ ____ _  _ 
|\ | |  | [__  |    |__| |\ | 
| \| |__| ___] |___ |  | | \| 
noscan version: ` + version.Version + `
`
	print(banner)
}

func Flag(config Configs) {
	Banner()
	//flag.StringVar(&, "h", "", "IP address of the host you want to scan,for example: 192.168.11.11 | 192.168.11.11-255 | 192.168.11.11,192.168.11.12")
	//flag.Int64Var(&Timeout, "time", 3, "Set timeout")
	flag.BoolVar(&config.jsonOutput, "json", true, "using json output")
	flag.BoolVar(&config.ciscn, "d", true, "to complete ciscn task")
	flag.StringVar(&config.filepath, "path", "../result/result.json", "output file path")
	flag.StringVar(&config.ScanType, "t", "tcp", "scan method, tcp | syn | fin | NULL")

	// TODO to add flags and corresponding var

	flag.PrintDefaults()

	flag.Parse()
}
