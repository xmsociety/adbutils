package adbutils

import (
	"bytes"
	"errors"
	"fmt"
	_ "github.com/xmsociety/adbutils/binaries"
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

func GetFreePort() int {
	conn, err := net.Listen("tcp", "127.0.0.1:0")
	defer conn.Close()
	if err != nil {
		log.Println("getFreePort error! ", err.Error())
		return 0
	}
	ipPort := strings.Split(conn.Addr().String(), ":")
	port, _ := strconv.Atoi(ipPort[len(ipPort)-1])
	return port
}

// AdbConnection region AdbConnection

type AdbConnection struct {
	Host string
	Port int
	Conn net.Conn
}

func (adbConnection AdbConnection) safeConnect() (*net.Conn, error) {
	conn, err := adbConnection.createSocket()
	if err != nil {
		switch reflect.TypeOf(err) {
		case reflect.TypeOf(&net.OpError{}):
			cmd := exec.Command(AdbPath(), "start-server")
			err = cmd.Start()
			if err != nil {
				log.Println("start adb error: ", err.Error())
				return nil, err
			}
			err = cmd.Wait()
			if err != nil {
				log.Println("start adb error: ", err.Error())
				return nil, err
			}
			conn, err = adbConnection.createSocket()
			if err != nil {
				log.Println("restart adb error! ", err.Error())
				return nil, err
			}
			return conn, nil
		default:
			log.Println("unknown error! ", err.Error())
			return nil, err
		}
		return nil, err
	}
	return conn, nil
}

func (adbConnection AdbConnection) SetTimeout(timeOut time.Duration) error {
	if timeOut != 0 {
		var err error
		err = adbConnection.Conn.SetDeadline(time.Now().Add(time.Second * timeOut))
		if err != nil {
			panic(err.Error())
			return err
		}
	}
	return nil
}

func (adbConnection AdbConnection) createSocket() (*net.Conn, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%d", adbConnection.Host, adbConnection.Port))
	if err != nil {
		return nil, err
	}
	return &conn, nil
}

func (adbConnection AdbConnection) Close() {
	err := adbConnection.Conn.Close()
	if err != nil {
		return
	}
}

func (adbConnection AdbConnection) Read(n int) []byte {
	return adbConnection.readFully(n)
}

func (adbConnection AdbConnection) readFully(n int) []byte {
	t := 0
	buffer := make([]byte, n)
	result := bytes.NewBuffer(nil)
	for t < n {
		length, err := adbConnection.Conn.Read(buffer[0:n])
		if err != nil {
			if err == io.EOF {
				break
			}
			break
		}
		if length == 0 {
			break
		}
		result.Write(buffer[0:length])
		t += length
	}
	return result.Bytes()
}

func (adbConnection AdbConnection) SendCommand(cmd string) {
	msg := fmt.Sprintf("%04x%s", len(cmd), cmd)
	_, err := adbConnection.Conn.Write([]byte(msg))
	if err != nil {
		log.Println("write error!", err.Error())
		return
	}
}

func (adbConnection AdbConnection) ReadString(n int) string {
	res := adbConnection.Read(n)
	return string(res)
}

func (adbConnection AdbConnection) ReadStringBlock() string {
	str := adbConnection.ReadString(4)
	if len(str) == 0 {
		log.Println("receive data error connection closed")
	}
	size, _ := strconv.ParseUint(str, 16, 32)
	return adbConnection.ReadString(int(size))
}

func (adbConnection AdbConnection) ReadUntilClose() string {
	buf := []byte{}
	for {
		chunk := adbConnection.Read(4096)
		if len(chunk) == 0 {
			break
		}

		buf = append(buf, chunk...)
	}
	return string(buf)
}

func (adbConnection AdbConnection) CheckOkay() {
	data := adbConnection.ReadString(4)
	if data == FAIL {
		log.Println(fmt.Sprintf("receive data: %v connection closed", data))
	} else if data == OKAY {
		return
	}
	log.Println(fmt.Sprintf("Unknown data: %v", data))
}

// end region AdbConnection

// AdbClient region AdbClient
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

func (adb *AdbClient) connect() *AdbConnection {
	adbConnection := &AdbConnection{
		Host: adb.Host,
		Port: adb.Port,
	}
	conn, err := adbConnection.safeConnect()
	if err != nil {
		log.Println("get connect error: ", err.Error())
	}
	adbConnection.Conn = *conn
	return adbConnection

}

func (adb *AdbClient) ServerVersion() int {
	c := adb.connect()
	c.SendCommand("host:version")
	c.CheckOkay()
	res := c.ReadStringBlock()
	l, _ := strconv.Atoi(res)
	return l + 16
}

func (adb *AdbClient) ServerKill() {
	if checkServer(adb.Host, adb.Port) {
		c := adb.connect()
		c.SendCommand("host:kill")
		c.CheckOkay()
	}
}

func (adb *AdbClient) WaitFor() {
	// pass
}

func (adb *AdbClient) Connect(addr string) string {
	//addr (str): adb remote address [eg: 191.168.0.1:5555]
	c := adb.connect()
	c.SendCommand("host:connect:" + addr)
	return c.ReadStringBlock()
}

func (adb *AdbClient) Disconnect(addr string, raiseErr bool) string {
	//addr (str): adb remote address [eg: 191.168.0.1:5555]
	c := adb.connect()
	c.SendCommand("host:disconnect:" + addr)
	return c.ReadStringBlock()
}

type SerialNTransportID struct {
	Serial      string
	TransportID int
}

func (adb *AdbClient) Shell(serial string, command string, stream bool) interface{} {
	snNtid := SerialNTransportID{Serial: serial}
	return adb.Device(snNtid).Shell(command, stream, adb.SocketTime)
}

func (adb *AdbClient) DeviceList() []AdbDevice {
	res := []AdbDevice{}
	c := adb.connect()
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
			log.Println("Error: Can't find any android device/emulator")
		} else if len(ds) > 1 {
			log.Println("more than one device/emulator, please specify the serial number")
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
