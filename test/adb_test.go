package test

import (
	"adbutils"
	"testing"
)

func TestConnect(t *testing.T) {
	adbutils.NewAdb("127.0.0.1", 5037, 10)
}
