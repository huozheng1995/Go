package myutil

import (
	"strconv"
	"strings"
)

type Int64Util interface {
	ToNum(string) int64
	ToString(int64) string
	GetDisplaySize() int
}

type HexUtil struct {
}

func (u HexUtil) ToNum(str string) int64 {
	val, _ := strconv.ParseInt(str, 16, 64)
	return val
}
func (u HexUtil) ToString(val int64) string {
	if val == 0 {
		return "0"
	}

	var tempVal byte
	var builder strings.Builder
	for i := 64; i > 0; i = i - 4 {
		tempVal = byte(val>>(i-4)) & 0x0F
		if tempVal > 0 || builder.Len() > 0 {
			builder.WriteByte(ByteHexMap[tempVal])
		}
	}

	return builder.String()
}
func (u HexUtil) GetDisplaySize() int {
	return 16
}

type DecUtil struct {
}

func (u DecUtil) ToNum(str string) int64 {
	val, _ := strconv.ParseInt(str, 10, 64)
	return val
}
func (u DecUtil) ToString(val int64) string {
	return strconv.FormatInt(val, 10)
}
func (u DecUtil) GetDisplaySize() int {
	return 32
}

type BinUtil struct {
}

func (u BinUtil) ToNum(str string) int64 {
	val, _ := strconv.ParseInt(str, 2, 64)
	return val
}
func (u BinUtil) ToString(val int64) string {
	if val == 0 {
		return "0"
	}

	var tempVal byte
	var builder strings.Builder
	for i := 64; i > 0; i-- {
		tempVal = byte(val>>(i-1)) & 0x01
		if tempVal > 0 || builder.Len() > 0 {
			builder.WriteByte(ByteBinMap[tempVal])
		}
	}

	return builder.String()
}
func (u BinUtil) GetDisplaySize() int {
	return 64
}
