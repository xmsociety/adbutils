package adbutils

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

type ShellMixin struct {
	Client      *adbClient
	Serial      string
	TransportID int
	properties  map[string]string
}

func (mixin ShellMixin) run(cmd string) interface{} {
	return mixin.Client.Shell(mixin.Serial, cmd, false, 10)
}

func (mixin ShellMixin) SayHello() string {
	content := "hello from " + mixin.Serial
	res := mixin.run("echo" + content)
	return res.(string)
}

func (mixin ShellMixin) SwitchScreen(status bool) {
	KeyMap := map[bool]string{
		true:  "224",
		false: "223",
	}
	mixin.KeyEvent(KeyMap[status])
}

func (mixin ShellMixin) SwitchAirPlane(status bool) {
	base := "settings put global airplane_mode_on"
	am := "am broadcast -a android.intent.action.AIRPLANE_MODE --ez state"
	if status {
		base += "1"
		am += "true"
	} else {
		base += "0"
		am += "false"
	}
	mixin.run(base)
	mixin.run(am)
}

func (mixin ShellMixin) SwitchWifi(status bool) {
	cmdMap := map[bool]string{
		true:  "svc wifi enable",
		false: "svc wifi disable",
	}
	mixin.run(cmdMap[status])
}

func (mixin ShellMixin) KeyEvent(keyCode string) string {
	res := mixin.run("input keyevent " + keyCode)
	return res.(string)
}

func (mixin ShellMixin) CLick(x, y int) {
	mixin.run(fmt.Sprintf("input tap %v %v", x, y))
}

func (mixin ShellMixin) Swipe(x, y, tox, toy int, duration time.Duration) {
	mixin.run(fmt.Sprintf("input swipe %v %v %v %v %v", x, y, tox, toy, duration*1000))
}

func (mixin ShellMixin) SendKeys(text string) {
	// TODO escapeSpecialCharacters
	mixin.run("input text " + text)
}

func (mixin ShellMixin) escapeSpecialCharacters(text string) {}

func (mixin ShellMixin) WlanIp() string {
	res := mixin.run("ifconfig wlan0")
	ipInfo := res.(string)
	// TODO regrex
	return ipInfo
}

func (mixin ShellMixin) install(pathOrUrl string, noLaunch bool, unInstall bool, silent bool, callBack func()) {
}

func (mixin ShellMixin) InstallRemote(remotePath string, clean bool) {
	res := mixin.run("pm install -r -t " + remotePath)
	resInfo := res.(string)
	if !strings.Contains(resInfo, "Success") {
		log.Fatalln(resInfo)
	}
	if clean {
		mixin.run("rm " + remotePath)
	}
}

func (mixin ShellMixin) Uninstall(packageName string) {
	mixin.run("pm uninstall " + packageName)
}

func (mixin ShellMixin) GetProp(prop string) string {
	res := mixin.run("getprop " + prop)
	return strings.TrimSpace(res.(string))
}

func (mixin ShellMixin) ListPackages() []string {
	result := []string{}
	res := mixin.run("pm list packages")
	output := res.(string)
	for _, packageName := range strings.Split(output, "\n") {
		p := strings.TrimSpace(strings.TrimPrefix(packageName, "package:"))
		if p == "" {
			continue
		}
		result = append(result, p)
	}
	return result
}

func (mixin ShellMixin) PackageInfo(packageName string) {
	// TODO
}

func (mixin ShellMixin) Rotation() {}

func (mixin ShellMixin) rawWindowSize() {}

func (mixin ShellMixin) WindowSize() {}

func (mixin ShellMixin) AppStart(packageName, activity string) {
	if activity != "" {
		mixin.run("am start -n " + packageName + "/" + activity)
	} else {
		mixin.run("monkey -p " + packageName + "-c" + "android.intent.category.LAUNCHER 1")
	}
}

func (mixin ShellMixin) AppStop(packageName string) {
	mixin.run("am force-stop " + packageName)
}

func (mixin ShellMixin) AppClear(packageName string) {
	mixin.run("pm clear " + packageName)
}

func (mixin ShellMixin) IsScreenOn() bool {
	res := mixin.run("dumpsys power")
	output := res.(string)
	return strings.Contains(output, "mHoldingDisplaySuspendBlocker=true")
}

func (mixin ShellMixin) OpenBrowser(url string) {
	mixin.run("am start -a android.intent.action.VIEW -d " + url)
}

func (mixin ShellMixin) DumpHierarchy() string {
	return ""
}

func (mixin ShellMixin) CurrentApp() string {
	return ""
}

func (mixin ShellMixin) Remove(path string) {
	mixin.run("rm " + path)
}

type AdbDevice struct {
	ShellMixin
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

// TODO chage adbStreamConnection 2 AdbConnection
func (device AdbDevice) openTransport(command string, timeout time.Duration) *adbStreamConnection {
	return &adbStreamConnection{}
}

func (device AdbDevice) Shell(command string, stream bool, timeout time.Duration) interface{} {
	c := device.openTransport(command, timeout)
	if stream {
		c = device.Client.connect(0)
	} else {
		c = device.Client.connect(timeout)
	}
	c.SendCommand("host:transport:" + command)
	c.CheckOkay()
	c.SendCommand("shell:" + command)
	c.CheckOkay()
	if stream {
		return c
	} else {
		return c.ReadUntilClose()
	}
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
