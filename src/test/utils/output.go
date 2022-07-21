package utils

import (
	"fmt"
	"strconv"
	"strings"
)

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
		builder.WriteString(fmt.Sprintf("row%s(%s, %s): %s\n", Fill0(strconv.Itoa(rowIndex), printLen),
			Fill0(strconv.Itoa(byteIndex), printLen), Fill0(strconv.Itoa(byteIndex+8), printLen),
			ByteArrayToLine(arr[byteIndex:byteIndex+rowCount])))
	}
	return builder.String()
}
