package test

import (
	"fmt"
	"testing"

	"github.com/xmsociety/adbutils"
)

var adb = adbutils.AdbClient{Host: "localhost", Port: 5037, SocketTime: 10}

func TestServerVersion(t *testing.T) {
	version := adb.ServerVersion()
	t.Logf("version: %d", version)
}

func TestConnect(t *testing.T) {
	// adb := adbutils.NewAdb("localhost", 5037, 10)
	for _, i := range adb.DeviceList() {
		adb.Connect(i.Serial)
		snNtid := adbutils.SerialNTransportID{
			Serial: i.Serial,
		}
		fmt.Println(adb.Device(snNtid).SayHello())
		// fmt.Println(adb.Device(snNtid).Push("/Users/sato/Desktop/go-scrcpy-client/scrcpy/scrcpy-server.jar", "/data/local/tmp/scrcpy-server.jar"))
	}

}
