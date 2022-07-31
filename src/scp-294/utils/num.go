package utils

import (
	"encoding/binary"
)

var Endian = binary.BigEndian
var ByteHexMap = map[byte]byte{0: 48, 1: 49, 2: 50, 3: 51, 4: 52, 5: 53, 6: 54, 7: 55, 8: 56, 9: 57, 10: 65, 11: 66, 12: 67, 13: 68, 14: 69, 15: 70}
var ByteBinMap = map[byte]byte{0: 48, 1: 49}
var nullByte = byte(126)

func ByteToHex(val byte) string {
	if val == 0 {
		return "00"
	}
	result := ""
	for tempVal := val; tempVal > 0; tempVal = tempVal >> 4 {
		result = string(ByteHexMap[tempVal%16]) + result
	}
	if len(result) == 1 {
		return "0" + result
	}
	return result
}

func DecToHex(int64Val int64) string {
	val := uint64(int64Val)
	if val == 0 {
		return "00"
	}
	result := ""
	for tempVal := val; tempVal > 0; tempVal = tempVal >> 4 {
		result = string(ByteHexMap[byte(tempVal%16)]) + result
	}
	if len(result) == 1 {
		return "0" + result
	}
	return result
}

func DecToBin(val int64) string {
	if val == 0 {
		return "0"
	}
	result := ""
	for tempVal := val; tempVal > 0; tempVal = tempVal >> 1 {
		result = string(ByteBinMap[byte(tempVal%2)]) + result
	}
	return result
}
