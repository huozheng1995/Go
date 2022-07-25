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

	inputTypes := []string{
		"Hex",
		"Dec",
		"Bin",
		"HexByteArray",
		"ByteArray",
		"Int8Array",
	}
	outputTypes := []string{
		"Hex",
		"Dec",
		"Bin",
		"HexByteArray",
		"ByteArray",
		"Int8Array",
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
		ProjName    string
		InputTypes  []string
		OutputTypes []string
		Groups      []model.Group
		Records     []model.Record
	}{common.ProjName, inputTypes, outputTypes, groups, records}
	t.ExecuteTemplate(w, "layout", data)
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
