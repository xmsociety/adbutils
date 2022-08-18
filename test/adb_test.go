package test

import (
	"adbutils"
	"fmt"
	"testing"
)

func TestConnect(t *testing.T) {
	adb := adbutils.NewAdb("localhost", 5037, 10)
	snNtid := adbutils.SerialNTransportID{
		Serial: "emulator-5558",
	}
	// fmt.Println(adb.Device(snNtid).CurrentApp())
	fmt.Println(adb.Device(snNtid).Shell("ls", false, 10))
	//fmt.Println(adbutils.AdbPath(), runtime.GOARCH)
}
