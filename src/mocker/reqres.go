package main

import (
	"io/ioutil"
	"strconv"
	"strings"
)

type ReqDataResData struct {
	ReqData *[]byte
	ResData *[]byte
}

type ResLenResData struct {
	ResLen  int
	ResData *[]byte
}

func AddReqDataResData(m *Mocker, reqData []byte, resData []byte) {
	m.PreSendSet.Add(&ReqDataResData{
		ReqData: &reqData,
		ResData: &resData,
	})
}

func AddResLenResData(m *Mocker, resLen int, resData []byte) {
	m.PostSendSet.Add(&ResLenResData{
		ResLen:  resLen,
		ResData: &resData,
	})
}

func HexFileToBytes(fileUri string) []byte {
	fileBytes, err := ioutil.ReadFile(fileUri)
	if err != nil {
		Log("Failed to read file: " + fileUri + ", error: " + err.Error())
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
	result = append(result, byte(10))

	return result
}
