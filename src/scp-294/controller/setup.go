package controller

import (
	"github.com/edward/scp-294/common"
	"github.com/edward/scp-294/inout"
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
	t, err := t.ParseFiles("./template/layout.html", "./template/header.html")
	if err != nil {
		logger.Logger.Log("Main", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	groups, err := model.ListGroups()
	if err != nil {
		logger.Logger.Log("Main", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	records, err := model.ListRecords()
	if err != nil {
		logger.Logger.Log("Main", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := struct {
		ProjName            string
		InoutTypes          []inout.Type
		InoutFormatMappings []inout.FormatMapping
		Groups              []model.Group
		Records             []model.Record
	}{
		"SCP-294",
		inout.CreateTypes(),
		inout.CreateFormatMappings(),
		groups,
		records,
	}
	t.ExecuteTemplate(w, "layout", data)
}

func reloadHeader(w http.ResponseWriter) {
	t := template.New("reloadHeader")
	t, err := t.ParseFiles("./template/header.html")
	if err != nil {
		common.RespondError(w, "Failed to parse files, error: "+err.Error())
		return
	}
	groups, err := model.ListGroups()
	if err != nil {
		common.RespondError(w, "Failed to list groups, error: "+err.Error())
		return
	}
	records, err := model.ListRecords()
	if err != nil {
		common.RespondError(w, "Failed to list records, error: "+err.Error())
		return
	}
	data := struct {
		Groups  []model.Group
		Records []model.Record
	}{groups, records}
	t.ExecuteTemplate(w, "header", data)
}
