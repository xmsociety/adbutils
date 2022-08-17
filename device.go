package adbutils

import (
	"net"
	"strings"
	"time"
)

type ShellMixin struct{}

type AdbDevice struct {
	Client *adbClient
	Serial string
}

func (device AdbDevice) getWithCommand(cmd string) string {
	c := device.Client.connect(10)
	c.SendCommand(strings.Join([]string{"host-serial", device.Serial, cmd}, ":"))
	c.CheckOkay()
	return c.ReadStringBlock()
}

func (device AdbDevice) GetState() string {
	return device.getWithCommand("get-state")
}

func (device AdbDevice) Shell(cmd string, stream bool, timeOut time.Duration) interface{} {
	ret := device.Client.Shell(device.Serial, cmd, stream, timeOut)
	if stream {
		return ret
	}
	return ret
}

func (device AdbDevice) ShellOutPut(cmd string) string {
	res := device.Client.Shell(device.Serial, cmd, false, 0)
	return res.(string)
}

func (device AdbDevice) Push(local, remote string) {
	// pass
}

func (device AdbDevice) CreateConnection(netWork, address string) net.Conn {
	c := device.Client.connect(10)
	c.SendCommand("host:transport:" + device.Serial)
	c.CheckOkay()
	switch netWork {
	case TCP:
		c.SendCommand("tcp" + address)
		c.CheckOkay()
	case UNIX, LOCALABSTRACT:
		c.SendCommand("localabstract" + address)
		c.CheckOkay()
	case LOCALFILESYSTEM, LOCAL, DEV, LOCALRESERVED:
		c.SendCommand(netWork + ":" + address)
		c.CheckOkay()
	default:
		panic("not support net work: " + netWork)
	}
	return c.Conn
}
