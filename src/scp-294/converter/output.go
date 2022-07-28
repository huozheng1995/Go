package converter

import (
	"fmt"
	"github.com/edward/scp-294/utils"
	"strconv"
	"strings"
)

const printLen = 5

func ByteArrayToString(arr []byte) string {
	rowSize := 16

	totalRow := len(arr) / rowSize
	lastRowCount := len(arr) % rowSize
	if lastRowCount > 0 {
		totalRow++
	}

	var builder strings.Builder
	for rowIndex := 0; rowIndex < totalRow; rowIndex++ {
		byteIndex := rowIndex * rowSize
		if rowIndex == totalRow-1 {
			rowSize = lastRowCount
		}
		builder.WriteString(fmt.Sprintf("Row%s(%s, %s): %s        %s\n", utils.Fill0(strconv.Itoa(rowIndex), printLen-1),
			utils.Fill0(strconv.Itoa(byteIndex), printLen), utils.Fill0(strconv.Itoa(byteIndex+8), printLen),
			utils.ByteArrayToLine(arr, byteIndex, rowSize), utils.ByteArrayToCharLine(arr, byteIndex, rowSize)))
	}
	return builder.String()
}

func Int8ArrayToString(arr []int8) string {
	rowSize := 16

	totalRow := len(arr) / rowSize
	lastRowCount := len(arr) % rowSize
	if lastRowCount > 0 {
		totalRow++
	}

	var builder strings.Builder
	for rowIndex := 0; rowIndex < totalRow; rowIndex++ {
		byteIndex := rowIndex * rowSize
		if rowIndex == totalRow-1 {
			rowSize = lastRowCount
		}
		builder.WriteString(fmt.Sprintf("Row%s(%s, %s): %s        %s\n", utils.Fill0(strconv.Itoa(rowIndex), printLen-1),
			utils.Fill0(strconv.Itoa(byteIndex), printLen), utils.Fill0(strconv.Itoa(byteIndex+8), printLen),
			utils.Int8ArrayToLine(arr[byteIndex:byteIndex+rowSize]), utils.Int8ArrayToCharLine(arr[byteIndex:byteIndex+rowSize])))
	}
	return builder.String()
}

func HexByteArrayToString(arr []string) string {
	rowSize := 16

	totalRow := len(arr) / rowSize
	lastRowCount := len(arr) % rowSize
	if lastRowCount > 0 {
		totalRow++
	}

	var builder strings.Builder
	for rowIndex := 0; rowIndex < totalRow; rowIndex++ {
		byteIndex := rowIndex * rowSize
		if rowIndex == totalRow-1 {
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

func FileByteArrayToString(arr []byte, isLastPacket bool, preArr []byte, preArrLen int) (str string, nextArr []byte, nextArrLen int) {
	rowSize := 16
	totalLen := preArrLen + len(arr)
	totalRow := totalLen / rowSize
	lastRowCount := totalLen % rowSize
	if lastRowCount > 0 {
		totalRow++
	}

	var builder strings.Builder
	for rowIndex := 0; rowIndex < totalRow; rowIndex++ {
		byteIndex := rowIndex * rowSize
		if rowIndex == totalRow-1 {
			if !isLastPacket && lastRowCount > 0 {
				break
			}
			rowSize = lastRowCount
		}
		if rowIndex == 0 {
			builder.WriteString(fmt.Sprintf("Row%s(%s, %s): %s        %s\n", utils.Fill0(strconv.Itoa(rowIndex), printLen-1),
				utils.Fill0(strconv.Itoa(byteIndex), printLen), utils.Fill0(strconv.Itoa(byteIndex+8), printLen),
				utils.TwoByteArraysToLine(preArr, 0, preArrLen, arr, byteIndex, rowSize),
				utils.TwoBytesArrayToCharLine(preArr, 0, preArrLen, arr, byteIndex, rowSize)))
		} else {
			builder.WriteString(fmt.Sprintf("Row%s(%s, %s): %s        %s\n", utils.Fill0(strconv.Itoa(rowIndex), printLen-1),
				utils.Fill0(strconv.Itoa(byteIndex), printLen), utils.Fill0(strconv.Itoa(byteIndex+8), printLen),
				utils.ByteArrayToLine(arr, byteIndex, rowSize),
				utils.ByteArrayToCharLine(arr, byteIndex, rowSize)))
		}
	}

	if !isLastPacket && lastRowCount > 0 {
		for i, j := totalRow*rowSize, 0; i < len(arr); i, j = i+1, j+1 {
			preArr[j] = arr[i]
		}
		nextArr = preArr
		nextArrLen = lastRowCount
	}

	return builder.String(), nextArr, nextArrLen
}
