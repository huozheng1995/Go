package controller

import (
	"encoding/json"
	"github.com/edward/scp-294/common"
	"github.com/edward/scp-294/logger"
	"github.com/edward/scp-294/model"
	"net/http"
	"net/http/httputil"
)

func registerRecordRoutes() {
	http.HandleFunc("/saveRecord", saveRecord)
}

func saveRecord(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		dec := json.NewDecoder(r.Body)
		record := model.Record{}
		err := dec.Decode(&record)
		if err != nil {
			logger.Log(err.Error())
			requestDump, _ := httputil.DumpRequest(r, true)
			logger.Log(string(requestDump))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = record.Insert()
		if err != nil {
			logger.Log(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		enc := json.NewEncoder(w)
		resData := common.ResData{
			Success: true,
			Message: "Record was saved!",
			Data:    nil,
		}
		err = enc.Encode(resData)
		if err != nil {
			logger.Log(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
