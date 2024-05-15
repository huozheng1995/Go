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

type Int64ToStr func(val int64) string

func Int64ToHexStr(val int64) string {
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

func Int64ToInt64Str(val int64) string {
	return strconv.FormatInt(val, 10)
}

func Int64ToBinStr(val int64) string {
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
