package controller

import (
	"encoding/json"
	"github.com/edward/scp-294/common"
	"github.com/edward/scp-294/logger"
	"github.com/edward/scp-294/model"
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
		group := model.Group{}
		err := dec.Decode(&group)
		if err != nil {
			logger.Log(err.Error())
			common.ResponseError(w, "Failed to decode data, error: "+err.Error())
			return
		}

		err = group.Insert()
		if err != nil {
			logger.Log(err.Error())
			common.ResponseError(w, "Failed to insert group, error: "+err.Error())
			return
		} else {
			reloadHeader(w)
		}
	default:
		common.ResponseError(w, "Failed to save group")
	}
}

func deleteGroup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		query := r.URL.Query()
		id, ok := query["GroupId"]
		if !ok {
			common.ResponseError(w, "Failed to get parameter 'GroupId'")
		}
		err := model.DeleteRecordsByGroupId(id[0])
		if err != nil {
			logger.Log(err.Error())
			common.ResponseError(w, "Failed to delete records, error: "+err.Error())
		}
		err = model.DeleteGroup(id[0])
		if err != nil {
			logger.Log(err.Error())
			common.ResponseError(w, "Failed to delete group, error: "+err.Error())
		} else {
			reloadHeader(w)
		}
	default:
		common.ResponseError(w, "Failed to delete group")
	}
}
