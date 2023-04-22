package utils

import (
	"bytes"
	"strconv"
	"strings"
)

// string to int64

type StrToInt64 func(str string) int64

func HexStrToInt64(str string) int64 {
	val, _ := strconv.ParseInt(str, 16, 64)
	return val
}

func Int64StrToInt64(str string) int64 {
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

// request string to int64 array

func ReqStrToInt64Array(text string, funcStrToInt64 StrToInt64) []int64 {
	result := make([]int64, 0, 4096)
	var val byte
	var builder strings.Builder
	for i := 0; i < len(text)+1; i++ {
		if i == len(text) {
			val = 0
		} else {
			val = text[i]
		}
		if (val >= '0' && val <= '9') || (val >= 'a' && val <= 'f') || (val >= 'A' && val <= 'F') || val == '-' {
			builder.WriteByte(val)
		} else {
			if builder.Len() > 0 {
				result = append(result, funcStrToInt64(builder.String()))
				builder.Reset()
			}
		}
	}

	return result
}

// int64 array to response

func Int64ArrayToResponse(arr []int64, funcInt64ToStr Int64ToStr, resBuf *bytes.Buffer) {
	for _, val := range arr {
		resBuf.WriteString(funcInt64ToStr(val))
		resBuf.WriteString(", ")
	}
}
