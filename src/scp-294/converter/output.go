package converter

import (
	"bytes"
	"fmt"
	"github.com/edward/scp-294/utils"
	"strconv"
	"strings"
)

const printLen = 5
const GlobalRowSize = 16

func ByteArrayToString(arr []byte) string {
	rowSize := GlobalRowSize

	totalRow := len(arr) / rowSize
	lastRowCount := len(arr) % rowSize
	if lastRowCount > 0 {
		totalRow++
	}

	var builder strings.Builder
	for rowIndex := 0; rowIndex < totalRow; rowIndex++ {
		byteIndex := rowIndex * rowSize
		if rowIndex == totalRow-1 && lastRowCount > 0 {
			rowSize = lastRowCount
		}
		builder.WriteString(fmt.Sprintf("Row%s(%s, %s): %s        %s\n", utils.Fill0(strconv.Itoa(rowIndex), printLen-1),
			utils.Fill0(strconv.Itoa(byteIndex), printLen), utils.Fill0(strconv.Itoa(byteIndex+8), printLen),
			utils.BytesToByteLine(arr, byteIndex, rowSize), utils.ByteArrayToCharLine(arr, byteIndex, rowSize)))
	}
	return builder.String()
}

func Int8ArrayToString(arr []int8) string {
	rowSize := GlobalRowSize

	totalRow := len(arr) / rowSize
	lastRowCount := len(arr) % rowSize
	if lastRowCount > 0 {
		totalRow++
	}

	var builder strings.Builder
	for rowIndex := 0; rowIndex < totalRow; rowIndex++ {
		byteIndex := rowIndex * rowSize
		if rowIndex == totalRow-1 && lastRowCount > 0 {
			rowSize = lastRowCount
		}
		builder.WriteString(fmt.Sprintf("Row%s(%s, %s): %s        %s\n", utils.Fill0(strconv.Itoa(rowIndex), printLen-1),
			utils.Fill0(strconv.Itoa(byteIndex), printLen), utils.Fill0(strconv.Itoa(byteIndex+8), printLen),
			utils.Int8ArrayToLine(arr[byteIndex:byteIndex+rowSize]), utils.Int8ArrayToCharLine(arr[byteIndex:byteIndex+rowSize])))
	}
	return builder.String()
}

func HexByteArrayToString(arr []string) string {
	rowSize := GlobalRowSize

	totalRow := len(arr) / rowSize
	lastRowCount := len(arr) % rowSize
	if lastRowCount > 0 {
		totalRow++
	}

	var builder strings.Builder
	for rowIndex := 0; rowIndex < totalRow; rowIndex++ {
		byteIndex := rowIndex * rowSize
		if rowIndex == totalRow-1 && lastRowCount > 0 {
			rowSize = lastRowCount
		}
		builder.WriteString(fmt.Sprintf("Row%s(%s, %s): %s        %s\n", utils.Fill0(strconv.Itoa(rowIndex), printLen-1),
			utils.Fill0(strconv.Itoa(byteIndex), printLen), utils.Fill0(strconv.Itoa(byteIndex+8), printLen),
			utils.HexByteArrayToLine(arr[byteIndex:byteIndex+rowSize]), utils.HexByteArrayToCharLine(arr[byteIndex:byteIndex+rowSize])))
	}
	return builder.String()
}

func HexArrayToString(arr []string) string {
	var builder strings.Builder
	for _, val := range arr {
		builder.WriteString(val)
		builder.WriteString(", ")
	}
	return builder.String()
}

func BinArrayToString(arr []string) string {
	return HexArrayToString(arr)
}

func DecArrayToString(arr []int64) string {
	var builder strings.Builder
	for _, val := range arr {
		builder.WriteString(strconv.FormatInt(val, 10))
		builder.WriteString(", ")
	}
	return builder.String()
}

func StreamBytesToStringBytes(inputArr []byte, globalRowIndex *int, funcBytesToLine utils.BytesToLine) (outputArr []byte) {
	rowSize := GlobalRowSize

	totalLen := len(inputArr)
	totalRow := totalLen / rowSize
	lastRowCount := totalLen % rowSize
	if lastRowCount > 0 {
		totalRow++
	}

	buffer := new(bytes.Buffer)
	buffer.Grow(len(inputArr) * 10)
	inputArrOff := 0
	for rowIndex := 0; rowIndex < totalRow; rowIndex++ {
		globalByteIndex := *globalRowIndex * GlobalRowSize
		if rowIndex == totalRow-1 && lastRowCount > 0 {
			rowSize = lastRowCount
		}
		buffer.WriteString(fmt.Sprintf("Row%s(%s, %s): %s        %s\n", utils.Fill0(strconv.Itoa(*globalRowIndex), printLen-1),
			utils.Fill0(strconv.Itoa(globalByteIndex), printLen), utils.Fill0(strconv.Itoa(globalByteIndex+8), printLen),
			funcBytesToLine(inputArr, inputArrOff, rowSize), utils.ByteArrayToCharLine(inputArr, inputArrOff, rowSize)))
		*globalRowIndex++
		inputArrOff += rowSize
	}
	return buffer.Bytes()
}
