/*
Copyright 2023 noname

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under tgo he License.
*/

package main

import (
	"fmt"
	"github.com/MayMistery/noscan/cmd"
	"github.com/MayMistery/noscan/lib/appfinger"
	"github.com/MayMistery/noscan/scan"
	"github.com/MayMistery/noscan/storage/bolt"
	"github.com/MayMistery/noscan/utils"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// export GODEBUG=detectdeadlocks=1
func main() {
	start := time.Now()
	Exec()
	t := time.Now().Sub(start)
	fmt.Printf("[*] Task done, Duration: %s\n", t)
}

func Exec() {
	handleSIGINT()

	cmd.Flag(&cmd.Config)
	dbWaitGroup := bolt.InitAsyncDatabase()
	utils.InitResultMap()

	defer func() {
		bolt.CloseDatabase(dbWaitGroup)
		dbWaitGroup.Wait()
	}()

	// if the Debug flag is set, start a goroutine that serves runtime profiling data via HTTP
	if cmd.Config.Debug {
		go func() {
			log.Println(http.ListenAndServe(":38899", nil))
		}()
	}

	// open the rules file
	fs, _ := os.Open(cmd.Config.RulesFilePath)
	if n, err := appfinger.InitDatabaseFS(fs); err != nil {
		fmt.Println("[-]指纹库加载失败，请检查【fingerprint.txt】文件", err)
	} else {
		fmt.Printf("[+]成功加载HTTP指纹:[%d]条\n", n)
	}

	// start scanning
	fmt.Println("[+]Start Scan, default output to /app/result/")
	err := scan.Scan()
	if err != nil {
		fmt.Println("[-]Scan error", err)
	}

	// output result
	handleOutput()

	fmt.Println("[+]Bye")
}

func handleSIGINT() {
	// 创建一个通道来接收信号
	sigs := make(chan os.Signal, 1)

	// 注册需要捕获的信号
	signal.Notify(sigs, syscall.SIGINT)

	go func() {
		// 阻塞直到信号被接收到
		sig := <-sigs
		fmt.Println(sig)
		fmt.Println("You pressed ctrl+C, Please have a wait.")
		fmt.Println("Press again to force quit")

		//TODO gracefully exit
		handleOutput()

		fmt.Println("[+]Bye (interrupted)")
		os.Exit(0)
	}()
}

func handleOutput() {
	if cmd.Config.JsonOutput {
		utils.OutputResultMap()
	} else {
		//TODO add terminal output
	}
}
