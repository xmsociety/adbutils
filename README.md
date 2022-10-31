# adbutils
[![GoDoc](https://pkg.go.dev/badge/github.com/xmsociety/adbutils?status.svg)](https://pkg.go.dev/github.com/xmsociety/adbutils?tab=doc)
[![Sourcegraph](https://sourcegraph.com/github.com/xmsociety/adbutils/-/badge.svg)](https://sourcegraph.com/github.com/xmsociety/adbutils?badge)
[![Goproxy.cn](https://goproxy.cn/stats/github.com/xmsociety/adbutils/badges/download-count.svg)]((https://goproxy.cn))

Transfer from [python adbutils](https://github.com/openatx/adbutils)


**Table of Contents**

<!--ts-->
   * [adbutils](#adbutils)
   * [Install](#install)
   * [Usage](#usage)
      * [x] [Connect ADB Server](#connect-adb-server)
      * [x] [List all the devices and get device object](#list-all-the-devices-and-get-device-object)
      * [ ] [Connect remote device](#connect-remote-device)
      * [ ] [adb forward and adb reverse](#adb-forward-and-adb-reverse)
      * [x] [Create socket connection to the device](#create-socket-connection-to-the-device)
      * [x] [Run shell command](#run-shell-command)
      * [ ] [Transfer files](#transfer-files)
      * [ ] [Extended Functions](#extended-functions)
      * [ ] [Run in command line 命令行使用](#run-in-command-line-命令行使用)
         * [x] [Environment variables](#environment-variables)
         * [x] [Color Logcat](#color-logcat)
      * [x] [Experiment](#experiment)
      * [x] [Examples](#examples)
   * [Develop](#develop)
      * [Watch adb socket data](#watch-adb-socket-data)
   * [Thanks](#thanks)
   * [Ref](#ref)
   * [LICENSE](#license)

<!-- Added by: shengxiang, at: 2021年 3月26日 星期五 15时05分04秒 CST -->

<!--te-->

# Install
- No development plan yet

# Usage
Example

## Connect ADB Server
```go
package test

import (
	"fmt"
	"testing"

	"github.com/xmsociety/adbutils"
)

var adb = adbutils.AdbClient{Host: "localhost", Port: 5037, SocketTime: 10}

func TestConnect(t *testing.T) {
	for _, i := range adb.DeviceList() {
		adb.Connect(i.Serial)
		snNtid := adbutils.SerialNTransportID{
			Serial: i.Serial,
		}
		fmt.Println(adb.Device(snNtid).SayHello())
	}

}
```

## List all the devices and get device object
```go
package main

import (
	"fmt"

	"github.com/xmsociety/adbutils"
)
adb := adbutils.NewAdb("localhost", 5037, 10)

func ShowSerials() {
	for _, device := range adb.DeviceList() {
		fmt.Println("", device.Serial)
	}
}

type SerialNTransportID struct {
    // you get this struct by adbutils.SerialNTransportID
	Serial      string
	TransportID int
}

just_serial := SerialNTransportID{Serial: "33ff22xx"}
adb.Device(just_serial)

// or
just_transport_id := SerialNTransportID{TransportID: 24}
adb.Device(just_transport_id) // transport_id can be found in: adb devices -l

// # You do not need to offer serial if only one device connected
// # RuntimeError will be raised if multi device connected
// d = adb.Device()
```

The following code will not write `from adbutils import adb` for short

## Connect or disconnect remote device
Same as command `adb connect`

```go
output := adb.Connect("127.0.0.1:5555")
// output: already connected to 127.0.0.1:5555

# connect with timeout
// timeout 10
adb := adbutils.NewAdb("localhost", 5037, 10)


// Disconnect
adb.Disconnect("127.0.0.1:5555")
adb.Disconnect("127.0.0.1:5555", raise_error=True) # if device is not present, AdbError will raise

// wait-for-device
// TODO
```

## Create socket connection to the device

For example

```go
func (adbConnection AdbConnection) ReadString(n int) string {
	res := adbConnection.Read(n)
	return string(res)
}

func (adbConnection AdbConnection) ReadStringBlock() string {
	str := adbConnection.ReadString(4)
	if len(str) == 0 {
		log.Fatal("receive data error connection closed")
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
```

```go
func (adbConnection AdbConnection) createSocket() (*net.Conn, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%d", adbConnection.Host, adbConnection.Port))
	if err != nil {
		return nil, err
	}
	return &conn, nil
}
```

There are many other usage, see [SERVICES.TXT](https://cs.android.com/android/platform/superproject/+/master:packages/modules/adb/SERVICES.TXT;l=175) for more details

## Run shell command
I assume there is only one device connected.

```go
package main

import (
	"fmt"

	"github.com/xmsociety/adbutils"
)

// 获取序列号
func GetServerVersion() {
	adb := adbutils.AdbClient{Host: "localhost", Port: 5037, SocketTime: 10}
	version := adb.ServerVersion()
	fmt.Printf("version: %d\n\n", version)
}

func Shell(arg string) {
	adb := adbutils.NewAdb("localhost", 5037, 10)
	for _, device := range adb.DeviceList() {
		fmt.Printf("Now show device: %s, ls: \n", device.Serial)
		fmt.Println(device.Shell(arg, false))
	}
}

func main() {
	GetServerVersion()
	Shell("ls")
}

```
**Other You Can Send:**
- Argument just support str
`Shell(["getprop", "ro.serial"])` - can't work

- Can't Set timeout for shell command
`Shell("sleep 1", timeout=0.5)` - Recommend you set timeout by adb's socketTime

- The advanced shell (returncode archieved by add command suffix: ;echo EXIT:$?)
```go
ret := device.Shell("echo 1")
fmt.Println(ret)
```

- show property, also based on d.shell
TODO


### Environment variables

```bash
ANDROID_SERIAL  serial number to connect to
ANDROID_ADB_SERVER_HOST adb server host to connect to
ANDROID_ADB_SERVER_PORT adb server port to connect to
```

### Color Logcat
- No development plan yet


## Experiment
TODO
<!-- Install Auto confirm supported(Beta), you need to famillar with [uiautomator2](https://github.com/openatx/uiautomator2) first

```bash
# Install with auto confirm (Experiment, based on github.com/openatx/uiautomator2)
$ python -m adbutils --install-confirm -i some.apk
``` -->

For more usage, please see the code for details.

## Examples
Record video using screenrecord

It is highly recommended that you follow [this](https://github.com/xmsociety/go-scrcpy-client)
```go
// TODO
```
<!-- stream = d.shell("screenrecord /sdcard/s.mp4", stream=True)
time.sleep(3) # record for 3 seconds
with stream:
	stream.send(b"\003") # send Ctrl+C
	stream.read_until_close()

start = time.time()
print("Video total time is about", time.time() - start)
d.sync.pull("/sdcard/s.mp4", "s.mp4") # pulling video -->

Reading Logcat

```go
// TODO
```
<!-- d.shell("logcat --clear")
stream = d.shell("logcat", stream=True)
with stream:
    f = stream.conn.makefile()
    for _ in range(100): # read 100 lines
        line = f.readline()
        print("Logcat:", line.rstrip())
    f.close() -->

# Develop
Make sure you can connect Github, Now you can edit code in `adbutils` and test with

```go
package test
import (
	"github.com/xmsociety/adbutils"
	"testing"
)
// .... test code here ...
```

Run tests requires one device connected to your computer

```sh
# change to repo directory
cd adbutils

go test test/*
```

# Environment
Some environment can affect the adbutils behavior

- ADBUTILS_ADB_PATH: specify adb path, default search from PATH
- ANDROID_SERIAL: default adb serial
- ANDROID_ADB_SERVER_HOST: default 127.0.0.1
- ANDROID_ADB_SERVER_PORT: default 5037

## Watch adb socket data
Watch the adb socket data using `socat`

```
$ socat -t100 -x -v TCP-LISTEN:5577,reuseaddr,fork TCP4:localhost:5037
```

open another terminal, type the following command then you will see the socket data

```bash
$ export ANDROID_ADB_SERVER_PORT=5577
$ adb devices
```

## Generate TOC
```bash
gh-md-toc --insert README.md
```

<https://github.com/ekalinin/github-markdown-toc>

# Thanks
- [python adbutils](https://github.com/openatx/adbutils)
- [swind pure-python-adb](https://github.com/Swind/pure-python-adb)
- [openstf/adbkit](https://github.com/openstf/adbkit)
- [ADB Source Code](https://github.com/aosp-mirror/platform_system_core/blob/master/adb)
- ADB Protocols [OVERVIEW.TXT](https://cs.android.com/android/platform/superproject/+/master:packages/modules/adb/OVERVIEW.TXT) [SERVICES.TXT](https://cs.android.com/android/platform/superproject/+/master:packages/modules/adb/SERVICES.TXT) [SYNC.TXT](https://cs.android.com/android/platform/superproject/+/master:packages/modules/adb/SYNC.TXT)
- [Awesome ADB](https://github.com/mzlogin/awesome-adb)
- [JakeWharton/pidcat](https://github.com/JakeWharton/pidcat)

# Alternative
just like [pure-python-adb](https://github.com/Swind/pure-python-adb)

![Alternative](https://raw.githubusercontent.com/Swind/pure-python-adb/master/docs/adb_pure_python_adb.png "Alternative png")

# Ref
- <https://github.com/imageio/imageio-ffmpeg/blob/80e37882d0/imageio_ffmpeg/_utils.py>

# LICENSE
[MIT](LICENSE)


