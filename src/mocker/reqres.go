package main

import (
	"io/ioutil"
	"strconv"
	"strings"
)

type ReqRes struct {
	Request  *[]byte
	Response *[]byte
}

func AddReqRes(m *Mocker, request []byte, response []byte) {
	m.MockSet.Add(&ReqRes{
		Request:  &request,
		Response: &response,
	})
}

func AddReqResFromFile(m *Mocker, reqFileUri string, resFileUri string) {
	// Read request file
	reqData, err := ioutil.ReadFile(reqFileUri)
	if err != nil {
		Log("Failed to read request file, error: " + err.Error())
		panic(err)
	}

	// Read response file
	resData, err := ioutil.ReadFile(resFileUri)
	if err != nil {
		Log("Failed to read response file" + err.Error())
		panic(err)
	}

	// Add request-response to mock set
	AddReqRes(m, FileBytesToByteArray(reqData), FileBytesToByteArray(resData))
}

func FileBytesToByteArray(fileBytes []byte) []byte {
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
	result = append(result, byte(10))

	return result
}
