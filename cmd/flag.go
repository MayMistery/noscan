package cmd

import (
	"flag"
	version "github.com/MayMistery/noscan/.version"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"
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
	_, ex, _, _ := runtime.Caller(0)
	exPath := filepath.Dir(filepath.Dir(ex))
	var timeout int

	Banner()
	flag.StringVar(&config.CIDR, "CIDR", "default", "multi CIDR to scan, split by \",\"")
	flag.StringVar(&config.Ports, "port", "common", "Ports to scan, choices: all | common | <port>[-<port>][,<ports>]")
	flag.BoolVar(&config.JsonOutput, "json", true, "using json output")
	flag.BoolVar(&config.Ciscn, "d", true, "to complete Ciscn task")
	flag.StringVar(&config.InputFilepath, "input", path.Join(exPath, "data/target"), "input file path")
	flag.StringVar(&config.OutputFilepath, "output", path.Join(exPath, "result/result.json"), "output file path")
	flag.StringVar(&config.RulesFilePath, "rule", path.Join(exPath, "data/fingerprint.txt"), "rule file path")
	flag.StringVar(&config.DBFilePath, "db", path.Join(exPath, "data/database.db"), "database file path")
	flag.StringVar(&config.ScanType, "t", "tcp", "scan method, tcp | syn | fin | NULL")
	flag.BoolVar(&config.Ping, "ping", false, "use system ping method to check host alive")
	flag.StringVar(&config.Proxy, "proxy", "", "proxy")
	flag.IntVar(&config.Threads, "thread", 1000, "threads num")
	flag.IntVar(&timeout, "timeout", 3, "connect time out")
	flag.BoolVar(&config.DeepInspection, "simple", true, "close deep identify")
	flag.BoolVar(&config.Debug, "debug", false, "open go pprof")
	flag.BoolVar(&config.help, "help", false, "print help info")

	config.Timeout = time.Duration(timeout) * time.Second
	// TODO to add flags and corresponding var

	flag.Parse()

	if config.help == true {
		flag.PrintDefaults()
		os.Exit(0)
		// maybe need to exit
	}
}
