package controller

import (
	"encoding/json"
	"github.com/edward/scp-294/common"
	"github.com/edward/scp-294/converter"
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
		"HexToDec",
		"DecToHex",
		"BinToDec",
		"DecToBin",
		"HexByteToDecByte",
		"DecByteToHexByte",
		"HexByteToInt8",
		"Int8ToHexByte",
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
		var outputData string
		switch convertReq.ConvertType {
		case "HexToDec":
			outputData = converter.DecArrayToString(converter.HexArrayToDecArray(convertReq.InputData))
		case "DecToHex":
			outputData = converter.HexArrayToString(converter.DecArrayToHexArray(convertReq.InputData))
		case "BinToDec":
			outputData = converter.DecArrayToString(converter.BinArrayToDecArray(convertReq.InputData))
		case "DecToBin":
			outputData = converter.BinArrayToString(converter.DecArrayToBinArray(convertReq.InputData))
		case "HexByteToDecByte":
			outputData = converter.ByteArrayToString(converter.HexByteArrayToDecByteArray(convertReq.InputData))
		case "DecByteToHexByte":
			outputData = converter.HexByteArrayToString(converter.DecByteArrayToHexByteArray(convertReq.InputData))
		case "HexByteToInt8":
			outputData = converter.Int8ArrayToString(converter.HexByteArrayToInt8Array(convertReq.InputData))
		case "Int8ToHexByte":
			outputData = converter.HexByteArrayToString(converter.Int8ArrayToHexByteArray(convertReq.InputData))
		default:
			outputData = convertReq.InputData
		}

		enc := json.NewEncoder(w)
		convertRes := common.ConvertRes{
			OutputData: outputData,
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
