package adbutils

import "time"

type DeviceEvent struct {
	Present bool
	Serial  string
	Status  string
}

type ForwardItem struct {
	Serial string
	Local  string
	Remote string
}

type ReverseItem struct {
	Remote string
	Local  string
}

type FileInfo struct {
	Mode  int
	Size  int
	Mtime *time.Time
	Path  string
}
type WindowSize struct {
	Width  int
	Height int
}

type RunningAppInfo struct {
	Package  string
	Activity string
	Pid      int
}

type ShellReturn struct {
	Command    string
	ReturnCode int
	Output     string
}

type AdbDeviceInfo struct {
	Serial string
	State  string
}
