package controller

import (
	"encoding/json"
	"github.com/edward/scp-294/common"
	"github.com/edward/scp-294/model"
	"net/http"
)

func registerRecordRoutes() {
	http.HandleFunc("/addRecord", addRecord)
	http.HandleFunc("/loadRecord", loadRecord)
	http.HandleFunc("/deleteRecord", deleteRecord)
}

func addRecord(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		dec := json.NewDecoder(r.Body)
		record := model.Record{}
		err := dec.Decode(&record)
		if err != nil {
			common.RespondError(w, "Failed to decode data, error: "+err.Error())
			return
		}

		err = record.Insert()
		if err != nil {
			common.RespondError(w, "Failed to insert record, error: "+err.Error())
			return
		} else {
			reloadHeader(w)
		}
	default:
		common.RespondError(w, "Failed to save record")
	}
}

func loadRecord(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		query := r.URL.Query()
		id, ok := query["RecordId"]
		if !ok {
			common.RespondError(w, "Failed to get parameter 'RecordId'")
			return
		}
		record, err := model.GetRecord(id[0])
		if err != nil {
			common.RespondError(w, "Failed to get record, error: "+err.Error())
			return
		}

		enc := json.NewEncoder(w)
		resData := common.ResData{
			Success: true,
			Message: "Record was loaded!",
			Data:    record,
		}
		err = enc.Encode(resData)
		if err != nil {
			common.RespondError(w, "Failed to encode data, error: "+err.Error())
			return
		}
	default:
		common.RespondError(w, "Failed to load record")
	}
}

func deleteRecord(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		query := r.URL.Query()
		id, ok := query["RecordId"]
		if !ok {
			common.RespondError(w, "Failed to get parameter 'RecordId'")
			return
		}
		err := model.DeleteRecord(id[0])
		if err != nil {
			common.RespondError(w, "Failed to delete record, error: "+err.Error())
			return
		} else {
			reloadHeader(w)
		}
	default:
		common.RespondError(w, "Failed to delete record")
	}
}
