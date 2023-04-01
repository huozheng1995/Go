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

		var InputType, InputFormat, OutputFormat common.NumType
		var InputData string
		for key, values := range form.Value {
			if key == "InputType" {
				intVal, _ := strconv.Atoi(values[0])
				InputType = common.NumType(intVal)
			} else if key == "InputFormat" {
				intVal, _ := strconv.Atoi(values[0])
				InputFormat = common.NumType(intVal)
			} else if key == "OutputFormat" {
				intVal, _ := strconv.Atoi(values[0])
				OutputFormat = common.NumType(intVal)
			} else if key == "InputData" {
				InputData = values[0]
			}
		}

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
			defer file.Close()
			convertFile(w, file, InputFormat, OutputFormat)
		} else {
			convertText(w, InputData, InputFormat, OutputFormat)
		}

	default:
		common.ResponseError(w, "Failed to convert data")
	}
}

func convertText(w http.ResponseWriter, InputData string, InputFormat, OutputFormat common.NumType) {
	var outputData string

	funcStrToInt64 := selectFuncStrToInt64(InputFormat)
	if funcStrToInt64 != nil {
		funcInt64ToStr := selectFuncInt64ToStr(OutputFormat)
		if funcInt64ToStr == nil {
			w.WriteHeader(http.StatusInternalServerError)
			common.ResponseError(w, "Cannot convert '"+common.InputFormatMap[InputFormat]+"' to '"+common.OutputFormatMap[OutputFormat]+"'")
			return
		}
		int64Array := utils.StringToInt64Array(InputData, funcStrToInt64)
		outputData = utils.Int64ArrayToRowString(int64Array, funcInt64ToStr)
		toResponse(w, outputData)
		return
	}

	funcStrToByte := selectFuncStrToByte(InputFormat)
	if funcStrToInt64 != nil {
		funcByteToStr, withDetails := selectFuncByteToStr(OutputFormat)
		if funcByteToStr == nil {
			w.WriteHeader(http.StatusInternalServerError)
			common.ResponseError(w, "Cannot convert '"+common.InputFormatMap[InputFormat]+"' to '"+common.OutputFormatMap[OutputFormat]+"'")
			return
		}
		byteArray := utils.StringToByteArray(InputData, funcStrToByte)
		globalRowIndex := 0
		outputData = utils.ByteArrayToRowString(byteArray, &globalRowIndex, funcByteToStr, withDetails)
		toResponse(w, outputData)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	common.ResponseError(w, "Cannot convert '"+common.InputFormatMap[InputFormat]+"' to '"+common.OutputFormatMap[OutputFormat]+"'")
}

func toResponse(w http.ResponseWriter, outputData string) {
	enc := json.NewEncoder(w)
	resData := common.ResData{
		Success: true,
		Message: "Data was converted!",
		Data:    outputData,
	}
	err := enc.Encode(resData)
	if err != nil {
		common.ResponseError(w, "Failed to encode data, error: "+err.Error())
		return
	}
}

func convertFile(w http.ResponseWriter, file multipart.File, InputFormat, OutputFormat common.NumType) {
	funcStrToInt64 := selectFuncStrToInt64(InputFormat)
	if funcStrToInt64 != nil {
		funcInt64ToStr := selectFuncInt64ToStr(OutputFormat)
		if funcInt64ToStr == nil {
			w.WriteHeader(http.StatusInternalServerError)
			common.ResponseError(w, "Cannot convert '"+common.InputFormatMap[InputFormat]+"' to '"+common.OutputFormatMap[OutputFormat]+"'")
			return
		}
		logger.Log("Begin parse, buffer count: " + strconv.Itoa(int(int64BufferCount)))
		exitChan, readChan := utils.FileStreamToInt64Array(file, int64BufferPool, funcStrToInt64)
		defer close(exitChan)
		utils.ReadInt64ArrayAndSendBody(w, readChan, funcInt64ToStr, int64BufferPool)
		logger.Log("End parse, buffer count: " + strconv.Itoa(int(int64BufferCount)))
		return
	}

	funcStrToByte := selectFuncStrToByte(InputFormat)
	if funcStrToByte != nil {
		funcByteToStr, withDetails := selectFuncByteToStr(OutputFormat)
		if funcByteToStr == nil {
			w.WriteHeader(http.StatusInternalServerError)
			common.ResponseError(w, "Cannot convert '"+common.InputFormatMap[InputFormat]+"' to '"+common.OutputFormatMap[OutputFormat]+"'")
			return
		}
		logger.Log("Begin parse, buffer count: " + strconv.Itoa(int(byteBufferCount)))
		exitChan, readChan := utils.FileStreamToPageBytes(file, byteBufferPool, funcStrToByte)
		defer close(exitChan)
		utils.ReadBytesAndSendBody(w, readChan, funcByteToStr, withDetails, byteBufferPool)
		logger.Log("End parse, buffer count: " + strconv.Itoa(int(byteBufferCount)))
		return
	}

	if InputFormat == common.RawBytes {
		funcByteToStr, withDetails := selectFuncByteToStr(OutputFormat)
		if funcByteToStr == nil {
			w.WriteHeader(http.StatusInternalServerError)
			common.ResponseError(w, "Cannot convert '"+common.InputFormatMap[InputFormat]+"' to '"+common.OutputFormatMap[OutputFormat]+"'")
			return
		}
		logger.Log("Begin parse, buffer count: " + strconv.Itoa(int(byteBufferCount)))
		exitChan, readChan := utils.FileStreamToRawBytes(file, byteBufferPool)
		defer close(exitChan)
		utils.ReadBytesAndSendBody(w, readChan, funcByteToStr, withDetails, byteBufferPool)
		logger.Log("End parse, buffer count: " + strconv.Itoa(int(byteBufferCount)))
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	common.ResponseError(w, "Cannot convert '"+common.InputFormatMap[InputFormat]+"' to '"+common.OutputFormatMap[OutputFormat]+"'")
}

var int64BufferCount int32
var int64BufferPool = &sync.Pool{
	New: func() interface{} {
		atomic.AddInt32(&int64BufferCount, 1)
		return make([]int64, 512)
	},
}

var byteBufferCount int32
var byteBufferPool = &sync.Pool{
	New: func() interface{} {
		atomic.AddInt32(&byteBufferCount, 1)
		return make([]byte, 4096)
	},
}

func selectFuncStrToInt64(InputFormat common.NumType) (funcStrToInt64 utils.StrToInt64) {
	switch InputFormat {
	case common.Hex:
		funcStrToInt64 = utils.HexStrToInt64
	case common.Dec:
		funcStrToInt64 = utils.Int64StrToInt64
	case common.Bin:
		funcStrToInt64 = utils.BinStrToInt64
	default:
		funcStrToInt64 = nil
	}
	return funcStrToInt64
}

func selectFuncInt64ToStr(OutputFormat common.NumType) (funcInt64ToStr utils.Int64ToStr) {
	switch OutputFormat {
	case common.Hex:
		funcInt64ToStr = utils.Int64ToHexStr
	case common.Dec:
		funcInt64ToStr = utils.Int64ToInt64Str
	case common.Bin:
		funcInt64ToStr = utils.Int64ToBinStr
	default:
		funcInt64ToStr = nil
	}
	return funcInt64ToStr
}

func selectFuncStrToByte(InputFormat common.NumType) (funcStrToByte utils.StrToByte) {
	switch InputFormat {
	case common.HexByte:
		funcStrToByte = utils.HexStrToByte
	case common.DecByte:
		funcStrToByte = utils.ByteStrToByte
	case common.DecInt8:
		funcStrToByte = utils.Int8StrToByte
	default:
		funcStrToByte = nil
	}
	return funcStrToByte
}

func selectFuncByteToStr(OutputFormat common.NumType) (funcByteToStr utils.ByteToStr, withDetails bool) {
	withDetails = false
	switch OutputFormat {
	case common.HexByte:
		funcByteToStr = utils.ByteToHexStr{}
	case common.DecByte:
		funcByteToStr = utils.ByteToByteStr{}
	case common.DecInt8:
		funcByteToStr = utils.ByteToInt8Str{}
	case common.HexByteFormatted:
		withDetails = true
		funcByteToStr = utils.ByteToHexStr{}
	case common.DecByteFormatted:
		withDetails = true
		funcByteToStr = utils.ByteToByteStr{}
	case common.DecInt8Formatted:
		withDetails = true
		funcByteToStr = utils.ByteToInt8Str{}
	default:
		funcByteToStr = nil
	}
	return funcByteToStr, withDetails
}
