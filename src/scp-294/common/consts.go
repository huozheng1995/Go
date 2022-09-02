package common

import (
	"encoding/json"
	"github.com/edward/scp-294/logger"
	"net/http"
)

type NumType int

const (
	Hex NumType = iota
	Dec
	Bin
	HexByte
	DecByte
	DecInt8
	File
)

var InputTypeMap = map[NumType]string{
	Hex:     "Hex numbers",
	Dec:     "Dec numbers",
	Bin:     "Bin numbers",
	HexByte: "Hex Byte numbers",
	DecByte: "Dec Byte numbers",
	DecInt8: "Dec Int8 numbers",
	File:    "File",
}

var OutputTypeMap = map[NumType]string{
	Hex:     "Hex numbers",
	Dec:     "Dec numbers",
	Bin:     "Bin numbers",
	HexByte: "Hex Byte numbers",
	DecByte: "Dec Byte numbers",
	DecInt8: "Dec Int8 numbers",
}

var TypeDescMap = map[NumType]string{
	Hex:     "ABAB5, 12EF1, 56, 75, CCCCC, 2CDD, DC11248, 05, 12, FE, FF, ",
	Dec:     "703157, 77553, 86, 117, 838860, 11485, 230756936, 5, 18, 254, 255, ",
	Bin:     "10101011101010110101, 10010111011110001, 1010110, 1110101, 11001100110011001100, 10110011011101, 1101110000010001001001001000, 101, 10010, 11111110, 11111111, ",
	HexByte: "AB, EF, 56, 75, CC, 2C, DC, BB, FE, FF, ",
	DecByte: "171, 239, 86, 117, 204, 44, 220, 187, 254, 255, ",
	DecInt8: "-85, -17, 86, 117, -52, 44, -36, -69, -2, -1, ",
	File:    "Select a file to parse",
}

type ResData struct {
	Success bool        `json:"Success"`
	Message string      `json:"Message"`
	Data    interface{} `json:"Data"`
}

type ConvertReq struct {
	InputType  string `json:"InputType"`
	OutputType string `json:"OutputType"`
	InputData  string `json:"InputData"`
}

func ResponseError(w http.ResponseWriter, message string) {
	logger.Log(message)
	enc := json.NewEncoder(w)
	resData := ResData{
		Success: false,
		Message: message,
		Data:    nil,
	}
	enc.Encode(resData)
}
