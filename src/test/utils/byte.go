package utils

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
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

func GetBytesData(bytes []byte) string {
	var result string
	count := 0
	for _, val := range bytes {
		count++
		if count > 0 && (count%8) == 0 {
			result = result + strconv.Itoa(int(val)) + ", "
		} else {
			result = result + strconv.Itoa(int(val)) + " "
		}
	}
	return result
}

func GetBytesDataHex(bytes []byte) string {
	var result string
	count := 0
	for _, val := range bytes {
		count++
		if count > 0 && (count%8) == 0 {
			result = result + DecToHex(int64(val)) + ", "
		} else {
			result = result + DecToHex(int64(val)) + " "
		}
	}
	return result
}

func PrintStringBytes(val string) {
	fmt.Println(GetBytesData(StringToBytes(val)))
}
