package utils

import (
	"strconv"
	"strings"
)

// Hex, Dec, Bin in int64

func HexArrayToDecArray(strArray []string) []int64 {
	arr := make([]int64, 0, len(strArray))
	for _, str := range strArray {
		val, _ := strconv.ParseInt(str, 16, 64)
		arr = append(arr, val)
	}
	return arr
}

func HexArrayToBinArray(strArray []string) []string {
	arr := make([]string, 0, len(strArray))
	for _, str := range strArray {
		val, _ := strconv.ParseInt(str, 16, 64)
		arr = append(arr, Int64ToBin(val))
	}
	return arr
}

func DecArrayToHexArray(strArray []string) []string {
	arr := make([]string, 0, len(strArray))
	for _, str := range strArray {
		val, _ := strconv.ParseInt(str, 10, 64)
		arr = append(arr, Int64ToHex(val))
	}
	return arr
}

func DecArrayToDecArray(strArray []string) []int64 {
	arr := make([]int64, 0, len(strArray))
	for _, str := range strArray {
		val, _ := strconv.ParseInt(str, 10, 64)
		arr = append(arr, val)
	}
	return arr
}

func DecArrayToBinArray(strArray []string) []string {
	arr := make([]string, 0, len(strArray))
	for _, str := range strArray {
		val, _ := strconv.ParseInt(str, 10, 64)
		arr = append(arr, Int64ToBin(val))
	}
	return arr
}

func BinArrayToHexArray(strArray []string) []string {
	arr := make([]string, 0, len(strArray))
	for _, str := range strArray {
		val, _ := strconv.ParseInt(str, 2, 64)
		arr = append(arr, Int64ToHex(val))
	}
	return arr
}

func BinArrayToDecArray(strArray []string) []int64 {
	arr := make([]int64, 0, len(strArray))
	for _, str := range strArray {
		val, _ := strconv.ParseInt(str, 2, 64)
		arr = append(arr, val)
	}
	return arr
}

// Hex, Dec, Bin in byte

type StrToDecByte func(str string) byte

func HexByteToDecByte(str string) byte {
	val, _ := strconv.ParseInt(str, 16, 64)
	return byte(val)
}

func DecByteToDecByte(str string) byte {
	val, _ := strconv.ParseInt(str, 10, 64)
	return byte(val)
}

func DecInt8ToDecByte(str string) byte {
	val, _ := strconv.ParseInt(str, 10, 64)
	return byte(val)
}

func HexByteArrayToDecByteArray(strArray []string) []byte {
	arr := make([]byte, 0, len(strArray))
	for _, str := range strArray {
		val, _ := strconv.ParseInt(str, 16, 64)
		arr = append(arr, byte(val))
	}
	return arr
}

func HexByteArrayToInt8Array(strArray []string) []int8 {
	arr := make([]int8, 0, len(strArray))
	for _, str := range strArray {
		val, _ := strconv.ParseInt(str, 16, 64)
		arr = append(arr, int8(val))
	}
	return arr
}

func DecByteArrayToHexByteArray(strArray []string) []string {
	arr := make([]string, 0, len(strArray))
	for _, str := range strArray {
		val, _ := strconv.ParseInt(str, 10, 64)
		arr = append(arr, ByteToHex(byte(val)))
	}
	return arr
}

func DecByteArrayToDecByteArray(strArray []string) []byte {
	arr := make([]byte, 0, len(strArray))
	for _, str := range strArray {
		val, _ := strconv.ParseInt(str, 10, 64)
		arr = append(arr, byte(val))
	}
	return arr
}

func DecByteArrayToDecInt8Array(strArray []string) []int8 {
	arr := make([]int8, 0, len(strArray))
	for _, str := range strArray {
		val, _ := strconv.ParseInt(str, 10, 64)
		arr = append(arr, int8(val))
	}
	return arr
}

func DecInt8ArrayToHexByteArray(strArray []string) []string {
	return DecByteArrayToHexByteArray(strArray)
}

func DecInt8ArrayToDecByteArray(strArray []string) []byte {
	arr := make([]byte, 0, len(strArray))
	for _, str := range strArray {
		val, _ := strconv.ParseInt(str, 10, 64)
		arr = append(arr, byte(val))
	}
	return arr
}

func DecInt8ArrayToDecInt8Array(strArray []string) []int8 {
	arr := make([]int8, 0, len(strArray))
	for _, str := range strArray {
		val, _ := strconv.ParseInt(str, 10, 64)
		arr = append(arr, int8(val))
	}
	return arr
}

func SplitInputString(input string) []string {
	strArray := make([]string, 0, 100)
	var val byte
	var builder strings.Builder
	for i := 0; i <= len(input); i++ {
		if i == len(input) {
			val = 0
		} else {
			val = input[i]
		}
		if (val >= '0' && val <= '9') || (val >= 'a' && val <= 'f') || (val >= 'A' && val <= 'F') || val == '-' {
			builder.WriteByte(val)
		} else {
			if builder.Len() > 0 {
				strArray = append(strArray, builder.String())
				builder.Reset()
			}
		}
	}

	return strArray
}
