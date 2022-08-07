package controller

import (
	"encoding/json"
	"github.com/edward/scp-294/common"
	"github.com/edward/scp-294/logger"
	"github.com/edward/scp-294/utils"
	"mime/multipart"
	"net/http"
	"strconv"
)

func registerConverterRoutes() {
	http.HandleFunc("/convert", convert)
}

func convert(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var err error
		err = r.ParseMultipartForm(100)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			common.ResponseError(w, "Failed to parse form data, error: "+err.Error())
			return
		}
		form := r.MultipartForm

		var InputType, OutputType, InputData string
		for key, values := range form.Value {
			if key == "InputType" {
				InputType = values[0]
			} else if key == "OutputType" {
				OutputType = values[0]
			} else if key == "InputData" {
				InputData = values[0]
			}
		}

		//Handle File
		if InputType == "File" {
			var file multipart.File
			for key, files := range form.File {
				if key == "InputFile" {
					file, err = files[0].Open()
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						common.ResponseError(w, "Failed to open file, error: "+err.Error())
						return
					}
				}
			}
			var funcBytesToRow utils.BytesToRow
			switch OutputType {
			case "HexByteArray":
				funcBytesToRow = utils.BytesToHexRow
			case "ByteArray":
				funcBytesToRow = utils.BytesToByteRow
			case "Int8Array":
				funcBytesToRow = utils.BytesToInt8Row
			default:
				w.WriteHeader(http.StatusInternalServerError)
				common.ResponseError(w, "Cannot convert File to '"+OutputType+"'")
				return
			}
			exitChan, dataChan := utils.FileStreamToChannel(file, utils.GlobalRowSize*256)
			readStreamAndSendBody(w, dataChan, funcBytesToRow)
			close(exitChan)
			return
		}

		var outputData string
		var strings = utils.SplitInputString(InputData)
		switch InputType {
		case "Hex":
			switch OutputType {
			case "Hex":
				outputData = utils.HexArrayToString(strings)
			case "Dec":
				outputData = utils.DecArrayToString(utils.HexArrayToDecArray(strings))
			case "Bin":
				outputData = utils.BinArrayToString(utils.HexArrayToBinArray(strings))
			default:
				common.ResponseError(w, "Cannot convert '"+InputType+"' to '"+OutputType+"'")
				return
			}
		case "Dec":
			switch OutputType {
			case "Hex":
				outputData = utils.HexArrayToString(utils.DecArrayToHexArray(strings))
			case "Dec":
				outputData = utils.DecArrayToString(utils.DecArrayToDecArray(strings))
			case "Bin":
				outputData = utils.BinArrayToString(utils.DecArrayToBinArray(strings))
			default:
				common.ResponseError(w, "Cannot convert '"+InputType+"' to '"+OutputType+"'")
				return
			}
		case "Bin":
			switch OutputType {
			case "Hex":
				outputData = utils.HexArrayToString(utils.BinArrayToHexArray(strings))
			case "Dec":
				outputData = utils.DecArrayToString(utils.BinArrayToDecArray(strings))
			case "Bin":
				outputData = utils.BinArrayToString(strings)
			default:
				common.ResponseError(w, "Cannot convert '"+InputType+"' to '"+OutputType+"'")
				return
			}
		case "HexByteArray":
			switch OutputType {
			case "HexByteArray":
				outputData = utils.HexByteArrayToRows(strings)
			case "ByteArray":
				outputData = utils.ByteArrayToRows(utils.HexByteArrayToDecByteArray(strings))
			case "Int8Array":
				outputData = utils.Int8ArrayToRows(utils.HexByteArrayToInt8Array(strings))
			default:
				common.ResponseError(w, "Cannot convert '"+InputType+"' to '"+OutputType+"'")
				return
			}
		case "ByteArray":
			switch OutputType {
			case "HexByteArray":
				outputData = utils.HexByteArrayToRows(utils.DecByteArrayToHexByteArray(strings))
			case "ByteArray":
				outputData = utils.ByteArrayToRows(utils.DecByteArrayToDecByteArray(strings))
			case "Int8Array":
				outputData = utils.Int8ArrayToRows(utils.DecByteArrayToDecInt8Array(strings))
			default:
				common.ResponseError(w, "Cannot convert '"+InputType+"' to '"+OutputType+"'")
				return
			}
		case "Int8Array":
			switch OutputType {
			case "HexByteArray":
				outputData = utils.HexByteArrayToRows(utils.DecInt8ArrayToHexByteArray(strings))
			case "ByteArray":
				outputData = utils.ByteArrayToRows(utils.DecInt8ArrayToDecByteArray(strings))
			case "Int8Array":
				outputData = utils.Int8ArrayToRows(utils.DecInt8ArrayToDecInt8Array(strings))
			default:
				common.ResponseError(w, "Cannot convert '"+InputType+"' to '"+OutputType+"'")
				return
			}
		default:
			common.ResponseError(w, "Unknown input type: '"+InputType+"'")
			return
		}

		enc := json.NewEncoder(w)
		resData := common.ResData{
			Success: true,
			Message: "Data was converted!",
			Data:    outputData,
		}
		err = enc.Encode(resData)
		if err != nil {
			common.ResponseError(w, "Failed to encode data, error: "+err.Error())
			return
		}
	default:
		common.ResponseError(w, "Failed to convert data")
	}
}

func readStreamAndSendBody(w http.ResponseWriter, dataChan <-chan []byte, funcBytesToRow utils.BytesToRow) {
	transferSize := 0
	globalRowIndex := 0
	for {
		data, ok := <-dataChan
		if !ok || len(data) <= 0 {
			logger.Log("Read stream done, total size: " + strconv.Itoa(transferSize) + "Byte")
			return
		}
		rowsBytes := utils.StreamBytesToRowsBytes(data, &globalRowIndex, funcBytesToRow)
		w.Write(rowsBytes)
		transferSize += len(data)
		logger.Log("Read stream size: " + strconv.Itoa(transferSize) + "Byte")
	}
}
