package test

import (
	"adbutils"
	"fmt"
	"testing"
)

func TestConnect(t *testing.T) {
	adb := adbutils.NewAdb("localhost", 5037, 10)
	fmt.Println(adb.Device("127.0.0.1:5555").Shell("pwd", false, 10))
	//fmt.Println(adbutils.AdbPath(), runtime.GOARCH)
}
