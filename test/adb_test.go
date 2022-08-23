package test

import (
	"adbutils"
	"fmt"
	"testing"
)

var adb = adbutils.AdbClient{Host: "localhost", Port: 5037, SocketTime: 10}

func TestServerVersion(t *testing.T) {
	version := adb.ServerVersion()
	t.Logf("version: %d", version)
}

func TestConnect(t *testing.T) {
	// adb := adbutils.NewAdb("localhost", 5037, 10)
	snNtid := adbutils.SerialNTransportID{
		Serial: "127.0.0.1:5555",
	}
	// fmt.Println(adb.Device(snNtid).CurrentApp())
	fmt.Println(adb.Device(snNtid).Shell("ls", false, 10))
	//fmt.Println(adbutils.AdbPath(), runtime.GOARCH)
}
