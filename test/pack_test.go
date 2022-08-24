package test

import (
	"encoding/binary"
	"fmt"
	"github.com/xmsociety/adbutils"
	"testing"
)

func TestPackAndUnpack(t *testing.T) {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(1111))
	fmt.Printf("%#v\n", bs)
	i := binary.LittleEndian.Uint32(bs)
	fmt.Println(i)
}

func TestPackAndUnpackMuti(t *testing.T) {
	bytes := []byte{}
	for i := 0; i <= 3; i++ {
		bs := make([]byte, 4)
		binary.LittleEndian.PutUint32(bs, uint32(5))
		//fmt.Printf("%#v\n", bs)
		bytes = append(bytes, bs...)
	}
	fmt.Println(bytes)
	for i := 0; i <= 3; i++ {
		item := binary.LittleEndian.Uint32(bytes[i*4 : (i+1)*4])
		fmt.Println(item)
	}
}

func TestGetFreePort(t *testing.T) {
	fmt.Println(adbutils.GetFreePort())
}
