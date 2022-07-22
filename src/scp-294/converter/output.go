package converter

import (
	"fmt"
	"github.com/edward/scp-294/utils"
	"strconv"
	"strings"
)

const printLen = 5

func ByteArrayToString(arr []byte) string {
	rowCount := 16

	totalRow := len(arr) / rowCount
	lastRowCount := len(arr) % rowCount
	if lastRowCount > 0 {
		totalRow++
	}

	var builder strings.Builder
	for rowIndex := 0; rowIndex < totalRow; rowIndex++ {
		byteIndex := rowIndex * rowCount
		if rowIndex == totalRow-1 {
			rowCount = lastRowCount
		}
		builder.WriteString(fmt.Sprintf("Row%s(%s, %s): %s        %s\n", utils.Fill0(strconv.Itoa(rowIndex), printLen-1),
			utils.Fill0(strconv.Itoa(byteIndex), printLen), utils.Fill0(strconv.Itoa(byteIndex+8), printLen),
			utils.ByteArrayToLine(arr[byteIndex:byteIndex+rowCount]), utils.ByteArrayToCharLine(arr[byteIndex:byteIndex+rowCount])))
	}
	return builder.String()
}

func Int8ArrayToString(arr []int8) string {
	rowCount := 16

	totalRow := len(arr) / rowCount
	lastRowCount := len(arr) % rowCount
	if lastRowCount > 0 {
		totalRow++
	}

	var builder strings.Builder
	for rowIndex := 0; rowIndex < totalRow; rowIndex++ {
		byteIndex := rowIndex * rowCount
		if rowIndex == totalRow-1 {
			rowCount = lastRowCount
		}
		builder.WriteString(fmt.Sprintf("Row%s(%s, %s): %s        %s\n", utils.Fill0(strconv.Itoa(rowIndex), printLen-1),
			utils.Fill0(strconv.Itoa(byteIndex), printLen), utils.Fill0(strconv.Itoa(byteIndex+8), printLen),
			utils.Int8ArrayToLine(arr[byteIndex:byteIndex+rowCount]), utils.Int8ArrayToCharLine(arr[byteIndex:byteIndex+rowCount])))
	}
	return builder.String()
}

func HexByteArrayToString(arr []string) string {
	rowCount := 16

	totalRow := len(arr) / rowCount
	lastRowCount := len(arr) % rowCount
	if lastRowCount > 0 {
		totalRow++
	}

	var builder strings.Builder
	for rowIndex := 0; rowIndex < totalRow; rowIndex++ {
		byteIndex := rowIndex * rowCount
		if rowIndex == totalRow-1 {
			rowCount = lastRowCount
		}
		builder.WriteString(fmt.Sprintf("Row%s(%s, %s): %s        %s\n", utils.Fill0(strconv.Itoa(rowIndex), printLen-1),
			utils.Fill0(strconv.Itoa(byteIndex), printLen), utils.Fill0(strconv.Itoa(byteIndex+8), printLen),
			utils.HexByteArrayToLine(arr[byteIndex:byteIndex+rowCount]), utils.HexByteArrayToCharLine(arr[byteIndex:byteIndex+rowCount])))
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
	var builder strings.Builder
	for _, val := range arr {
		builder.WriteString(val)
		builder.WriteString(", ")
	}
	return builder.String()
}

func DecArrayToString(arr []int64) string {
	var builder strings.Builder
	for _, val := range arr {
		builder.WriteString(strconv.FormatInt(val, 10))
		builder.WriteString(", ")
	}
	return builder.String()
}
