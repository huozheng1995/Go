package common

import (
	"encoding/json"
	"net/http"
)

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
	enc := json.NewEncoder(w)
	resData := ResData{
		Success: false,
		Message: message,
		Data:    nil,
	}
	enc.Encode(resData)
}
