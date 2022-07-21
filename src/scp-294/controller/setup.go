package controller

import (
	"github.com/edward/scp-294/common"
	"github.com/edward/scp-294/logger"
	"github.com/edward/scp-294/model"
	"html/template"
	"net/http"
)

func RegisterRoutes() {
	registerRoutes()
	registerGroupRoutes()
	registerRecordRoutes()
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
