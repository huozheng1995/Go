package controller

import (
	"github.com/edward/scp-294/common"
	"github.com/edward/scp-294/logger"
	"github.com/edward/scp-294/model"
	"html/template"
	"net/http"
)

func RegisterRoutes() {
	http.HandleFunc("/", loadMainPage)
	registerConverterRoutes()
	registerGroupRoutes()
	registerRecordRoutes()
}

func loadMainPage(w http.ResponseWriter, r *http.Request) {
	t := template.New("layout")
	t, err := t.ParseFiles("./templates/layout.html", "./templates/header.html")
	if err != nil {
		logger.Log(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	groups, err := model.ListGroups()
	if err != nil {
		logger.Log(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	records, err := model.ListRecords()
	if err != nil {
		logger.Log(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data := struct {
		ProjName      string
		InputTypeMap  map[common.NumType]string
		OutputTypeMap map[common.NumType]string
		TypeDescMap   map[common.NumType]string
		Groups        []model.Group
		Records       []model.Record
	}{common.ProjName, common.InputTypeMap, common.OutputTypeMap, common.TypeDescMap, groups, records}
	t.ExecuteTemplate(w, "layout", data)
}

func reloadHeader(w http.ResponseWriter) {
	t := template.New("reloadHeader")
	t, err := t.ParseFiles("./templates/header.html")
	if err != nil {
		common.ResponseError(w, "Failed to parse files, error: "+err.Error())
		return
	}
	groups, err := model.ListGroups()
	if err != nil {
		common.ResponseError(w, "Failed to list groups, error: "+err.Error())
		return
	}
	records, err := model.ListRecords()
	if err != nil {
		common.ResponseError(w, "Failed to list records, error: "+err.Error())
		return
	}
	data := struct {
		Groups  []model.Group
		Records []model.Record
	}{groups, records}
	t.ExecuteTemplate(w, "header", data)
}
