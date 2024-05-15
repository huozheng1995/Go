package utils

import (
	"myutil"
	"strings"
)

// text to byte array

func ReqTextToByteArray(str string, funcStrToByte myutil.StrToByte) []byte {
	result := make([]byte, 0, len(str))
	var val byte
	var builder strings.Builder
	for i := 0; i < len(str)+1; i++ {
		if i == len(str) {
			val = 0
		} else {
			val = str[i]
		}
		if (val >= '0' && val <= '9') || (val >= 'a' && val <= 'f') || (val >= 'A' && val <= 'F') || val == '-' {
			builder.WriteByte(val)
		} else {
			if builder.Len() > 0 {
				result = append(result, funcStrToByte(builder.String()))
				builder.Reset()
			}
		}
	}

	return result
}

// text to int64 array

func ReqTextToInt64Array(text string, funcStrToInt64 myutil.StrToInt64) []int64 {
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
