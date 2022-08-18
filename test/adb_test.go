package test

import (
	"adbutils"
	"fmt"
	"testing"
)

func TestConnect(t *testing.T) {
	adb := adbutils.NewAdb("localhost", 5037, 10)
	fmt.Println(adb.Device("").CurrentApp())
	//fmt.Println(adbutils.AdbPath(), runtime.GOARCH)
}
