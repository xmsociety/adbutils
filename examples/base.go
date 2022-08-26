/**
Test localhost Only one Device
**/
package main

import (
	"fmt"

	"github.com/xmsociety/adbutils"
)

func GetServerVersion() {
	adb := adbutils.AdbClient{Host: "localhost", Port: 5037, SocketTime: 10}
	version := adb.ServerVersion()
	fmt.Printf("version: %d\n\n", version)
}

func Shell(arg string) {
	adb := adbutils.NewAdb("localhost", 5037, 10)
	for _, device := range adb.DeviceList() {
		fmt.Printf("Now show device: %s, ls: \n", device.Serial)
		fmt.Printf("Now show device: %s, ls: \n", device.Properties)
		fmt.Println(device.Shell(arg, false))
	}
}

func main() {
	GetServerVersion()
	Shell("ls")
}
