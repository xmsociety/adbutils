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
		Serial: "emulator-5554",
	}
	adb.Connect("emulator-5554")
	fmt.Println(adb.Device(snNtid).AdbOut("push /Users/sato/Desktop/go-scrcpy-client/scrcpy/scrcpy-server.jar /data/local/tmp/scrcpy-server.jar"))

}
