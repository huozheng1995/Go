package utils

import "strings"

func HexByteArrayToDecByteArray(input string) []byte {
	strArray := SplitInputString(input)
	byteArray := make([]byte, 0, len(strArray))
	for _, str := range strArray {
		byteArray = append(byteArray, Int64ToBytes(HexToInt64(str))[7])
	}
	return byteArray
}

func HexArrayToDecArray(input string) (decArr []int64) {
	strArray := SplitInputString(input)
	decArray := make([]int64, 0, len(strArray))
	for _, str := range strArray {
		decArray = append(decArray, HexToInt64(str))
	}
	return decArray
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
