package utils

import (
	"bytes"
	"myutil"
	"strconv"
)

// int64 array to response

func Int64ArrayToResponse(arr []int64, funcInt64ToStr myutil.Int64ToStr) *bytes.Buffer {
	resBuf := new(bytes.Buffer)
	for _, val := range arr {
		resBuf.WriteString(funcInt64ToStr(val))
		resBuf.WriteString(", ")
	}

	return resBuf
}

// byte array to response

func ByteArrayToResponse(arr []byte, globalRowIndex *int, byteToStr myutil.ByteToStr, withDetails bool) *bytes.Buffer {
	resBuf := new(bytes.Buffer)
	//To RawBytes
	if _, ok := byteToStr.(myutil.ByteToRawBytes); ok {
		resBuf.Write(arr)
		return resBuf
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

	return resBuf
}

func WriteRowData(buffer *bytes.Buffer, byteToStr myutil.ByteToStr, arr []byte, off int, len2 int, withDetails bool) {
	count := 0
	for i := 0; i < len2; i++ {
		count++
		WriteStringWithSpace(buffer, byteToStr.ToString(arr[off+i]), byteToStr.GetWidth())
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
			buffer.WriteString(myutil.CharNULL)
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
