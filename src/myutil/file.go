package myutil

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

const (
	beforePrint, printing = 1, 2
	printLen              = 5
)

var PrintState = beforePrint

func PrintFile(fileName string) {
	const BufferSize = 64
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	buffer := make([]byte, BufferSize)

	var rowIndex = 1
	for {
		count, err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}

		if count < BufferSize {
			fmt.Printf("%s\n", string(buffer[:count]))
		} else {
			fmt.Printf("%s\n", string(buffer))
		}
		rowIndex++
	}
}

func PrintFileBytesDec(fileName string) {
	printFileBytes(fileName, 1, -1, ByteArrayToLine)
}

func PrintFileBytesHex(fileName string) {
	printFileBytes(fileName, 1, -1, ByteArrayToHex)
}

func PrintFileBytesDec2(fileName string, beginIndex int, len int) {
	printFileBytes(fileName, beginIndex, len, ByteArrayToLine)
}

func PrintFileBytesHex2(fileName string, beginIndex int, len int) {
	printFileBytes(fileName, beginIndex, len, ByteArrayToHex)
}

func printFileBytes(fileName string, beginIndex int, len int, bytesDataToNum BytesDataToNum) {
	const BufferSize = 32
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	defer file.Close()

	buffer := make([]byte, BufferSize)

	if beginIndex < 1 {
		beginIndex = 1
	}
	var rowIndex = 1
	for {
		count, err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				log.Fatal(err.Error())
			}
			break
		}

		switch PrintState {
		case beforePrint:
			if rowIndex >= beginIndex {
				PrintState = printing
			}
		case printing:
			if len > 0 && rowIndex >= beginIndex+len {
				return
			}
		}

		if PrintState == printing {
			byteIndex := (rowIndex - 1) * BufferSize
			if count < BufferSize {
				fmt.Printf("row%s(%s, %s, %s, %s): %s\n", Fill0(strconv.Itoa(rowIndex), printLen),
					Fill0(strconv.Itoa(byteIndex), printLen), Fill0(strconv.Itoa(byteIndex+8), printLen),
					Fill0(strconv.Itoa(byteIndex+16), printLen), Fill0(strconv.Itoa(byteIndex+24), printLen),
					bytesDataToNum(buffer[:count]))
			} else {
				fmt.Printf("row%s(%s, %s, %s, %s): %s\n", Fill0(strconv.Itoa(rowIndex), printLen),
					Fill0(strconv.Itoa(byteIndex), printLen), Fill0(strconv.Itoa(byteIndex+8), printLen),
					Fill0(strconv.Itoa(byteIndex+16), printLen), Fill0(strconv.Itoa(byteIndex+24), printLen),
					bytesDataToNum(buffer))
			}
		}
		rowIndex++
	}
}
