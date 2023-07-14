package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"
)

// ErrLog Uniformly output error messages to a log file
func ErrLog(format string, a ...interface{}) {
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05") + " : "
	errStr := formattedTime + fmt.Sprintf(format, a...)

	go func() {
		//TODO change filepath to flag
		_, ex, _, _ := runtime.Caller(0)
		exPath := path.Join(filepath.Dir(ex), "..")
		errorLogPath := path.Join(exPath, "result/err_log")
		file, err := os.OpenFile(errorLogPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			fmt.Println("Error opening file:", err, errStr)
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

// ResultLog Log successful scan events
func ResultLog(format string, a ...interface{}) {
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05") + " : "
	errStr := formattedTime + fmt.Sprintf(format, a...)

	go func() {
		//TODO change filepath to flag
		_, ex, _, _ := runtime.Caller(0)
		exPath := path.Join(filepath.Dir(ex), "..")
		resultLogPath := path.Join(exPath, "result/result_log")
		file, err := os.OpenFile(resultLogPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			fmt.Println("Error opening file:", err, errStr)
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
