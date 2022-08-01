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
		"File",
	}
	outputTypes := []string{
		"Hex",
		"Dec",
		"Bin",
		"HexByteArray",
		"ByteArray",
		"Int8Array",
	}
	inputCases := []string{
		"ABAB5648, 12EF51,56 75,,CCCCCC  2CDD, DC11248, 05, 12, FE, FF",
		"2880132680, 1240913,86 117,,13421772  11485, 230756936, 5, 18, 254, 255",
		"10101011101010110101011001001000, 100101110111101010001,1010110 1110101,,110011001100110011001100  10110011011101, 1101110000010001001001001000, 101, 10010, 11111110, 11111111",
		"AB, EF,56 75,,CC  2C, DC, BB, FE, FF",
		"171, 239,86 117,,204  44, 220, 187, 254, 255",
		"-85, -17,86 117,,-52  44, -36, -69, -2, -1",
		"Select a file to parse",
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
		InputCases  []string
		Groups      []model.Group
		Records     []model.Record
	}{common.ProjName, inputTypes, outputTypes, inputCases, groups, records}
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
