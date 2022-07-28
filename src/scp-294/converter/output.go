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
			utils.ByteArrayToLine(arr, byteIndex, rowSize), utils.ByteArrayToCharLine(arr, byteIndex, rowSize)))
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

func StreamByteArrayToStringByteArray(inputArr []byte, isLastPacket bool, preArr []byte, preArrLen int, globalRowIndex *int) (outputArr []byte, nextArr []byte, nextArrLen int) {
	buffer := new(bytes.Buffer)
	buffer.Grow(4096)
	arrOff := 0
	//Handle the previous bytes
	if preArrLen > 0 {
		if preArrLen+len(inputArr) > GlobalRowSize {
			for i, j := preArrLen, 0; i < GlobalRowSize; i, j = i+1, j+1 {
				preArr[i] = inputArr[j]
				arrOff++
			}
			appendStreamByteArrayToBuffer(preArr, 0, GlobalRowSize, *globalRowIndex, buffer)
			*globalRowIndex++
		} else {
			rowSize := preArrLen
			for i, j := preArrLen, 0; i < preArrLen+len(inputArr); i, j = i+1, j+1 {
				preArr[i] = inputArr[j]
				rowSize++
				arrOff++
			}
			appendStreamByteArrayToBuffer(preArr, 0, rowSize, *globalRowIndex, buffer)
			*globalRowIndex++
			outputArr = buffer.Bytes()
			return
		}
	}
	nextArr = preArr
	nextArrLen = 0
	//Handle input bytes
	rowSize := GlobalRowSize
	totalLen := len(inputArr) - arrOff
	totalRow := totalLen / rowSize
	lastRowCount := totalLen % rowSize
	if lastRowCount > 0 {
		totalRow++
	}

	for rowIndex := 0; rowIndex < totalRow; rowIndex++ {
		if rowIndex == totalRow-1 && lastRowCount > 0 {
			if !isLastPacket {
				for i, j := arrOff, 0; i < len(inputArr); i, j = i+1, j+1 {
					preArr[j] = inputArr[i]
				}
				nextArr = preArr
				nextArrLen = lastRowCount
				outputArr = buffer.Bytes()
				return
			} else {
				rowSize = lastRowCount
			}
		}
		appendStreamByteArrayToBuffer(inputArr, arrOff, rowSize, *globalRowIndex, buffer)
		*globalRowIndex++
		arrOff += rowSize
	}
	outputArr = buffer.Bytes()
	return
}

func appendStreamByteArrayToBuffer(arr []byte, off int, len int, rowIndex int, buffer *bytes.Buffer) {
	byteIndex := rowIndex * GlobalRowSize
	buffer.WriteString(fmt.Sprintf("Row%s(%s, %s): %s        %s\n", utils.Fill0(strconv.Itoa(rowIndex), printLen-1),
		utils.Fill0(strconv.Itoa(byteIndex), printLen), utils.Fill0(strconv.Itoa(byteIndex+8), printLen),
		utils.ByteArrayToLine(arr, off, len), utils.ByteArrayToCharLine(arr, off, len)))
}
