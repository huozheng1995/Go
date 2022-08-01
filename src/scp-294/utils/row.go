package utils

import (
	"strconv"
	"strings"
)

type BytesToRow func(arr []byte, off int, len int) string

func BytesToByteRow(arr []byte, off int, len int) string {
	var result strings.Builder
	count := 0
	for i := off; i < off+len; i++ {
		count++
		if count%8 == 0 {
			result.WriteString(FillSpace(strconv.Itoa(int(arr[i])), 3))
			result.WriteString(", ")
		} else {
			result.WriteString(FillSpace(strconv.Itoa(int(arr[i])), 3))
			result.WriteString(" ")
		}
	}
	return result.String()
}

func BytesToInt8Row(arr []byte, off int, len int) string {
	var result strings.Builder
	count := 0
	for i := off; i < off+len; i++ {
		count++
		if count%8 == 0 {
			result.WriteString(FillSpace(strconv.Itoa(int(int8(arr[i]))), 4))
			result.WriteString(", ")
		} else {
			result.WriteString(FillSpace(strconv.Itoa(int(int8(arr[i]))), 4))
			result.WriteString(" ")
		}
	}
	return result.String()
}

func BytesToHexRow(arr []byte, off int, len int) string {
	var result strings.Builder
	count := 0
	for i := off; i < off+len; i++ {
		count++
		if count%8 == 0 {
			result.WriteString(FillSpace(ByteToHex(arr[i]), 2))
			result.WriteString(", ")
		} else {
			result.WriteString(FillSpace(ByteToHex(arr[i]), 2))
			result.WriteString(" ")
		}
	}
	return result.String()
}

func Int8ArrayToRow(arr []int8, off int, len int) string {
	var result strings.Builder
	count := 0
	for i := off; i < off+len; i++ {
		count++
		if count%8 == 0 {
			result.WriteString(FillSpace(strconv.Itoa(int(arr[i])), 4))
			result.WriteString(", ")
		} else {
			result.WriteString(FillSpace(strconv.Itoa(int(arr[i])), 4))
			result.WriteString(" ")
		}
	}
	return result.String()
}

func HexByteArrayToRow(arr []string, off int, len int) string {
	var result strings.Builder
	count := 0
	for i := off; i < off+len; i++ {
		count++
		if count%8 == 0 {
			result.WriteString(FillSpace(arr[i], 2))
			result.WriteString(", ")
		} else {
			result.WriteString(FillSpace(arr[i], 2))
			result.WriteString(" ")
		}
	}
	return result.String()
}

func ByteArrayToCharRow(arr []byte, off int, len int) string {
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

func Int8ArrayToCharRow(arr []int8, off int, len int) string {
	var result strings.Builder
	count := 0
	for i := off; i < off+len; i++ {
		count++
		if arr[i] >= 32 && arr[i] <= 126 {
			result.WriteByte(byte(arr[i]))
		} else {
			result.WriteString(CharNULL)
		}
		if count%8 == 0 {
			result.WriteString(", ")
		}
	}
	return result.String()
}

func HexByteArrayToCharRow(arr []string, off int, len int) string {
	var result strings.Builder
	count := 0
	for i := off; i < off+len; i++ {
		count++
		val, _ := strconv.ParseInt(arr[i], 16, 64)
		var byteVal = byte(val)
		if byteVal >= 32 && byteVal <= 126 {
			result.WriteByte(byteVal)
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