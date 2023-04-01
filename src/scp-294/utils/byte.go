package utils

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// to byte

type StrToByte func(str string) byte

func HexStrToByte(str string) byte {
	val, _ := strconv.ParseInt(str, 16, 64)
	return byte(val)
}

func ByteStrToByte(str string) byte {
	val, _ := strconv.ParseInt(str, 10, 64)
	return byte(val)
}

func Int8StrToByte(str string) byte {
	val, _ := strconv.ParseInt(str, 10, 64)
	return byte(val)
}

// byte to

type ByteToStr interface {
	toString(byte) string
	getWidth() int
}

type ByteToHexStr struct {
}

func (byteToStr ByteToHexStr) toString(val byte) string {
	if val == 0 {
		return "00"
	}

	var builder strings.Builder
	builder.WriteByte(ByteHexMap[val>>4&0x0F])
	builder.WriteByte(ByteHexMap[val&0x0F])

	return builder.String()
}
func (byteToStr ByteToHexStr) getWidth() int {
	return 2
}

type ByteToByteStr struct {
}

func (byteToStr ByteToByteStr) toString(val byte) string {
	return strconv.Itoa(int(val))
}
func (byteToStr ByteToByteStr) getWidth() int {
	return 3
}

type ByteToInt8Str struct {
}

func (byteToStr ByteToInt8Str) toString(val byte) string {
	return strconv.Itoa(int(int8(val)))
}
func (byteToStr ByteToInt8Str) getWidth() int {
	return 4
}

// to byte array

func StringToByteArray(text string, funcStrToByte StrToByte) []byte {
	result := make([]byte, 0, 65536)
	var val byte
	var builder strings.Builder
	for i := 0; i < len(text)+1; i++ {
		if i == len(text) {
			val = 0
		} else {
			val = text[i]
		}
		if (val >= '0' && val <= '9') || (val >= 'a' && val <= 'f') || (val >= 'A' && val <= 'F') || val == '-' {
			builder.WriteByte(val)
		} else {
			if builder.Len() > 0 {
				result = append(result, funcStrToByte(builder.String()))
				builder.Reset()
			}
		}
	}

	return result
}

// byte array to

const printLen = 5
const GlobalRowSize = 16

func ByteArrayToRowBytes(arr []byte, globalRowIndex *int, byteToStr ByteToStr, withDetails bool) []byte {
	rowSize := GlobalRowSize

	totalLen := len(arr)
	totalRow := totalLen / rowSize
	lastRowCount := totalLen % rowSize
	if lastRowCount > 0 {
		totalRow++
	}

	buffer := new(bytes.Buffer)
	buffer.Grow(len(arr) * 10)
	for rowIndex := 0; rowIndex < totalRow; rowIndex++ {
		globalByteIndex := *globalRowIndex * rowSize
		byteIndex := rowIndex * rowSize
		if rowIndex == totalRow-1 && lastRowCount > 0 {
			rowSize = lastRowCount
		}
		if withDetails {
			buffer.WriteString(fmt.Sprintf("Row%s(%s, %s): %s        %s\n",
				Fill0(strconv.Itoa(*globalRowIndex), printLen-1),
				Fill0(strconv.Itoa(globalByteIndex), printLen),
				Fill0(strconv.Itoa(globalByteIndex+8), printLen),
				ByteArrayToRow(byteToStr, arr, byteIndex, rowSize, withDetails),
				ByteArrayToRowDetails(arr, byteIndex, rowSize)))
		} else {
			buffer.WriteString(ByteArrayToRow(byteToStr, arr, byteIndex, rowSize, withDetails))
			buffer.WriteString("\n")
		}
		*globalRowIndex++
	}
	return buffer.Bytes()
}

func ByteArrayToRowString(arr []byte, globalRowIndex *int, byteToStr ByteToStr, withDetails bool) string {
	rowSize := GlobalRowSize

	totalLen := len(arr)
	totalRow := totalLen / rowSize
	lastRowCount := totalLen % rowSize
	if lastRowCount > 0 {
		totalRow++
	}

	var builder strings.Builder
	for rowIndex := 0; rowIndex < totalRow; rowIndex++ {
		globalByteIndex := *globalRowIndex * rowSize
		byteIndex := rowIndex * rowSize
		if rowIndex == totalRow-1 && lastRowCount > 0 {
			rowSize = lastRowCount
		}
		if withDetails {
			builder.WriteString(fmt.Sprintf("Row%s(%s, %s): %s        %s\n",
				Fill0(strconv.Itoa(*globalRowIndex), printLen-1),
				Fill0(strconv.Itoa(globalByteIndex), printLen),
				Fill0(strconv.Itoa(globalByteIndex+8), printLen),
				ByteArrayToRow(byteToStr, arr, byteIndex, rowSize, withDetails),
				ByteArrayToRowDetails(arr, byteIndex, rowSize)))
		} else {
			builder.WriteString(ByteArrayToRow(byteToStr, arr, byteIndex, rowSize, withDetails))
			builder.WriteString("\n")
		}
		*globalRowIndex++
	}
	return builder.String()
}
