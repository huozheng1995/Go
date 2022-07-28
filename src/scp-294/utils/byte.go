package utils

import (
	"encoding/binary"
	"math"
	"strconv"
	"strings"
)

var Endian = binary.BigEndian

func FloatToBytes(val float32) []byte {
	bits := math.Float32bits(val)
	arr := make([]byte, 4)
	Endian.PutUint32(arr, bits)
	return arr
}

func Float64ToBytes(val float64) []byte {
	bits := math.Float64bits(val)
	arr := make([]byte, 8)
	Endian.PutUint64(arr, bits)
	return arr
}

func IntToBytes(val int64) []byte {
	var arr = make([]byte, 4)
	Endian.PutUint32(arr, uint32(val))
	return arr
}

func Int64ToBytes(val int64) []byte {
	var arr = make([]byte, 8)
	Endian.PutUint64(arr, uint64(val))
	return arr
}

func StringToBytes(val string) []byte {
	return []byte(val)
}

func BytesToFloat(arr []byte) float32 {
	bits := Endian.Uint32(arr)
	return math.Float32frombits(bits)
}

func BytesToFloat64(arr []byte) float64 {
	bits := Endian.Uint64(arr)
	return math.Float64frombits(bits)
}

func BytesToInt(arr []byte) int {
	return int(Endian.Uint32(arr))
}

func BytesToInt64(arr []byte) int64 {
	return int64(Endian.Uint64(arr))
}

func BytesToString(arr []byte) string {
	return string(arr)
}

type BytesDataToNum func(arr []byte) string

func ByteArrayToInt8Array(arr []byte) []int8 {
	len := len(arr)
	var result = make([]int8, len)
	for i := 0; i < len; i++ {
		result[i] = int8(arr[i])
	}
	return result
}

func Int8ArrayToByteArray(arr []int8) []byte {
	len := len(arr)
	var result = make([]byte, len)
	for i := 0; i < len; i++ {
		result[i] = byte(arr[i])
	}
	return result
}

func ByteArrayToLine(arr []byte, off int, len int) string {
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

func TwoByteArraysToLine(arr1 []byte, off1 int, len1 int, arr2 []byte, off2 int, len2 int) string {
	var result strings.Builder
	count := 0
	if arr1 != nil {
		for i := off1; i < off1+len1; i++ {
			count++
			if count%8 == 0 {
				result.WriteString(FillSpace(strconv.Itoa(int(arr1[i])), 3))
				result.WriteString(", ")
			} else {
				result.WriteString(FillSpace(strconv.Itoa(int(arr1[i])), 3))
				result.WriteString(" ")
			}
		}
	}
	if arr2 != nil {
		for i := off2; i < off2+len2; i++ {
			count++
			if count%8 == 0 {
				result.WriteString(FillSpace(strconv.Itoa(int(arr1[i])), 3))
				result.WriteString(", ")
			} else {
				result.WriteString(FillSpace(strconv.Itoa(int(arr1[i])), 3))
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

func TwoBytesArrayToCharLine(arr1 []byte, off1 int, len1 int, arr2 []byte, off2 int, len2 int) string {
	var result strings.Builder
	count := 0

	if arr1 != nil {
		for i := off1; i < off1+len1; i++ {
			count++
			if arr1[i] >= 32 && arr1[i] <= 126 {
				result.WriteByte(arr1[i])
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
	if arr2 != nil {
		for i := off2; i < off2+len2; i++ {
			count++
			if arr1[i] >= 32 && arr1[i] <= 126 {
				result.WriteByte(arr1[i])
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

func ByteArrayToHex(arr []byte) string {
	var result string
	count := 0
	for _, val := range arr {
		count++
		if count%8 == 0 {
			result = result + DecToHex(int64(val)) + ", "
		} else {
			result = result + DecToHex(int64(val)) + " "
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
