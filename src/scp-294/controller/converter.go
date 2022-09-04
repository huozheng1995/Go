package controller

import (
	"encoding/json"
	"github.com/edward/scp-294/common"
	"github.com/edward/scp-294/logger"
	"github.com/edward/scp-294/utils"
	"mime/multipart"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
)

func registerConverterRoutes() {
	http.HandleFunc("/convert", convert)
}

var fileBufferCount int32
var fileBufferPool = &sync.Pool{
	New: func() interface{} {
		atomic.AddInt32(&fileBufferCount, 1)
		return make([]byte, 4096)
	},
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

		var InputType, OutputType common.NumType
		var InputData string
		for key, values := range form.Value {
			if key == "InputType" {
				intVal, _ := strconv.Atoi(values[0])
				InputType = common.NumType(intVal)
			} else if key == "OutputType" {
				intVal, _ := strconv.Atoi(values[0])
				OutputType = common.NumType(intVal)
			} else if key == "InputData" {
				InputData = values[0]
			}
		}

		//Handle File
		if InputType == common.File {
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
			var funcBytesToRow utils.ByteArrayToRow
			var format = false
			switch OutputType {
			case common.HexByte:
				funcBytesToRow = utils.ByteArrayToHexByteRow
			case common.DecByte:
				funcBytesToRow = utils.ByteArrayToByteRow
			case common.DecInt8:
				funcBytesToRow = utils.ByteArrayToInt8Row
			case common.HexByteFormatted:
				format = true
				funcBytesToRow = utils.ByteArrayToHexByteRow
			case common.DecByteFormatted:
				format = true
				funcBytesToRow = utils.ByteArrayToByteRow
			case common.DecInt8Formatted:
				format = true
				funcBytesToRow = utils.ByteArrayToInt8Row
			default:
				w.WriteHeader(http.StatusInternalServerError)
				common.ResponseError(w, "Cannot convert '"+common.InputTypeMap[InputType]+"' to '"+common.OutputTypeMap[OutputType]+"'")
				return
			}
			logger.Log("Begin parse, buffer count: " + strconv.Itoa(int(fileBufferCount)))
			exitChannel, readChan := utils.FileStreamToChannel(file, fileBufferPool)
			readStreamAndSendBody(w, readChan, funcBytesToRow, format, fileBufferPool)
			logger.Log("End parse, buffer count: " + strconv.Itoa(int(fileBufferCount)))
			close(exitChannel)
			return
		}

		var outputData string
		var strings = utils.SplitInputString(InputData)
		switch InputType {
		case common.Hex:
			switch OutputType {
			case common.Hex:
				outputData = utils.HexArrayToString(strings)
			case common.Dec:
				outputData = utils.DecArrayToString(utils.HexArrayToDecArray(strings))
			case common.Bin:
				outputData = utils.BinArrayToString(utils.HexArrayToBinArray(strings))
			default:
				common.ResponseError(w, "Cannot convert '"+common.InputTypeMap[InputType]+"' to '"+common.OutputTypeMap[OutputType]+"'")
				return
			}
		case common.Dec:
			switch OutputType {
			case common.Hex:
				outputData = utils.HexArrayToString(utils.DecArrayToHexArray(strings))
			case common.Dec:
				outputData = utils.DecArrayToString(utils.DecArrayToDecArray(strings))
			case common.Bin:
				outputData = utils.BinArrayToString(utils.DecArrayToBinArray(strings))
			default:
				common.ResponseError(w, "Cannot convert '"+common.InputTypeMap[InputType]+"' to '"+common.OutputTypeMap[OutputType]+"'")
				return
			}
		case common.Bin:
			switch OutputType {
			case common.Hex:
				outputData = utils.HexArrayToString(utils.BinArrayToHexArray(strings))
			case common.Dec:
				outputData = utils.DecArrayToString(utils.BinArrayToDecArray(strings))
			case common.Bin:
				outputData = utils.BinArrayToString(strings)
			default:
				common.ResponseError(w, "Cannot convert '"+common.InputTypeMap[InputType]+"' to '"+common.OutputTypeMap[OutputType]+"'")
				return
			}
		case common.HexByte:
			switch OutputType {
			case common.HexByte:
				outputData = utils.HexByteArrayToRows(strings, false)
			case common.DecByte:
				outputData = utils.ByteArrayToRows(utils.HexByteArrayToDecByteArray(strings), false)
			case common.DecInt8:
				outputData = utils.Int8ArrayToRows(utils.HexByteArrayToInt8Array(strings), false)
			case common.HexByteFormatted:
				outputData = utils.HexByteArrayToRows(strings, true)
			case common.DecByteFormatted:
				outputData = utils.ByteArrayToRows(utils.HexByteArrayToDecByteArray(strings), true)
			case common.DecInt8Formatted:
				outputData = utils.Int8ArrayToRows(utils.HexByteArrayToInt8Array(strings), true)
			default:
				common.ResponseError(w, "Cannot convert '"+common.InputTypeMap[InputType]+"' to '"+common.OutputTypeMap[OutputType]+"'")
				return
			}
		case common.DecByte:
			switch OutputType {
			case common.HexByte:
				outputData = utils.HexByteArrayToRows(utils.DecByteArrayToHexByteArray(strings), false)
			case common.DecByte:
				outputData = utils.ByteArrayToRows(utils.DecByteArrayToDecByteArray(strings), false)
			case common.DecInt8:
				outputData = utils.Int8ArrayToRows(utils.DecByteArrayToDecInt8Array(strings), false)
			case common.HexByteFormatted:
				outputData = utils.HexByteArrayToRows(utils.DecByteArrayToHexByteArray(strings), true)
			case common.DecByteFormatted:
				outputData = utils.ByteArrayToRows(utils.DecByteArrayToDecByteArray(strings), true)
			case common.DecInt8Formatted:
				outputData = utils.Int8ArrayToRows(utils.DecByteArrayToDecInt8Array(strings), true)
			default:
				common.ResponseError(w, "Cannot convert '"+common.InputTypeMap[InputType]+"' to '"+common.OutputTypeMap[OutputType]+"'")
				return
			}
		case common.DecInt8:
			switch OutputType {
			case common.HexByte:
				outputData = utils.HexByteArrayToRows(utils.DecInt8ArrayToHexByteArray(strings), false)
			case common.DecByte:
				outputData = utils.ByteArrayToRows(utils.DecInt8ArrayToDecByteArray(strings), false)
			case common.DecInt8:
				outputData = utils.Int8ArrayToRows(utils.DecInt8ArrayToDecInt8Array(strings), false)
			case common.HexByteFormatted:
				outputData = utils.HexByteArrayToRows(utils.DecInt8ArrayToHexByteArray(strings), true)
			case common.DecByteFormatted:
				outputData = utils.ByteArrayToRows(utils.DecInt8ArrayToDecByteArray(strings), true)
			case common.DecInt8Formatted:
				outputData = utils.Int8ArrayToRows(utils.DecInt8ArrayToDecInt8Array(strings), true)
			default:
				common.ResponseError(w, "Cannot convert '"+common.InputTypeMap[InputType]+"' to '"+common.OutputTypeMap[OutputType]+"'")
				return
			}
		default:
			common.ResponseError(w, "Unknown input type: '"+common.InputTypeMap[InputType]+"'")
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

func readStreamAndSendBody(w http.ResponseWriter, readChan <-chan []byte,
	funcBytesToRow utils.ByteArrayToRow, format bool, bufferPool *sync.Pool) {
	readSize := 0
	writeSize := 0
	globalRowIndex := 0
	for {
		data, ok := <-readChan
		if !ok || len(data) <= 0 {
			logger.Log("Read channel done, total size: " + strconv.Itoa(readSize) + "Byte")
			logger.Log("Write stream done, total size: " + strconv.Itoa(writeSize) + "Byte")
			return
		}
		rowsBytes := utils.StreamBytesToRowsBytes(data, &globalRowIndex, funcBytesToRow, format)
		bufferPool.Put(data)
		w.Write(rowsBytes)
		readSize += len(data)
		writeSize += len(rowsBytes)
		//logger.Log("Read stream size: " + strconv.Itoa(readSize) + "Byte")
	}
}
