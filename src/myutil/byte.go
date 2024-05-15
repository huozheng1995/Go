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
	toString(byte) string
	getWidth() int
}

type ByteToHexStr struct {
}

func (byteToStr ByteToHexStr) toString(val byte) string {
	if val == 0 {
		return "00"
	}

	arr := make([]byte, 2)
	arr[0] = ByteHexMap[val>>4&0x0F]
	arr[1] = ByteHexMap[val&0x0F]

	return string(arr)
}
func (byteToStr ByteToHexStr) getWidth() int {
	return 2
}

type ByteToByteStr struct {
}

func (byteToStr ByteToByteStr) toString(val byte) string {
	return strconv.Itoa(int(val))
}
func (byteToStr ByteToByteStr) getWidth() int {
	return 3
}

type ByteToInt8Str struct {
}

func (byteToStr ByteToInt8Str) toString(val byte) string {
	return strconv.Itoa(int(int8(val)))
}
func (byteToStr ByteToInt8Str) getWidth() int {
	return 4
}

type ByteToRawBytes struct {
}

func (byteToStr ByteToRawBytes) toString(val byte) string {
	return ""
}
func (byteToStr ByteToRawBytes) getWidth() int {
	return 0
}
