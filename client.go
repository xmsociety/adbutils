package adbutils

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	OKAY            = "OKAY"
	FAIL            = "FAIL"
	DENT            = "DENT"
	DONE            = "DONE"
	DATA            = "DATA"
	TCP             = "tcp"
	UNIX            = "unix"
	DEV             = "dev"
	LOCAL           = "local"
	LOCALRESERVED   = "localreserved"
	LOCALFILESYSTEM = "localfilesystem"
	LOCALABSTRACT   = "localabstract"
)

func checkServer(host string, port int) bool {
	_, err := net.Dial("tcp", fmt.Sprintf("%v:%v", host, port))
	return err == nil
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func getParentDirectory(dirctory string) string {
	return substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}
func getCurrentFile() string {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		panic(errors.New("Can not get current file info"))
	}
	return getParentDirectory(file)
}

// adbStreamConnection region adbStreamConnection

type adbStreamConnection struct {
	Host string
	Port int
	Conn net.Conn
}

func (adbStream adbStreamConnection) safeConnect(timeOut time.Duration) (*net.Conn, error) {
	conn, err := adbStream.createSocket(timeOut)
	if err != nil {
		switch reflect.TypeOf(err) {
		case reflect.TypeOf(&net.OpError{}):
			cmd := exec.Command(AdbPath(), "start-server")
			err = cmd.Start()
			if err != nil {
				panic(err.Error())
				return nil, err
			}
			err = cmd.Wait()
			if err != nil {
				panic(err.Error())
				return nil, err
			}
			conn, err = adbStream.createSocket(timeOut)
			if err != nil {
				log.Fatal("restart adb error!")
			}
			return conn, nil
		default:
			panic(err.Error())
		}
		return nil, err
	}
	return conn, nil
}

func (adbStream adbStreamConnection) SetTimeout(timeOut time.Duration) error {
	if timeOut != 0 {
		var err error
		err = adbStream.Conn.SetReadDeadline(time.Now().Add(time.Second * timeOut))
		if err != nil {
			panic(err.Error())
			return err
		}
		err = adbStream.Conn.SetWriteDeadline(time.Now().Add(time.Second * timeOut))
		if err != nil {
			panic(err.Error())
			return err
		}
	}
	return nil
}

func (adbStream adbStreamConnection) createSocket(timeOut time.Duration) (*net.Conn, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%d", adbStream.Host, adbStream.Port))
	if err != nil {
		return nil, err
	}
	if timeOut != 0 {
		var err error
		err = conn.SetReadDeadline(time.Now().Add(time.Second * timeOut))
		if err != nil {
			panic(err.Error())
		}
		err = conn.SetWriteDeadline(time.Now().Add(time.Second * timeOut))
		if err != nil {
			panic(err.Error())
		}
	}
	return &conn, nil
}

func (adbStream adbStreamConnection) Close() {
	err := adbStream.Conn.Close()
	if err != nil {
		return
	}
}

func (adbStream adbStreamConnection) Read(n int) []byte {
	return adbStream.readFully(n)
}

func (adbStream adbStreamConnection) readFully(n int) []byte {
	t := 0
	buffer := make([]byte, n)
	result := bytes.NewBuffer(nil)
	for t < n {
		length, err := adbStream.Conn.Read(buffer[0:n])
		if length == 0 {
			break
		}
		result.Write(buffer[0:length])
		t += length
		if err != nil {
			if err == io.EOF {
				break
			}
		}
	}
	return result.Bytes()
}

func (adbStream adbStreamConnection) SendCommand(cmd string) {
	msg := fmt.Sprintf("%04x%s", len(cmd), cmd)
	_, err := adbStream.Conn.Write([]byte(msg))
	if err != nil {
		log.Fatal("write error!")
		return
	}
}

func (adbStream adbStreamConnection) ReadString(n int) string {
	res := adbStream.Read(n)
	return string(res)
}

func (adbStream adbStreamConnection) ReadStringBlock() string {
	str := adbStream.ReadString(4)
	if len(str) == 0 {
		log.Fatal("connection closed")
	}
	size, _ := strconv.ParseUint(str, 16, 32)
	return adbStream.ReadString(int(size))
}

