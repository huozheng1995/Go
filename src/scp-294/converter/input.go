package converter

import (
	"github.com/edward/scp-294/utils"
	"strconv"
	"strings"
)

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
		arr = append(arr, utils.DecToBin(val))
	}
	return arr
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
		arr = append(arr, int8(byte(val)))
	}
	return arr
}

func DecArrayToHexArray(strArray []string) []string {
	arr := make([]string, 0, len(strArray))
	for _, str := range strArray {
		val, _ := strconv.ParseInt(str, 10, 64)
		arr = append(arr, utils.DecToHex(val))
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
		arr = append(arr, utils.DecToBin(val))
	}
	return arr
}

func DecByteArrayToHexByteArray(strArray []string) []string {
	arr := make([]string, 0, len(strArray))
	for _, str := range strArray {
		val, _ := strconv.ParseInt(str, 10, 64)
		arr = append(arr, utils.DecToHex(val))
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

func DecByteArrayToInt8Array(strArray []string) []int8 {
	arr := make([]int8, 0, len(strArray))
	for _, str := range strArray {
		val, _ := strconv.ParseInt(str, 10, 64)
		arr = append(arr, int8(byte(val)))
	}
	return arr
}

func Int8ArrayToHexByteArray(strArray []string) []string {
	return DecByteArrayToHexByteArray(strArray)
}

func Int8ArrayToDecByteArray(strArray []string) []byte {
	arr := make([]byte, 0, len(strArray))
	for _, str := range strArray {
		val, _ := strconv.ParseInt(str, 10, 64)
		arr = append(arr, byte(val))
	}
	return arr
}

func Int8ArrayToInt8Array(strArray []string) []int8 {
	arr := make([]int8, 0, len(strArray))
	for _, str := range strArray {
		val, _ := strconv.ParseInt(str, 10, 64)
		arr = append(arr, int8(byte(val)))
	}
	return arr
}

func BinArrayToHexArray(strArray []string) []string {
	arr := make([]string, 0, len(strArray))
	for _, str := range strArray {
		val, _ := strconv.ParseInt(str, 2, 64)
		arr = append(arr, utils.DecToHex(val))
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

func SplitInputString(input string) []string {
	strArray := make([]string, 0, 100)
	var val byte
	var builder strings.Builder
	for i := 0; i < len(input)+1; i++ {
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
