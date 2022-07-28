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
		if rowIndex == totalRow-1 && lastRowCount > 0 {
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
	rowSize := 16

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

const GlobalRowSize = 16

func FileByteArrayToString(arr []byte, isLastPacket bool, preArr []byte, preArrLen int, globalRowIndex *int) (str string, nextArr []byte, nextArrLen int) {
	rowSize := GlobalRowSize
	totalLen := preArrLen + len(arr)
	totalRow := totalLen / rowSize
	lastRowCount := totalLen % rowSize
	if lastRowCount > 0 {
		totalRow++
	}

	var builder strings.Builder
	byteIndex := 0
	for rowIndex := 0; rowIndex < totalRow; rowIndex++ {
		if rowIndex == totalRow-1 && lastRowCount > 0 {
			if !isLastPacket {
				for i, j := byteIndex, 0; i < len(arr); i, j = i+1, j+1 {
					preArr[j] = arr[i]
				}
				nextArr = preArr
				nextArrLen = lastRowCount
				break
			} else {
				rowSize = lastRowCount
			}
		}
		*globalRowIndex++
		globalByteIndex := *globalRowIndex * GlobalRowSize
		if rowIndex == 0 {
			for i, j := preArrLen, 0; i < len(preArr); i, j = i+1, j+1 {
				preArr[i] = arr[j]
				byteIndex++
			}
			builder.WriteString(fmt.Sprintf("Row%s(%s, %s): %s        %s\n", utils.Fill0(strconv.Itoa(*globalRowIndex), printLen-1),
				utils.Fill0(strconv.Itoa(globalByteIndex), printLen), utils.Fill0(strconv.Itoa(globalByteIndex+8), printLen),
				utils.ByteArrayToLine(preArr, 0, rowSize),
				utils.ByteArrayToCharLine(preArr, 0, rowSize)))
		} else {
			builder.WriteString(fmt.Sprintf("Row%s(%s, %s): %s        %s\n", utils.Fill0(strconv.Itoa(*globalRowIndex), printLen-1),
				utils.Fill0(strconv.Itoa(globalByteIndex), printLen), utils.Fill0(strconv.Itoa(globalByteIndex+8), printLen),
				utils.ByteArrayToLine(arr, byteIndex, rowSize),
				utils.ByteArrayToCharLine(arr, byteIndex, rowSize)))
			byteIndex += rowSize
		}
	}
	str = builder.String()
	return
}
