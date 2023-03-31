package utils

import (
	"strconv"
	"strings"
)

// input

type StrToByte func(str string) byte

func HexStrToByte(str string) byte {
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

// output

type ByteToStr func(val byte) string

func ByteToHexStr(val byte) string {
	if val == 0 {
		return "00"
	}

	var builder strings.Builder
	builder.WriteByte(ByteHexMap[val>>4&0x0F])
	builder.WriteByte(ByteHexMap[val&0x0F])

	return builder.String()
}

func ByteToByteStr(val byte) string {
	return strconv.Itoa(int(val))
}

func ByteToInt8Str(val byte) string {
	return strconv.Itoa(int(int8(val)))
}
