package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

func ReadFile(fileName string) {
	const BufferSize = 100
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
			fmt.Println("row" + strconv.Itoa(rowIndex) + ": " + string(buffer[:count]))
		} else {
			fmt.Println("row" + strconv.Itoa(rowIndex) + ": " + string(buffer))
		}
		rowIndex++
	}
}

func ReadFileBytes(fileName string) {
	const BufferSize = 100
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
			fmt.Println("row" + strconv.Itoa(rowIndex) + ": " + ReadBytesData(buffer[:count]))
		} else {
			fmt.Println("row" + strconv.Itoa(rowIndex) + ": " + ReadBytesData(buffer))
		}
		rowIndex++
	}
}
