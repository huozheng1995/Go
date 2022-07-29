package utils

import (
	"encoding/binary"
	"strconv"
	"strings"
)

var Endian = binary.BigEndian
var HexMap = map[byte]byte{0: 48, 1: 49, 2: 50, 3: 51, 4: 52, 5: 53, 6: 54, 7: 55, 8: 56, 9: 57, 10: 65, 11: 66, 12: 67, 13: 68, 14: 69, 15: 70}

func ByteToHex(val byte) string {
	if val == 0 {
		return "00"
	}
	s := ""
	for q := val; q > 0; q = q >> 4 {
		m := q % 16
		s = string(HexMap[m]) + s
	}
	if val>>4 == 0 {
		return "0" + s
	}
	return s
}

type BytesToLine func(arr []byte, off int, len int) string

func BytesToByteLine(arr []byte, off int, len int) string {
	var result strings.Builder
	count := 0
	if arr != nil {
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
	}
	return result.String()
}

func BytesToInt8Line(arr []byte, off int, len int) string {
	var result strings.Builder
	count := 0
	if arr != nil {
		for i := off; i < off+len; i++ {
			count++
			if count%8 == 0 {
				result.WriteString(FillSpace(strconv.Itoa(int(int8(arr[i]))), 3))
				result.WriteString(", ")
			} else {
				result.WriteString(FillSpace(strconv.Itoa(int(int8(arr[i]))), 3))
				result.WriteString(" ")
			}
		}
	}
	return result.String()
}

func BytesToHexLine(arr []byte, off int, len int) string {
	var result strings.Builder
	count := 0
	if arr != nil {
		for i := off; i < off+len; i++ {
			count++
			if count%8 == 0 {
				result.WriteString(FillSpace(ByteToHex(arr[i]), 3))
				result.WriteString(", ")
			} else {
				result.WriteString(FillSpace(ByteToHex(arr[i]), 3))
				result.WriteString(" ")
			}
		}
	}
	return result.String()
}

func Int8ArrayToLine(arr []int8) string {
	var result string
	count := 0
	for _, val := range arr {
		count++
		if count%8 == 0 {
			result = result + FillSpace(strconv.Itoa(int(val)), 4) + ", "
		} else {
			result = result + FillSpace(strconv.Itoa(int(val)), 4) + " "
		}
	}
	return result
}

func HexByteArrayToLine(arr []string) string {
	var result string
	count := 0
	for _, val := range arr {
		count++
		if count%8 == 0 {
			result = result + FillSpace(val, 3) + ", "
		} else {
			result = result + FillSpace(val, 3) + " "
		}
	}
	return result
}

func ByteArrayToCharLine(arr []byte, off int, len int) string {
	var result strings.Builder
	count := 0
	if arr != nil {
		for i := off; i < off+len; i++ {
			count++
			if arr[i] >= 32 && arr[i] <= 126 {
				result.WriteByte(arr[i])
			} else {
				result.WriteByte(0)
			}
			if count%8 == 0 {
				result.WriteString(", ")
			} else {
				result.WriteString(" ")
			}
		}
	}
	return result.String()
}

func Int8ArrayToCharLine(arr []int8) string {
	var result string
	count := 0
	for _, val := range arr {
		var charVal string
		if val >= 32 && val <= 126 {
			charVal = string(val)
		} else {
			charVal = string(0)
		}
		count++
		if count%8 == 0 {
			result = result + charVal + ", "
		} else {
			result = result + charVal + " "
		}
	}
	return result
}

func HexByteArrayToCharLine(arr []string) string {
	var result string
	count := 0
	for _, val := range arr {
		var byteVal = byte(HexToDec(val))
		var charVal string
		if byteVal >= 32 && byteVal <= 126 {
			charVal = string(byteVal)
		} else {
			charVal = string(0)
		}
		count++
		if count%8 == 0 {
			result = result + charVal + ", "
		} else {
			result = result + charVal + " "
		}
	}
	return result
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
