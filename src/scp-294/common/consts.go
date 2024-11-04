package common

import (
	"encoding/json"
	"github.com/edward/scp-294/logger"
	"net/http"
)

type ResData struct {
	Success bool        `json:"Success"`
	Message string      `json:"Message"`
	Data    interface{} `json:"Data"`
}

func RespondError(w http.ResponseWriter, message string) {
	logger.Logger.Log("Main", message)
	enc := json.NewEncoder(w)
	resData := ResData{
		Success: false,
		Message: message,
		Data:    nil,
	}
	enc.Encode(resData)
}
