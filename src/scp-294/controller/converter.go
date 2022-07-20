package controller

import (
	"encoding/json"
	"github.com/edward/scp-294/common"
	"github.com/edward/scp-294/logger"
	"github.com/edward/scp-294/model"
	"html/template"
	"net/http"
)

func registerRoutes() {
	http.HandleFunc("/", loadMainPage)
	http.HandleFunc("/loadMainPage", loadMainPage)
	http.HandleFunc("/convert", convert)
}

func loadMainPage(w http.ResponseWriter, r *http.Request) {
	t := template.New("layout")
	t, err := t.ParseFiles("./templates/layout.html", "./templates/header.html")
	if err != nil {
		logger.Log(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	convertTypes := []string{
		"DecToHex",
		"HexToDec",
		"DecToBinary",
		"BinaryToDec",
		"HexToBinary",
		"BinaryToHex",
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
		ProjName     string
		ConvertTypes []string
		Groups       []model.Group
		Records      []model.Record
	}{common.ProjName, convertTypes, groups, records}
	t.ExecuteTemplate(w, "layout", data)
}

func convert(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		dec := json.NewDecoder(r.Body)
		convertReq := common.ConvertReq{}
		err := dec.Decode(&convertReq)
		if err != nil {
			logger.Log(err.Error())
			common.ResponseError(w, "Failed to decode data, error: "+err.Error())
			return
		}

		enc := json.NewEncoder(w)
		convertRes := common.ConvertRes{
			OutputData: convertReq.InputData,
		}
		resData := common.ResData{
			Success: true,
			Message: "Data was converted!",
			Data:    convertRes,
		}
		err = enc.Encode(resData)
		if err != nil {
			logger.Log(err.Error())
			common.ResponseError(w, "Failed to encode data, error: "+err.Error())
			return
		}
	default:
		common.ResponseError(w, "Failed to convert data")
	}
}
