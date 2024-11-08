package handler

import (
	"encoding/json"
	"github.com/edward/scp-294/internal"
	"net/http"
)

type resData struct {
	Success bool        `json:"Success"`
	Message string      `json:"Message"`
	Data    interface{} `json:"Data"`
}

func respondError(w http.ResponseWriter, message string) {
	internal.Logger.Log("Main", message)
	enc := json.NewEncoder(w)
	resData := resData{
		Success: false,
		Message: message,
		Data:    nil,
	}
	enc.Encode(resData)
}
