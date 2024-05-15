package myutil

import (
	"strconv"
)

// string to byte

type StrToByte func(str string) byte

func Hex2StrToByte(str string) byte {
	val, _ := strconv.ParseInt(str, 16, 64)
	return byte(val)
}

func ByteStrToByte(str string) byte {
	val, _ := strconv.ParseInt(str, 10, 64)
	return byte(val)
}

func Int8StrToByte(str string) byte {
	val, _ := strconv.ParseInt(str, 10, 64)
	return byte(val)
}

// byte to string

type ByteToStr interface {
	ToString(byte) string
	GetWidth() int
}

type ByteToHexStr struct {
}

func (byteToStr ByteToHexStr) ToString(val byte) string {
	if val == 0 {
		return "00"
	}

	arr := make([]byte, 2)
	arr[0] = ByteHexMap[val>>4&0x0F]
	arr[1] = ByteHexMap[val&0x0F]

	return string(arr)
}
func (byteToStr ByteToHexStr) GetWidth() int {
	return 2
}

type ByteToByteStr struct {
}

func (byteToStr ByteToByteStr) ToString(val byte) string {
	return strconv.Itoa(int(val))
}
func (byteToStr ByteToByteStr) GetWidth() int {
	return 3
}

type ByteToInt8Str struct {
}

func (byteToStr ByteToInt8Str) ToString(val byte) string {
	return strconv.Itoa(int(int8(val)))
}
func (byteToStr ByteToInt8Str) GetWidth() int {
	return 4
}

type ByteToRawBytes struct {
}

func (byteToStr ByteToRawBytes) ToString(val byte) string {
	return ""
}
func (byteToStr ByteToRawBytes) GetWidth() int {
	return 0
}
