package main

import (
	"strconv"
)

func ByteToBit(val byte) string {
	var result string
	for i := 0; i < 8; i++ {
		var t int
		t = (val >> (7 - i)) & 0x01
		result += strconv.Itoa(t)
	}
	return result
}

func BytesToBit(bytes []byte) string {
	var result string
	for _, val := range bytes {
		result += ByteToBit(val) + " "
	}
	return result
}

func FloatToBit(val float32) string {
	bytes := FloatToBytes(val)
	return BytesToBit(bytes)
}

func Float64ToBit(val float64) string {
	bytes := Float64ToBytes(val)
	return BytesToBit(bytes)
}

func IntToBit(val int) string {
	var result string
	for i := 0; i < 32; i++ {
		var t int
		t = (val >> (31 - i)) & 0x01
		result += strconv.Itoa(t)
	}
	return result
}

func Int64ToBit(val int64) string {
	var result string
	for i := 0; i < 64; i++ {
		var t int
		t = (val >> (63 - i)) & 0x01
		result += strconv.Itoa(t)
	}
	return result
}
