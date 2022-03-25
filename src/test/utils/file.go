package utils

import (
	"fmt"
	"io"
	"log"
	"os"
)

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

func PrintFileBytes2(fileName string, beginIndex int, len int) {
	const BufferSize = 32
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	defer file.Close()

	buffer := make([]byte, BufferSize)

	var rowIndex = 1
	for {
		count, err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				log.Fatal(err.Error())
			}
			break
		}

		if beginIndex < 1 {
			beginIndex = 1
		}

		print := false
		if len < 0 {
			if rowIndex >= beginIndex {
				print = true
			}
		} else if rowIndex >= beginIndex && rowIndex < beginIndex+len {
			print = true
		}

		if print {
			byteIndex := (rowIndex - 1) * BufferSize
			if count < BufferSize {
				fmt.Printf("row%d(%d, %d, %d, %d): %s\n", rowIndex, byteIndex, byteIndex+8, byteIndex+16, byteIndex+24, GetBytesData(buffer[:count]))
			} else {
				fmt.Printf("row%d(%d, %d, %d, %d): %s\n", rowIndex, byteIndex, byteIndex+8, byteIndex+16, byteIndex+24, GetBytesData(buffer))
			}
		}
		rowIndex++
	}
}

func PrintFileBytes(fileName string) {
	PrintFileBytes2(fileName, 1, -1)
}

func PrintFileBytesHex2(fileName string, beginIndex int, len int) {
	const BufferSize = 32
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	defer file.Close()

	buffer := make([]byte, BufferSize)

	var rowIndex = 1
	for {
		count, err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				log.Fatal(err.Error())
			}
			break
		}

		if beginIndex < 1 {
			beginIndex = 1
		}

		print := false
		if len < 0 {
			if rowIndex >= beginIndex {
				print = true
			}
		} else if rowIndex >= beginIndex && rowIndex < beginIndex+len {
			print = true
		}

		if print {
			byteIndex := (rowIndex - 1) * BufferSize
			if count < BufferSize {
				fmt.Printf("row%d(%d, %d, %d, %d): %s\n", rowIndex, byteIndex, byteIndex+8, byteIndex+16, byteIndex+24, GetBytesDataHex(buffer[:count]))
			} else {
				fmt.Printf("row%d(%d, %d, %d, %d): %s\n", rowIndex, byteIndex, byteIndex+8, byteIndex+16, byteIndex+24, GetBytesDataHex(buffer))
			}
		}
		rowIndex++
	}
}

func PrintFileBytesHex(fileName string) {
	PrintFileBytesHex2(fileName, 1, -1)
}
