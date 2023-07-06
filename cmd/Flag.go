package cmd

import (
	"flag"
	version "github.com/MayMistery/noscan/.version"
	"os"
	"path"
	"path/filepath"
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

func Flag(config *Configs) {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	Banner()
	//flag.StringVar(&, "h", "", "IP address of the host you want to scan,for example: 192.168.11.11 | 192.168.11.11-255 | 192.168.11.11,192.168.11.12")
	//flag.Int64Var(&Timeout, "time", 3, "Set timeout")
	flag.BoolVar(&config.JsonOutput, "json", true, "using json output")
	flag.BoolVar(&config.Ciscn, "d", true, "to complete Ciscn task")
	flag.StringVar(&config.InputFilepath, "input", path.Join(exPath, "data/target"), "input file path")
	flag.StringVar(&config.OutputFilepath, "output", path.Join(exPath, "result/result.json"), "output file path")
	flag.StringVar(&config.ScanType, "t", "tcp", "scan method, tcp | syn | fin | NULL")
	flag.StringVar(&config.DBFilePath, "db", path.Join(exPath, "data/database.db"), "database file path")

	// TODO to add flags and corresponding var

	flag.PrintDefaults()

	flag.Parse()
}
