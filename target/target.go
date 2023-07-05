package target

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	// 读取target文件
	data, err := os.ReadFile("target")
	if err != nil {
		fmt.Println("无法读取文件：", err)
		return
	}

	// 将文件内容分割为字符串数组
	ipAddresses := strings.Split(string(data), "\n")

	// 打印IP地址
	for _, ip := range ipAddresses {
		fmt.Println(ip)
	}
}
