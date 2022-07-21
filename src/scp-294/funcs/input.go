package funcs

import (
	"github.com/edward/scp-294/utils"
	"strconv"
	"strings"
)

func HexArrayToDecArray(input string) []int64 {
	strArray := SplitInputString(input)
	arr := make([]int64, 0, len(strArray))
	for _, str := range strArray {
		arr = append(arr, utils.HexToDec(str))
	}
	return arr
}

func HexByteArrayToDecByteArray(input string) []byte {
	strArray := SplitInputString(input)
	arr := make([]byte, 0, len(strArray))
	for _, str := range strArray {
		arr = append(arr, byte(utils.HexToDec(str)))
	}
	return arr
}

func HexByteArrayToInt8Array(input string) []int8 {
	strArray := SplitInputString(input)
	arr := make([]int8, 0, len(strArray))
	for _, str := range strArray {
		arr = append(arr, int8(utils.HexToDec(str)))
	}
	return arr
}

func DecArrayToHexArray(input string) []string {
	strArray := SplitInputString(input)
	arr := make([]string, 0, len(strArray))
	for _, str := range strArray {
		val, _ := strconv.ParseInt(str, 10, 64)
		arr = append(arr, utils.DecToHex(val))
	}
	return arr
}

func DecByteArrayToHexByteArray(input string) []string {
	strArray := SplitInputString(input)
	arr := make([]string, 0, len(strArray))
	for _, str := range strArray {
		val, _ := strconv.ParseInt(str, 10, 8)
		arr = append(arr, utils.DecToHex(val))
	}
	return arr
}

func Int8ArrayToHexByteArray(input string) []string {
	strArray := SplitInputString(input)
	arr := make([]string, 0, len(strArray))
	for _, str := range strArray {
		val, _ := strconv.ParseInt(str, 10, 8)
		arr = append(arr, utils.DecToHex(val))
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
