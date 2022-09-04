package utils

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

const printLen = 5
const GlobalRowSize = 16

func ByteArrayToRows(arr []byte, format bool) string {
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
		if format {
			builder.WriteString(fmt.Sprintf("Row%s(%s, %s): %s        %s\n", Fill0(strconv.Itoa(rowIndex), printLen-1),
				Fill0(strconv.Itoa(byteIndex), printLen), Fill0(strconv.Itoa(byteIndex+8), printLen),
				ByteArrayToByteRow(arr, byteIndex, rowSize, format), ByteArrayToCharRow(arr, byteIndex, rowSize)))
		} else {
			builder.WriteString(ByteArrayToByteRow(arr, byteIndex, rowSize, format))
			builder.WriteString("\n")
		}
	}
	return builder.String()
}

func Int8ArrayToRows(arr []int8, format bool) string {
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
		if format {
			builder.WriteString(fmt.Sprintf("Row%s(%s, %s): %s        %s\n", Fill0(strconv.Itoa(rowIndex), printLen-1),
				Fill0(strconv.Itoa(byteIndex), printLen), Fill0(strconv.Itoa(byteIndex+8), printLen),
				Int8ArrayToRow(arr, byteIndex, rowSize, format), Int8ArrayToCharRow(arr, byteIndex, rowSize)))
		} else {
			builder.WriteString(Int8ArrayToRow(arr, byteIndex, rowSize, format))
			builder.WriteString("\n")
		}
	}
	return builder.String()
}

func HexByteArrayToRows(arr []string, format bool) string {
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
		if format {
			builder.WriteString(fmt.Sprintf("Row%s(%s, %s): %s        %s\n", Fill0(strconv.Itoa(rowIndex), printLen-1),
				Fill0(strconv.Itoa(byteIndex), printLen), Fill0(strconv.Itoa(byteIndex+8), printLen),
				HexByteArrayToRow(arr, byteIndex, rowSize, format), HexByteArrayToCharRow(arr, byteIndex, rowSize)))
		} else {
			builder.WriteString(HexByteArrayToRow(arr, byteIndex, rowSize, format))
			builder.WriteString("\n")
		}
	}
	return builder.String()
}

func StreamBytesToRowsBytes(arr []byte, globalRowIndex *int, funcBytesToRow ByteArrayToRow, format bool) []byte {
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
		if format {
			buffer.WriteString(fmt.Sprintf("Row%s(%s, %s): %s        %s\n", Fill0(strconv.Itoa(*globalRowIndex), printLen-1),
				Fill0(strconv.Itoa(globalByteIndex), printLen), Fill0(strconv.Itoa(globalByteIndex+8), printLen),
				funcBytesToRow(arr, byteIndex, rowSize, format), ByteArrayToCharRow(arr, byteIndex, rowSize)))
		} else {
			buffer.WriteString(funcBytesToRow(arr, byteIndex, rowSize, format))
			buffer.WriteString("\n")
		}
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

func ByteArrayToString(arr []byte) string {
	var builder strings.Builder
	for _, val := range arr {
		builder.WriteString(strconv.Itoa(int(val)))
		builder.WriteString(", ")
	}
	return builder.String()
}

func Int8ArrayToString(arr []int8) string {
	var builder strings.Builder
	for _, val := range arr {
		builder.WriteString(strconv.Itoa(int(val)))
		builder.WriteString(", ")
	}
	return builder.String()
}
