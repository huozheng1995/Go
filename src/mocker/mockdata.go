package main

import (
	"github.com/edward/mocker/logger"
	"os"
	"strconv"
	"strings"
)

type ReqDataResFiles struct {
	ReqData  *[]byte
	FileUris []string
}

type ResLenResFiles struct {
	ResLen   int
	FileUris []string
}

func HexFileToBytes(fileUri string) []byte {
	fileBytes, err := os.ReadFile(fileUri)
	if err != nil {
		logger.Log("Failed to read file: " + fileUri + ", error: " + err.Error())
		panic(err)
	}

	result := make([]byte, 0, len(fileBytes)>>1)
	var val byte
	var builder strings.Builder
	for i := 0; i < len(fileBytes)+1; i++ {
		if i == len(fileBytes) {
			val = 0
		} else {
			val = fileBytes[i]
		}
		if (val >= '0' && val <= '9') || (val >= 'a' && val <= 'f') || (val >= 'A' && val <= 'F') || val == '-' {
			builder.WriteByte(val)
		} else {
			if builder.Len() > 0 {
				intVal, _ := strconv.ParseInt(builder.String(), 16, 64)
				builder.Reset()
				result = append(result, byte(intVal))
			}
		}
	}

	return result
}
