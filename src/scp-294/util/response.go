package util

import (
	"bytes"
	"myutil"
	"strconv"
)

// num array to response

type NumsToResp[T any] interface {
	GetBytes() int
	ToResp(arr []T) []byte
}

type Int64sToResp struct {
	NumToStr myutil.Int64ToStr
}

func (toResp *Int64sToResp) GetBytes() int {
	return 8
}

func (toResp *Int64sToResp) ToResp(arr []int64) []byte {
	resBuf := new(bytes.Buffer)
	for _, val := range arr {
		resBuf.WriteString(toResp.NumToStr.ToString(val))
		resBuf.WriteString(", ")
	}

	return resBuf.Bytes()
}

type BytesToResp struct {
	NumToStr       myutil.ByteToStr
	WithDetails    bool
	GlobalRowIndex int
}

func (toResp *BytesToResp) GetBytes() int {
	return 1
}

func (toResp *BytesToResp) ToResp(arr []byte) []byte {
	resBuf := new(bytes.Buffer)
	//To RawBytes
	if _, ok := toResp.NumToStr.(myutil.ByteToRawBytes); ok {
		resBuf.Write(arr)
		return resBuf.Bytes()
	}

	//To String
	printLen := 5
	rowSize := 16
	totalLen := len(arr)
	totalRow := totalLen / rowSize
	lastRowCount := totalLen % rowSize
	if lastRowCount > 0 {
		totalRow++
	}

	for rowIndex := 0; rowIndex < totalRow; rowIndex++ {
		globalByteIndex := toResp.GlobalRowIndex * rowSize
		byteIndex := rowIndex * rowSize
		if rowIndex == totalRow-1 && lastRowCount > 0 {
			rowSize = lastRowCount
		}
		if toResp.WithDetails {
			resBuf.WriteString("Row")
			writeStringWith0(resBuf, strconv.Itoa(toResp.GlobalRowIndex), printLen-1)
			resBuf.WriteString("(")
			writeStringWith0(resBuf, strconv.Itoa(globalByteIndex), printLen)
			resBuf.WriteString(", ")
			writeStringWith0(resBuf, strconv.Itoa(globalByteIndex+8), printLen)
			resBuf.WriteString("): ")
			writeRowData(resBuf, toResp.NumToStr, arr, byteIndex, rowSize, toResp.WithDetails)
			resBuf.WriteString("        ")
			writeRowDetails(resBuf, arr, byteIndex, rowSize)
			resBuf.WriteString("\n")
		} else {
			writeRowData(resBuf, toResp.NumToStr, arr, byteIndex, rowSize, toResp.WithDetails)
			resBuf.WriteString("\n")
		}
		toResp.GlobalRowIndex++
	}

	return resBuf.Bytes()
}

func writeRowData(buffer *bytes.Buffer, byteToStr myutil.ByteToStr, arr []byte, off int, len2 int, withDetails bool) {
	count := 0
	for i := 0; i < len2; i++ {
		count++
		writeStringWithSpace(buffer, byteToStr.ToString(arr[off+i]), byteToStr.GetWidth())
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

func writeRowDetails(buffer *bytes.Buffer, arr []byte, off int, len int) {
	count := 0
	for i := off; i < off+len; i++ {
		count++
		if arr[i] >= 32 && arr[i] <= 126 {
			buffer.WriteByte(arr[i])
		} else {
			buffer.WriteString(myutil.CharNULL)
		}
		if count&0x0F == 0 {
			buffer.WriteString(", ")
		}
	}
}

func writeStringWithChar(buffer *bytes.Buffer, str string, expectedLen int, char rune) {
	diff := expectedLen - len(str)
	if diff > 0 {
		for i := 0; i < diff; i++ {
			buffer.WriteString(string(char))
		}
	}

	buffer.WriteString(str)
}

func writeStringWithSpace(buffer *bytes.Buffer, str string, expectedLen int) {
	writeStringWithChar(buffer, str, expectedLen, ' ')
}

func writeStringWith0(buffer *bytes.Buffer, str string, expectedLen int) {
	writeStringWithChar(buffer, str, expectedLen, '0')
}
