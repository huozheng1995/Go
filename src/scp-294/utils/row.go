package utils

import (
	"strings"
)

func ByteArrayToRow(byteToStr ByteToStr, arr []byte, off int, len int, withDetails bool) string {
	var result strings.Builder
	count := 0
	for i := 0; i < len; i++ {
		count++
		result.WriteString(FillSpace(byteToStr.toString(arr[off+i]), byteToStr.getWidth()))
		if withDetails {
			if count&0x0F == 0 {
				result.WriteString(", ")
			} else {
				result.WriteString(" ")
			}
		} else {
			result.WriteString(", ")
		}
	}
	return result.String()
}

func ByteArrayToRowDetails(arr []byte, off int, len int) string {
	var result strings.Builder
	count := 0
	for i := off; i < off+len; i++ {
		count++
		if arr[i] >= 32 && arr[i] <= 126 {
			result.WriteByte(arr[i])
		} else {
			result.WriteString(CharNULL)
		}
		if count%8 == 0 {
			result.WriteString(", ")
		}
	}
	return result.String()
}

func FillChar(val string, expectedLen int, char rune) string {
	diff := expectedLen - len(val)
	if diff > 0 {
		return strings.Repeat(string(char), diff) + val
	} else {
		return val
	}
}

func FillSpace(val string, expectedLen int) string {
	return FillChar(val, expectedLen, ' ')
}

func Fill0(val string, expectedLen int) string {
	return FillChar(val, expectedLen, '0')
}
