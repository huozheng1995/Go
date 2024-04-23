package utils

import (
	"bytes"
	"strconv"
	"strings"
)

// string to byte

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

// byte to string

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

type ByteToRawBytes struct {
}

func (byteToStr ByteToRawBytes) toString(val byte) string {
	return ""
}
func (byteToStr ByteToRawBytes) getWidth() int {
	return 0
}

// request string to byte array

func ReqStrToByteArray(str string, funcStrToByte StrToByte) []byte {
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

// byte array to response

const printLen = 5
const GlobalRowSize = 16

func ByteArrayToResponse(arr []byte, globalRowIndex *int, byteToStr ByteToStr, withDetails bool, resBuf *bytes.Buffer) {
	//To RawBytes
	if _, ok := byteToStr.(ByteToRawBytes); ok {
		resBuf.Write(arr)
		return
	}

	//To String
	rowSize := GlobalRowSize
	totalLen := len(arr)
	totalRow := totalLen / rowSize
	lastRowCount := totalLen % rowSize
	if lastRowCount > 0 {
		totalRow++
	}

	for rowIndex := 0; rowIndex < totalRow; rowIndex++ {
		globalByteIndex := *globalRowIndex * rowSize
		byteIndex := rowIndex * rowSize
		if rowIndex == totalRow-1 && lastRowCount > 0 {
			rowSize = lastRowCount
		}
		if withDetails {
			resBuf.WriteString("Row")
			WriteStringWith0(resBuf, strconv.Itoa(*globalRowIndex), printLen-1)
			resBuf.WriteString("(")
			WriteStringWith0(resBuf, strconv.Itoa(globalByteIndex), printLen)
			resBuf.WriteString(", ")
			WriteStringWith0(resBuf, strconv.Itoa(globalByteIndex+8), printLen)
			resBuf.WriteString("): ")
			WriteRowData(resBuf, byteToStr, arr, byteIndex, rowSize, withDetails)
			resBuf.WriteString("        ")
			WriteRowDetails(resBuf, arr, byteIndex, rowSize)
			resBuf.WriteString("\n")
		} else {
			WriteRowData(resBuf, byteToStr, arr, byteIndex, rowSize, withDetails)
			resBuf.WriteString("\n")
		}
		*globalRowIndex++
	}
}

func WriteRowData(buffer *bytes.Buffer, byteToStr ByteToStr, arr []byte, off int, len2 int, withDetails bool) {
	count := 0
	for i := 0; i < len2; i++ {
		count++
		WriteStringWithSpace(buffer, byteToStr.toString(arr[off+i]), byteToStr.getWidth())
		if withDetails {
			if count&0x0F == 0 {
				buffer.WriteString(", ")
			} else {
				buffer.WriteString(" ")
			}
		} else {
			buffer.WriteString(", ")
		}
	}
}

func WriteRowDetails(buffer *bytes.Buffer, arr []byte, off int, len int) {
	count := 0
	for i := off; i < off+len; i++ {
		count++
		if arr[i] >= 32 && arr[i] <= 126 {
			buffer.WriteByte(arr[i])
		} else {
			buffer.WriteString(CharNULL)
		}
		if count&0x0F == 0 {
			buffer.WriteString(", ")
		}
	}
}

func WriteStringWithChar(buffer *bytes.Buffer, str string, expectedLen int, char rune) {
	diff := expectedLen - len(str)
	if diff > 0 {
		for i := 0; i < diff; i++ {
			buffer.WriteString(string(char))
		}
	}

	buffer.WriteString(str)
}

func WriteStringWithSpace(buffer *bytes.Buffer, str string, expectedLen int) {
	WriteStringWithChar(buffer, str, expectedLen, ' ')
}

func WriteStringWith0(buffer *bytes.Buffer, str string, expectedLen int) {
	WriteStringWithChar(buffer, str, expectedLen, '0')
}
