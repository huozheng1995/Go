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

//func printFileBytes(fileName string, file multipart.File, beginIndex int, len int) {
//	const BufferSize = 32
//	defer file.Close()
//
//	buffer := make([]byte, BufferSize)
//
//	if beginIndex < 1 {
//		beginIndex = 1
//	}
//	var rowIndex = 1
//	for {
//		count, err := file.Read(buffer)
//		if err != nil {
//			if err != io.EOF {
//				log.Fatal(err.Error())
//			}
//			break
//		}
//
//		switch PrintState {
//		case beforePrint:
//			if rowIndex >= beginIndex {
//				PrintState = printing
//			}
//		case printing:
//			if len > 0 && rowIndex >= beginIndex+len {
//				return
//			}
//		}
//
//		if PrintState == printing {
//			byteIndex := (rowIndex - 1) * BufferSize
//			if count < BufferSize {
//				fmt.Printf("row%s(%s, %s, %s, %s): %s\n", Fill0(strconv.Itoa(rowIndex), printLen),
//					Fill0(strconv.Itoa(byteIndex), printLen), Fill0(strconv.Itoa(byteIndex+8), printLen),
//					Fill0(strconv.Itoa(byteIndex+16), printLen), Fill0(strconv.Itoa(byteIndex+24), printLen),
//					bytesDataToNum(buffer[:count]))
//			} else {
//				fmt.Printf("row%s(%s, %s, %s, %s): %s\n", Fill0(strconv.Itoa(rowIndex), printLen),
//					Fill0(strconv.Itoa(byteIndex), printLen), Fill0(strconv.Itoa(byteIndex+8), printLen),
//					Fill0(strconv.Itoa(byteIndex+16), printLen), Fill0(strconv.Itoa(byteIndex+24), printLen),
//					bytesDataToNum(buffer))
//			}
//		}
//		rowIndex++
//	}
//}
