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
	ConvertType string `json:"ConvertType"`
	InputData   string `json:"InputData"`
}

type ConvertRes struct {
	OutputData string `json:"OutputData"`
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
