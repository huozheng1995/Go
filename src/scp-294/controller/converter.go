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

func registerRoutes() {
	http.HandleFunc("/", loadMainPage)
	http.HandleFunc("/convert", convert)
}

func loadMainPage(w http.ResponseWriter, r *http.Request) {
	records, err := model.ListRecords()
	if err != nil {
		logger.Log(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	t := template.New("layout")
	t, err = t.ParseFiles("./templates/layout.html", "./templates/records.html")
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
	data := struct {
		ProjName     string
		ConvertTypes []string
		Records      []model.Record
	}{common.ProjName, convertTypes, records}
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
			requestDump, _ := httputil.DumpRequest(r, true)
			logger.Log(string(requestDump))
			w.WriteHeader(http.StatusInternalServerError)
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
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
