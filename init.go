package adbutils

import (
	"fmt"
	"net"
	"os/exec"
	"reflect"
	"time"
)

const (
	TCP             = "tcp"
	UNIX            = "unix"
	DEV             = "dev"
	LOCAL           = "local"
	LOCALRESERVED   = "localreserved"
	LOCALFILESYSTEM = "localfilesystem"
	LOCALABSTRACT   = "localabstract"
)

type adbStreamConnection struct {
	Host string
	Port int
	Conn net.Conn
}

func (adbStream adbStreamConnection) safeConnect(timeOut time.Duration) (net.Conn, error) {
	conn, err := adbStream.createSocket(timeOut)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (adbStream adbStreamConnection) createSocket(timeOut time.Duration) (net.Conn, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", adbStream.Host, adbStream.Port))
	if err != nil {
		return nil, err
	}
	if timeOut != 0 {
		var err error
		err = conn.SetReadDeadline(time.Now().Add(time.Second * timeOut))
		if err != nil {
			panic(err.Error())
		}
		err = conn.SetWriteDeadline(time.Now())
		if err != nil {
			panic(err.Error())
		}
	}
	return conn, nil
}

type adbClient struct {
	Host       string
	Port       int
	SocketTime time.Time
}

func AdbPath() {

}
func (adb adbClient) connect(timeout time.Duration) {
	adbStream := adbStreamConnection{
		Host: adb.Host,
		Port: adb.Port,
	}
	conn, err := adbStream.safeConnect(timeout)
	if err != nil {
		switch reflect.TypeOf(err) {
		case reflect.TypeOf(&net.OpError{}):
			// TODO
			cmd := exec.Command("adb_path", "start-server")
			err = cmd.Start()
			if err != nil {
				panic(err.Error())
				return
			}
			err := cmd.Wait()
			if err != nil {
				panic(err.Error())
				return
			}
		default:
			panic(err.Error())
		}
		return
	}
	adbStream.Conn = conn

}

func NewAdb(host string, port int, timeOut time.Duration) *adbClient {
	adb := &adbClient{Host: host, Port: port}
	adb.connect(timeOut)
	return adb
}
