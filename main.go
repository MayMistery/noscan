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
	"github.com/MayMistery/noscan/scan"
	"github.com/MayMistery/noscan/storage/bolt"
	"github.com/MayMistery/noscan/utils"
	"time"
)

func main() {
	start := time.Now()
	Exec()
	t := time.Now().Sub(start)
	fmt.Printf("[*] Task done, Duration: %s\n", t)
}

func Exec() {

	cmd.Flag(&cmd.Config)
	bolt.InitDatabase()

	err := scan.Scan()
	if err != nil {
		//TODO Handle error
	}
	if cmd.Config.JsonOutput {
		utils.OutputJsonResult()
	} else {
		//TODO add terminal output
	}
}
