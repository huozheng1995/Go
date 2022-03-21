package utils

import (
	"fmt"
	"io"
	"os"
	"strconv"
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

func PrintFileBytes(fileName string) {
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
			fmt.Println("row" + strconv.Itoa(rowIndex) + ": " + GetBytesData(buffer[:count]))
		} else {
			fmt.Println("row" + strconv.Itoa(rowIndex) + ": " + GetBytesData(buffer))
		}
		rowIndex++
	}
}

func PrintFileBytesHex(fileName string) {
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
			fmt.Println("row" + strconv.Itoa(rowIndex) + ": " + GetBytesDataHex(buffer[:count]))
		} else {
			fmt.Println("row" + strconv.Itoa(rowIndex) + ": " + GetBytesDataHex(buffer))
		}
		rowIndex++
	}
}
