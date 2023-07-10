package utils

import (
	"fmt"
	"os"
)

func errLog(format string, a ...interface{}) {
	errStr := fmt.Sprintf(format, a...)

	go func() {
		file, err := os.OpenFile("../result/err_log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()

		_, err = file.WriteString(errStr + "\n")
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}()
}
