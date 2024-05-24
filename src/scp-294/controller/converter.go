package controller

import (
	"encoding/json"
	"github.com/edward/scp-294/common"
	"github.com/edward/scp-294/logger"
	"github.com/edward/scp-294/util"
	"mime/multipart"
	"myutil"
	myfile "myutil/file"
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
			common.RespondError(w, "Failed to parse form data, error: "+err.Error())
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
						common.RespondError(w, "Failed to open file, error: "+err.Error())
						return
					}
				}
			}
			if file != nil {
				defer file.Close()
			}
			convertFile(file, InputFormat, OutputFormat, w)
		} else {
			convertText(InputData, InputFormat, OutputFormat, w)
		}

	default:
		common.RespondError(w, "Failed to convert data")
	}
}

func convertText(InputData string, InputFormat, OutputFormat common.NumType, w http.ResponseWriter) {
	var response string
	funcStrToInt64 := selectFuncStrToInt64(InputFormat)
	if funcStrToInt64 != nil {
		funcInt64ToStr := selectFuncInt64ToStr(OutputFormat)
		if funcInt64ToStr == nil {
			w.WriteHeader(http.StatusInternalServerError)
			common.RespondError(w, "Cannot convert '"+common.InputFormatMap[InputFormat]+"' to '"+common.OutputFormatMap[OutputFormat]+"'")
			return
		}
		int64Array := util.TextToNums[int64](InputData, funcStrToInt64)
		numsToResp := &util.Int64sToResp{
			NumToStr: funcInt64ToStr,
		}
		resPageBuf := numsToResp.ToResp(int64Array)
		response = string(resPageBuf.Bytes())
		writeResponse(w, response)
		return
	}

	funcStrToByte := selectFuncStrToByte(InputFormat)
	if funcStrToByte != nil {
		funcByteToStr, withDetails := selectFuncByteToStr(OutputFormat)
		if funcByteToStr == nil {
			w.WriteHeader(http.StatusInternalServerError)
			common.RespondError(w, "Cannot convert '"+common.InputFormatMap[InputFormat]+"' to '"+common.OutputFormatMap[OutputFormat]+"'")
			return
		}
		byteArray := util.TextToNums[byte](InputData, funcStrToByte)
		numsToResp := &util.BytesToResp{
			NumToStr:       funcByteToStr,
			WithDetails:    withDetails,
			GlobalRowIndex: 0,
		}
		resPageBuf := numsToResp.ToResp(byteArray)
		response = string(resPageBuf.Bytes())
		writeResponse(w, response)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	common.RespondError(w, "Cannot convert '"+common.InputFormatMap[InputFormat]+"' to '"+common.OutputFormatMap[OutputFormat]+"'")
}

func writeResponse(w http.ResponseWriter, response string) {
	enc := json.NewEncoder(w)
	resData := common.ResData{
		Success: true,
		Message: "Data was converted!",
		Data:    response,
	}
	err := enc.Encode(resData)
	if err != nil {
		common.RespondError(w, "Failed to encode data, error: "+err.Error())
		return
	}
}

func convertFile(file multipart.File, InputFormat, OutputFormat common.NumType, w http.ResponseWriter) {
	//Int64
	strToInt64File := selectStrToInt64File(file, InputFormat)
	if strToInt64File != nil {
		funcInt64ToStr := selectFuncInt64ToStr(OutputFormat)
		if funcInt64ToStr == nil {
			w.WriteHeader(http.StatusInternalServerError)
			common.RespondError(w, "Cannot convert '"+common.InputFormatMap[InputFormat]+"' to '"+common.OutputFormatMap[OutputFormat]+"'")
			return
		}
		readChan := make(chan []int64)
		go util.StrNumFileToNums(strToInt64File, reqInt64BufferPool, readChan)
		numsToResp := &util.Int64sToResp{
			NumToStr: funcInt64ToStr,
		}
		util.ReadFromChannelAndRespond(readChan, numsToResp, reqInt64BufferPool, w)
		return
	}

	//Bytes
	strToByteFile := selectStrToByteFile(file, InputFormat)
	if strToByteFile != nil {
		funcByteToStr, withDetails := selectFuncByteToStr(OutputFormat)
		if funcByteToStr == nil {
			w.WriteHeader(http.StatusInternalServerError)
			common.RespondError(w, "Cannot convert '"+common.InputFormatMap[InputFormat]+"' to '"+common.OutputFormatMap[OutputFormat]+"'")
			return
		}
		readChan := make(chan []byte)
		go util.StrNumFileToNums(strToByteFile, reqByteBufferPool, readChan)
		numsToResp := &util.BytesToResp{
			NumToStr:       funcByteToStr,
			WithDetails:    withDetails,
			GlobalRowIndex: 0,
		}
		util.ReadFromChannelAndRespond(readChan, numsToResp, reqByteBufferPool, w)
		return
	}

	//RawBytes
	if InputFormat == common.RawBytes {
		funcByteToStr, withDetails := selectFuncByteToStr(OutputFormat)
		if funcByteToStr == nil {
			w.WriteHeader(http.StatusInternalServerError)
			common.RespondError(w, "Cannot convert '"+common.InputFormatMap[InputFormat]+"' to '"+common.OutputFormatMap[OutputFormat]+"'")
			return
		}
		readChan := make(chan []byte)
		go util.RawBytesFileToBytes(file, reqByteBufferPool, readChan)
		numsToResp := &util.BytesToResp{
			NumToStr:       funcByteToStr,
			WithDetails:    withDetails,
			GlobalRowIndex: 0,
		}
		util.ReadFromChannelAndRespond(readChan, numsToResp, reqByteBufferPool, w)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	common.RespondError(w, "Cannot convert '"+common.InputFormatMap[InputFormat]+"' to '"+common.OutputFormatMap[OutputFormat]+"'")
}

var byteBufferCount int32
var reqByteBufferPool = &sync.Pool{
	New: func() interface{} {
		atomic.AddInt32(&byteBufferCount, 1)
		logger.Logger.Log("Main", "reqByteBufferPool: Count of new buffer: "+strconv.Itoa(int(byteBufferCount)))
		return make([]byte, 4096)
	},
}

var int64BufferCount int32
var reqInt64BufferPool = &sync.Pool{
	New: func() interface{} {
		atomic.AddInt32(&int64BufferCount, 1)
		logger.Logger.Log("Main", "reqInt64BufferPool: Count of new buffer: "+strconv.Itoa(int(int64BufferCount)))
		return make([]int64, 4096>>3)
	},
}

func selectFuncStrToInt64(InputFormat common.NumType) (funcStrToInt64 myutil.StrToInt64) {
	switch InputFormat {
	case common.Hex:
		funcStrToInt64 = myutil.HexStrToInt64{}
	case common.Dec:
		funcStrToInt64 = myutil.DecStrToInt64{}
	case common.Bin:
		funcStrToInt64 = myutil.BinStrToInt64{}
	default:
		funcStrToInt64 = nil
	}
	return funcStrToInt64
}

func selectStrToInt64File(file multipart.File, InputFormat common.NumType) (newFile *myfile.StrNumFile[int64]) {
	switch InputFormat {
	case common.Hex:
		newFile = myfile.NewStrHexFile(file)
	case common.Dec:
		newFile = myfile.NewStrDecFile(file)
	case common.Bin:
		newFile = myfile.NewStrBinFile(file)
	default:
		newFile = nil
	}
	return newFile
}

func selectFuncInt64ToStr(OutputFormat common.NumType) (funcInt64ToStr myutil.Int64ToStr) {
	switch OutputFormat {
	case common.Hex:
		funcInt64ToStr = myutil.Int64ToHexStr{}
	case common.Dec:
		funcInt64ToStr = myutil.Int64ToInt64Str{}
	case common.Bin:
		funcInt64ToStr = myutil.Int64ToBinStr{}
	default:
		funcInt64ToStr = nil
	}
	return funcInt64ToStr
}

func selectFuncStrToByte(InputFormat common.NumType) (funcStrToByte myutil.StrToByte) {
	switch InputFormat {
	case common.HexByte:
		funcStrToByte = myutil.Hex2StrToByte{}
	case common.DecByte:
		funcStrToByte = myutil.ByteStrToByte{}
	case common.DecInt8:
		funcStrToByte = myutil.Int8StrToByte{}
	default:
		funcStrToByte = nil
	}
	return funcStrToByte
}

func selectStrToByteFile(file multipart.File, InputFormat common.NumType) (newFile *myfile.StrNumFile[byte]) {
	switch InputFormat {
	case common.HexByte:
		newFile = myfile.NewStrHex2File(file)
	case common.DecByte:
		newFile = myfile.NewStrByteFile(file)
	case common.DecInt8:
		newFile = myfile.NewStrInt8File(file)
	default:
		newFile = nil
	}
	return newFile
}

func selectFuncByteToStr(OutputFormat common.NumType) (funcByteToStr myutil.ByteToStr, withDetails bool) {
	withDetails = false
	switch OutputFormat {
	case common.HexByte:
		funcByteToStr = myutil.ByteToHexStr{}
	case common.DecByte:
		funcByteToStr = myutil.ByteToByteStr{}
	case common.DecInt8:
		funcByteToStr = myutil.ByteToInt8Str{}
	case common.HexByteFormatted:
		withDetails = true
		funcByteToStr = myutil.ByteToHexStr{}
	case common.DecByteFormatted:
		withDetails = true
		funcByteToStr = myutil.ByteToByteStr{}
	case common.DecInt8Formatted:
		withDetails = true
		funcByteToStr = myutil.ByteToInt8Str{}
	case common.RawBytes:
		funcByteToStr = myutil.ByteToRawBytes{}
	default:
		funcByteToStr = nil
	}
	return funcByteToStr, withDetails
}
