package controller

import (
	"encoding/json"
	"github.com/edward/scp-294/common"
	"github.com/edward/scp-294/logger"
	"github.com/edward/scp-294/model"
	"html/template"
	"net/http"
	"net/http/httputil"
)

func registerRecordRoutes() {
	http.HandleFunc("/saveRecord", saveRecord)
	http.HandleFunc("/loadRecord", loadRecord)
	http.HandleFunc("/deleteRecord", deleteRecord)
}

func saveRecord(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		dec := json.NewDecoder(r.Body)
		record := model.Record{}
		err := dec.Decode(&record)
		if err != nil {
			requestDump, _ := httputil.DumpRequest(r, true)
			logger.Log(string(requestDump))
			logger.Log(err.Error())
			common.ResponseError(w, "Failed to decode data, error: "+err.Error())
			return
		}

		err = record.Insert()
		if err != nil {
			logger.Log(err.Error())
			common.ResponseError(w, "Failed to insert record, error: "+err.Error())
			return
		} else {
			reloadHeader(w)
		}
	default:
		common.ResponseError(w, "Failed to save record")
	}
}

func reloadHeader(w http.ResponseWriter) {
	t := template.New("reloadHeader")
	t, err := t.ParseFiles("./templates/header.html")
	if err != nil {
		logger.Log(err.Error())
		common.ResponseError(w, "Failed to parse files, error: "+err.Error())
		return
	}
	groups, err := model.ListGroups()
	if err != nil {
		logger.Log(err.Error())
		common.ResponseError(w, "Failed to list groups, error: "+err.Error())
		return
	}
	records, err := model.ListRecords()
	if err != nil {
		logger.Log(err.Error())
		common.ResponseError(w, "Failed to list records, error: "+err.Error())
		return
	}
	data := struct {
		Groups  []model.Group
		Records []model.Record
	}{groups, records}
	t.ExecuteTemplate(w, "header", data)
}

func loadRecord(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		query := r.URL.Query()
		id, ok := query["RecordId"]
		if !ok {
			common.ResponseError(w, "Failed to get parameter 'RecordId'")
		}
		record, err := model.GetRecord(id[0])
		if err != nil {
			logger.Log(err.Error())
			common.ResponseError(w, "Failed to get record, error: "+err.Error())
		}

		enc := json.NewEncoder(w)
		resData := common.ResData{
			Success: true,
			Message: "Record was loaded!",
			Data:    record,
		}
		err = enc.Encode(resData)
		if err != nil {
			logger.Log(err.Error())
			common.ResponseError(w, "Failed to encode data, error: "+err.Error())
			return
		}
	default:
		common.ResponseError(w, "Failed to load record")
	}
}

func deleteRecord(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		query := r.URL.Query()
		id, ok := query["RecordId"]
		if !ok {
			common.ResponseError(w, "Failed to get parameter 'RecordId'")
		}
		err := model.DeleteRecord(id[0])
		if err != nil {
			logger.Log(err.Error())
			common.ResponseError(w, "Failed to delete record, error: "+err.Error())
		} else {
			reloadHeader(w)
		}
	default:
		common.ResponseError(w, "Failed to delete record")
	}
}
