package cmd

import (
	"testing"
	"time"
)

func TestErrLog(t *testing.T) {
	err := "hello"

	ErrLog("Failed to retrieve data from the database: %v %v %v %v", err, err, err, err)

	ErrLog("hello")
	ErrLog("hello")
	ErrLog("hello")
	ErrLog("hello")
	ErrLog("hello")
	ErrLog("hello")
	time.Sleep(5 * time.Second)
}
