package myutil

import (
	"strconv"
	"strings"
)

// string to int64

type StrToInt64 func(str string) int64

func HexStrToInt64(str string) int64 {
	val, _ := strconv.ParseInt(str, 16, 64)
	return val
}

func DecStrToInt64(str string) int64 {
	val, _ := strconv.ParseInt(str, 10, 64)
	return val
}

func BinStrToInt64(str string) int64 {
	val, _ := strconv.ParseInt(str, 2, 64)
	return val
}

// int64 to string

type Int64ToStr interface {
	ToString(int64) string
	GetWidth() int
}

type Int64ToHexStr struct {
}

func (toStr Int64ToHexStr) ToString(val int64) string {
	if val == 0 {
		return "0"
	}

	var tempVal byte
	var builder strings.Builder
	for i := 64; i > 0; i = i - 4 {
		tempVal = byte(val>>(i-4)) & 0x0F
		if tempVal > 0 || builder.Len() > 0 {
			builder.WriteByte(ByteHexMap[tempVal])
		}
	}

	return builder.String()
}
func (toStr Int64ToHexStr) GetWidth() int {
	return 16
}

type Int64ToInt64Str struct {
}

func (toStr Int64ToInt64Str) ToString(val int64) string {
	return strconv.FormatInt(val, 10)
}
func (toStr Int64ToInt64Str) GetWidth() int {
	return 32
}

type Int64ToBinStr struct {
}

func (toStr Int64ToBinStr) ToString(val int64) string {
	if val == 0 {
		return "0"
	}

	var tempVal byte
	var builder strings.Builder
	for i := 64; i > 0; i-- {
		tempVal = byte(val>>(i-1)) & 0x01
		if tempVal > 0 || builder.Len() > 0 {
			builder.WriteByte(ByteBinMap[tempVal])
		}
	}

	return builder.String()
}
func (toStr Int64ToBinStr) GetWidth() int {
	return 64
}
