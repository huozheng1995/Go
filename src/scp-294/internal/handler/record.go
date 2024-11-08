package handler

import (
	"encoding/json"
	"github.com/edward/scp-294/internal/dbaccess"
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
		record := dbaccess.Record{}
		err := dec.Decode(&record)
		if err != nil {
			respondError(w, "Failed to decode data, error: "+err.Error())
			return
		}

		err = record.Insert()
		if err != nil {
			respondError(w, "Failed to insert record, error: "+err.Error())
			return
		} else {
			reloadHeader(w)
		}
	default:
		respondError(w, "Failed to save record")
	}
}

func loadRecord(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		query := r.URL.Query()
		id, ok := query["RecordId"]
		if !ok {
			respondError(w, "Failed to get parameter 'RecordId'")
			return
		}
		record, err := dbaccess.GetRecord(id[0])
		if err != nil {
			respondError(w, "Failed to get record, error: "+err.Error())
			return
		}

		enc := json.NewEncoder(w)
		resData := resData{
			Success: true,
			Message: "Record was loaded!",
			Data:    record,
		}
		err = enc.Encode(resData)
		if err != nil {
			respondError(w, "Failed to encode data, error: "+err.Error())
			return
		}
	default:
		respondError(w, "Failed to load record")
	}
}

func deleteRecord(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		query := r.URL.Query()
		id, ok := query["RecordId"]
		if !ok {
			respondError(w, "Failed to get parameter 'RecordId'")
			return
		}
		err := dbaccess.DeleteRecord(id[0])
		if err != nil {
			respondError(w, "Failed to delete record, error: "+err.Error())
			return
		} else {
			reloadHeader(w)
		}
	default:
		respondError(w, "Failed to delete record")
	}
}