func (adbStream adbStreamConnection) ReadUntilClose() string {
	buf := []byte{}
	for {
		chunk := adbStream.Read(4096)
		if len(chunk) == 0 {
			break
		}

		buf = append(buf, chunk...)
	}
	return string(buf)
}

func (adbStream adbStreamConnection) CheckOkay() {
	data := adbStream.ReadString(4)
	if data == FAIL {
		log.Fatal("connection closed")
	} else if data == OKAY {
		return
	}
	log.Fatal(fmt.Sprintf("Unknown data: %v", data))
}

// end region adbStreamConnection

// region AdbClient
type AdbClient struct {
	Host       string
	Port       int
	SocketTime time.Duration
}

func AdbPath() string {
	currentPath := getCurrentFile()
	var adb = "adb"
	if runtime.GOOS == "windows" {
		adb = "adb.exe"
	}
	abs, err := filepath.Abs(path.Join(currentPath, "binaries", adb))
	if err != nil {
		return ""
	}
	return abs
}

func (adb *AdbClient) connect(timeout time.Duration) *adbStreamConnection {
	adbStream := &adbStreamConnection{
		Host: adb.Host,
		Port: adb.Port,
	}
	conn, err := adbStream.safeConnect(timeout)
	if err != nil {
		panic(err.Error())
	}
	adbStream.Conn = *conn
	return adbStream

}

func (adb *AdbClient) ServerVersion() int {
	c := adb.connect(10)
	c.SendCommand("host:version")
	c.CheckOkay()
	res := c.ReadStringBlock()
	l, _ := strconv.Atoi(res)
	return l + 16
}

func (adb *AdbClient) ServerKill() {
	if checkServer(adb.Host, adb.Port) {
		c := adb.connect(10)
		c.SendCommand("host:kill")
		c.CheckOkay()
	}
}

func (adb *AdbClient) WaitFor() {
	// pass
}

func (adb *AdbClient) Connect(addr string, timeOut time.Duration) string {
	//addr (str): adb remote address [eg: 191.168.0.1:5555]
	c := adb.connect(timeOut)
	c.SendCommand("host:connect:" + addr)
	return c.ReadStringBlock()
}

func (adb *AdbClient) Disconnect(addr string, raiseErr bool) string {
	//addr (str): adb remote address [eg: 191.168.0.1:5555]
	c := adb.connect(10)
	c.SendCommand("host:disconnect:" + addr)
	return c.ReadStringBlock()
}

type SerialNTransportID struct {
	Serial      string
	TransportID int
}

func (adb *AdbClient) Shell(serial string, command string, stream bool, timeout time.Duration) interface{} {
	snNtid := SerialNTransportID{Serial: serial}
	return adb.Device(snNtid).Shell(command, stream, timeout)
}

func (adb *AdbClient) DeviceList() []AdbDevice {
	res := []AdbDevice{}
	c := adb.connect(10)
	c.SendCommand("host:devices")
	c.CheckOkay()
	outPut := c.ReadStringBlock()
	outPuts := strings.Split(outPut, "\n")
	for _, line := range outPuts {
		parts := strings.Split(strings.TrimSpace(line), "\t")
		if len(parts) != 2 {
			continue
		}
		if parts[1] == "device" {
			res = append(res, AdbDevice{ShellMixin{Client: adb, Serial: parts[0]}})
		}
	}
	return res
}

func (adb *AdbClient) Device(snNtid SerialNTransportID) AdbDevice {
	if snNtid.Serial != "" || snNtid.TransportID != 0 {
		return AdbDevice{ShellMixin{Client: adb, Serial: snNtid.Serial, TransportID: snNtid.TransportID}}
	}
	serial := os.Getenv("ANDROID_SERIAL")
	if serial != "" {
		ds := adb.DeviceList()
		if len(ds) == 0 {
			log.Fatal("Error: Can't find any android device/emulator")
		} else if len(ds) > 1 {
			log.Fatal("more than one device/emulator, please specify the serial number")
		} else {
			return ds[0]
		}
	}
	return AdbDevice{ShellMixin{Client: adb, Serial: snNtid.Serial, TransportID: snNtid.TransportID}}
}

func NewAdb(host string, port int, timeOut time.Duration) *AdbClient {
	adb := &AdbClient{Host: host, Port: port, SocketTime: time.Second * timeOut}
	return adb
}

// end region AdbClient
