package handler

import (
	"encoding/json"
	"github.com/edward/scp-294/internal/dbaccess"
	"net/http"
)

func registerGroupRoutes() {
	http.HandleFunc("/addGroup", addGroup)
	http.HandleFunc("/deleteGroup", deleteGroup)
}

func addGroup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		dec := json.NewDecoder(r.Body)
		group := dbaccess.Group{}
		err := dec.Decode(&group)
		if err != nil {
			respondError(w, "Failed to decode data, error: "+err.Error())
			return
		}

		err = group.Insert()
		if err != nil {
			respondError(w, "Failed to insert group, error: "+err.Error())
			return
		} else {
			reloadHeader(w)
		}
	default:
		respondError(w, "Failed to save group")
	}
}

func deleteGroup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		query := r.URL.Query()
		id, ok := query["GroupId"]
		if !ok {
			respondError(w, "Failed to get parameter 'GroupId'")
			return
		}
		err := dbaccess.DeleteRecordsByGroupId(id[0])
		if err != nil {
			respondError(w, "Failed to delete records, error: "+err.Error())
			return
		}
		err = dbaccess.DeleteGroup(id[0])
		if err != nil {
			respondError(w, "Failed to delete group, error: "+err.Error())
			return
		} else {
			reloadHeader(w)
		}
	default:
		respondError(w, "Failed to delete group")
	}
}
