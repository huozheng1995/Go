package myutil

import (
	"strconv"
)

type ByteUtil interface {
	ToNum(string) byte
	ToString(byte) string
	GetDisplaySize() int
}

type Hex8Util struct{}

func (u Hex8Util) ToNum(str string) byte {
	val, _ := strconv.ParseInt(str, 16, 64)
	return byte(val)
}
func (u Hex8Util) ToString(val byte) string {
	if val == 0 {
		return "00"
	}

	arr := make([]byte, 2)
	arr[0] = ByteHexMap[val>>4&0x0F]
	arr[1] = ByteHexMap[val&0x0F]

	return string(arr)
}
func (u Hex8Util) GetDisplaySize() int {
	return 2
}

type Byte8Util struct{}

func (u Byte8Util) ToNum(str string) byte {
	val, _ := strconv.ParseInt(str, 10, 64)
	return byte(val)
}
func (u Byte8Util) ToString(val byte) string {
	return strconv.Itoa(int(val))
}
func (u Byte8Util) GetDisplaySize() int {
	return 3
}

type Int8Util struct{}

func (u Int8Util) ToNum(str string) byte {
	val, _ := strconv.ParseInt(str, 10, 64)
	return byte(val)
}
func (u Int8Util) ToString(val byte) string {
	return strconv.Itoa(int(int8(val)))
}
func (u Int8Util) GetDisplaySize() int {
	return 4
}

type RawBytesUtil struct{}

func (u RawBytesUtil) ToNum(str string) byte {
	return 0
}
func (u RawBytesUtil) ToString(val byte) string {
	return ""
}
func (u RawBytesUtil) GetDisplaySize() int {
	return 0
}
