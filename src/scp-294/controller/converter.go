package controller

import (
	"encoding/json"
	"errors"
	"github.com/edward/scp-294/common"
	"github.com/edward/scp-294/converter"
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
			logger.Log(err.Error())
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
					file, _ = files[0].Open()
				}
			}
			exitChan, dataChan := converter.FileStreamToChannel(file, converter.GlobalRowSize*256)
			err = readStreamAndSendBody(w, dataChan, OutputType)
			if err != nil {
				close(exitChan)
				logger.Log(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				common.ResponseError(w, "Failed to parse file data, error: "+err.Error())
				return
			}
			close(exitChan)
			return
		}

		var outputData string
		var strings = converter.SplitInputString(InputData)
		switch InputType {
		case "Hex":
			switch OutputType {
			case "Hex":
				outputData = converter.HexArrayToString(strings)
			case "Dec":
				outputData = converter.DecArrayToString(converter.HexArrayToDecArray(strings))
			case "Bin":
				outputData = converter.BinArrayToString(converter.HexArrayToBinArray(strings))
			default:
				common.ResponseError(w, "Cannot convert '"+InputType+"' to '"+OutputType+"'")
				return
			}
		case "Dec":
			switch OutputType {
			case "Hex":
				outputData = converter.HexArrayToString(converter.DecArrayToHexArray(strings))
			case "Dec":
				outputData = converter.DecArrayToString(converter.DecArrayToDecArray(strings))
			case "Bin":
				outputData = converter.BinArrayToString(converter.DecArrayToBinArray(strings))
			default:
				common.ResponseError(w, "Cannot convert '"+InputType+"' to '"+OutputType+"'")
				return
			}
		case "Bin":
			switch OutputType {
			case "Hex":
				outputData = converter.HexArrayToString(converter.BinArrayToHexArray(strings))
			case "Dec":
				outputData = converter.DecArrayToString(converter.BinArrayToDecArray(strings))
			case "Bin":
				outputData = converter.BinArrayToString(strings)
			default:
				common.ResponseError(w, "Cannot convert '"+InputType+"' to '"+OutputType+"'")
				return
			}
		case "HexByteArray":
			switch OutputType {
			case "HexByteArray":
				outputData = converter.HexByteArrayToRows(strings)
			case "ByteArray":
				outputData = converter.ByteArrayToRows(converter.HexByteArrayToDecByteArray(strings))
			case "Int8Array":
				outputData = converter.Int8ArrayToRows(converter.HexByteArrayToInt8Array(strings))
			default:
				common.ResponseError(w, "Cannot convert '"+InputType+"' to '"+OutputType+"'")
				return
			}
		case "ByteArray":
			switch OutputType {
			case "HexByteArray":
				outputData = converter.HexByteArrayToRows(converter.DecByteArrayToHexByteArray(strings))
			case "ByteArray":
				outputData = converter.ByteArrayToRows(converter.DecByteArrayToDecByteArray(strings))
			case "Int8Array":
				outputData = converter.Int8ArrayToRows(converter.DecByteArrayToDecInt8Array(strings))
			default:
				common.ResponseError(w, "Cannot convert '"+InputType+"' to '"+OutputType+"'")
				return
			}
		case "Int8Array":
			switch OutputType {
			case "HexByteArray":
				outputData = converter.HexByteArrayToRows(converter.DecInt8ArrayToHexByteArray(strings))
			case "ByteArray":
				outputData = converter.ByteArrayToRows(converter.DecInt8ArrayToDecByteArray(strings))
			case "Int8Array":
				outputData = converter.Int8ArrayToRows(converter.DecInt8ArrayToDecInt8Array(strings))
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
			logger.Log(err.Error())
			common.ResponseError(w, "Failed to encode data, error: "+err.Error())
			return
		}
	default:
		common.ResponseError(w, "Failed to convert data")
	}
}

func readStreamAndSendBody(w http.ResponseWriter, dataChan <-chan []byte, outputType string) error {
	transferSize := 0
	globalRowIndex := 0
	var funcBytesToLine utils.BytesToRow
	switch outputType {
	case "HexByteArray":
		funcBytesToLine = utils.BytesToHexRow
	case "ByteArray":
		funcBytesToLine = utils.BytesToByteRow
	case "Int8Array":
		funcBytesToLine = utils.BytesToInt8Row
	default:
		return errors.New("Cannot convert File to '" + outputType + "'")
	}
	for {
		data, ok := <-dataChan
		if !ok || len(data) <= 0 {
			logger.Log("Read stream done, total size: " + strconv.Itoa(transferSize) + "Byte")
			return nil
		}
		rowsBytes := converter.StreamBytesToRowsBytes(data, &globalRowIndex, funcBytesToLine)
		w.Write(rowsBytes)
		transferSize += len(data)
		logger.Log("Read stream size: " + strconv.Itoa(transferSize) + "Byte")
	}
}
