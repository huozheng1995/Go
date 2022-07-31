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

func ByteArrayToRows(arr []byte) string {
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
			utils.BytesToByteRow(arr, byteIndex, rowSize), utils.ByteArrayToCharRow(arr, byteIndex, rowSize)))
	}
	return builder.String()
}

func Int8ArrayToRows(arr []int8) string {
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
			utils.Int8ArrayToRow(arr, byteIndex, rowSize), utils.Int8ArrayToCharRow(arr, byteIndex, rowSize)))
	}
	return builder.String()
}

func HexByteArrayToRows(arr []string) string {
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
			utils.HexByteArrayToRow(arr, byteIndex, rowSize), utils.HexByteArrayToCharRow(arr, byteIndex, rowSize)))
	}
	return builder.String()
}

func StreamBytesToRowsBytes(arr []byte, globalRowIndex *int, funcBytesToRow utils.BytesToRow) []byte {
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
		globalByteIndex := *globalRowIndex * GlobalRowSize
		byteIndex := rowIndex * rowSize
		if rowIndex == totalRow-1 && lastRowCount > 0 {
			rowSize = lastRowCount
		}
		buffer.WriteString(fmt.Sprintf("Row%s(%s, %s): %s        %s\n", utils.Fill0(strconv.Itoa(*globalRowIndex), printLen-1),
			utils.Fill0(strconv.Itoa(globalByteIndex), printLen), utils.Fill0(strconv.Itoa(globalByteIndex+8), printLen),
			funcBytesToRow(arr, byteIndex, rowSize), utils.ByteArrayToCharRow(arr, byteIndex, rowSize)))
		*globalRowIndex++
	}
	return buffer.Bytes()
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
