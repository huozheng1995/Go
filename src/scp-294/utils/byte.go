package utils

import (
	"bytes"
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

	arr := make([]byte, 2)
	arr[0] = ByteHexMap[val>>4&0x0F]
	arr[1] = ByteHexMap[val&0x0F]

	return string(arr)
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

func StringToByteArray(str string, funcStrToByte StrToByte) []byte {
	result := make([]byte, 0, len(str))
	var val byte
	var builder strings.Builder
	for i := 0; i < len(str)+1; i++ {
		if i == len(str) {
			val = 0
		} else {
			val = str[i]
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

func ByteArrayToOutputBytes(arr []byte, globalRowIndex *int, byteToStr ByteToStr, withDetails bool) []byte {
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
			buffer.WriteString("Row")
			AppendStringWith0(buffer, strconv.Itoa(*globalRowIndex), printLen-1)
			buffer.WriteString("(")
			AppendStringWith0(buffer, strconv.Itoa(globalByteIndex), printLen)
			buffer.WriteString(", ")
			AppendStringWith0(buffer, strconv.Itoa(globalByteIndex+8), printLen)
			buffer.WriteString("): ")
			buffer.Write(ByteArrayToRowBytes(byteToStr, arr, byteIndex, rowSize, withDetails))
			buffer.WriteString("        ")
			buffer.Write(ByteArrayToRowDetailsBytes(arr, byteIndex, rowSize))
			buffer.WriteString("\n")
		} else {
			buffer.Write(ByteArrayToRowBytes(byteToStr, arr, byteIndex, rowSize, withDetails))
			buffer.WriteString("\n")
		}
		*globalRowIndex++
	}
	return buffer.Bytes()
}

func ByteArrayToOutputString(arr []byte, globalRowIndex *int, byteToStr ByteToStr, withDetails bool) string {
	return string(ByteArrayToOutputBytes(arr, globalRowIndex, byteToStr, withDetails))
}
