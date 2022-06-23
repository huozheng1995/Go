package utils

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
	"strings"
)

var Endian = binary.BigEndian

func FloatToBytes(val float32) []byte {
	bits := math.Float32bits(val)
	bytes := make([]byte, 4)
	Endian.PutUint32(bytes, bits)
	return bytes
}

func Float64ToBytes(val float64) []byte {
	bits := math.Float64bits(val)
	bytes := make([]byte, 8)
	Endian.PutUint64(bytes, bits)
	return bytes
}

func IntToBytes(val int64) []byte {
	var bytes = make([]byte, 4)
	Endian.PutUint32(bytes, uint32(val))
	return bytes
}

func Int64ToBytes(val int64) []byte {
	var bytes = make([]byte, 8)
	Endian.PutUint64(bytes, uint64(val))
	return bytes
}

func StringToBytes(val string) []byte {
	return []byte(val)
}

func BytesToFloat(bytes []byte) float32 {
	bits := Endian.Uint32(bytes)
	return math.Float32frombits(bits)
}

func BytesToFloat64(bytes []byte) float64 {
	bits := Endian.Uint64(bytes)
	return math.Float64frombits(bits)
}

func BytesToInt(bytes []byte) int {
	return int(Endian.Uint32(bytes))
}

func BytesToInt64(bytes []byte) int64 {
	return int64(Endian.Uint64(bytes))
}

func BytesToString(bytes []byte) string {
	return string(bytes)
}

type BytesDataToNum func(bytes []byte) string

func ByteArrayToInt8Array(bytes []byte) []int8 {
	len := len(bytes)
	var result = make([]int8, len)
	for i := 0; i < len; i++ {
		result[i] = int8(bytes[i])
	}
	return result
}

func ByteArrayToLine(bytes []byte) string {
	var result string
	count := 0
	for _, val := range bytes {
		count++
		if count%8 == 0 {
			result = result + FillSpace(strconv.Itoa(int(val)), 3) + ", "
		} else {
			result = result + FillSpace(strconv.Itoa(int(val)), 3) + " "
		}
	}
	return result
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

func BytesDataToHex(bytes []byte) string {
	var result string
	count := 0
	for _, val := range bytes {
		count++
		if count%8 == 0 {
			result = result + DecToHex(int64(val)) + ", "
		} else {
			result = result + DecToHex(int64(val)) + " "
		}
	}
	return result
}

func PrintStringBytes(val string) {
	fmt.Println(ByteArrayToLine(StringToBytes(val)))
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
