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
	strToInt64 := selectStrToInt64(InputFormat)
	if strToInt64 != nil {
		int64ToStr := selectInt64ToStr(OutputFormat)
		if int64ToStr == nil {
			w.WriteHeader(http.StatusInternalServerError)
			common.RespondError(w, "Cannot convert '"+common.InputFormatMap[InputFormat]+"' to '"+common.OutputFormatMap[OutputFormat]+"'")
			return
		}
		int64Array := util.TextToNums[int64](InputData, strToInt64)
		numsToResp := &util.Int64sToResp{
			NumToStr: int64ToStr,
		}
		resp := numsToResp.ToResp(int64Array)
		writeResponse(w, string(resp))
		return
	}

	strToByte := selectStrToByte(InputFormat)
	if strToByte != nil {
		byteToStr, withDetails := selectByteToStr(OutputFormat)
		if byteToStr == nil {
			w.WriteHeader(http.StatusInternalServerError)
			common.RespondError(w, "Cannot convert '"+common.InputFormatMap[InputFormat]+"' to '"+common.OutputFormatMap[OutputFormat]+"'")
			return
		}
		byteArray := util.TextToNums[byte](InputData, strToByte)
		numsToResp := &util.BytesToResp{
			NumToStr:       byteToStr,
			WithDetails:    withDetails,
			GlobalRowIndex: 0,
		}
		resp := numsToResp.ToResp(byteArray)
		writeResponse(w, string(resp))
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
		int64ToStr := selectInt64ToStr(OutputFormat)
		if int64ToStr == nil {
			w.WriteHeader(http.StatusInternalServerError)
			common.RespondError(w, "Cannot convert '"+common.InputFormatMap[InputFormat]+"' to '"+common.OutputFormatMap[OutputFormat]+"'")
			return
		}
		readChan := make(chan []int64)
		go util.FileToNums[int64](strToInt64File, reqInt64BufferPool, readChan)
		numsToResp := &util.Int64sToResp{
			NumToStr: int64ToStr,
		}
		util.ReadFromChannelAndRespond(readChan, numsToResp, reqInt64BufferPool, w)
		return
	}

	//Bytes
	strToByteFile := selectStrToByteFile(file, InputFormat)
	if strToByteFile != nil {
		byteToStr, withDetails := selectByteToStr(OutputFormat)
		if byteToStr == nil {
			w.WriteHeader(http.StatusInternalServerError)
			common.RespondError(w, "Cannot convert '"+common.InputFormatMap[InputFormat]+"' to '"+common.OutputFormatMap[OutputFormat]+"'")
			return
		}
		readChan := make(chan []byte)
		go util.FileToNums[byte](strToByteFile, reqByteBufferPool, readChan)
		numsToResp := &util.BytesToResp{
			NumToStr:       byteToStr,
			WithDetails:    withDetails,
			GlobalRowIndex: 0,
		}
		util.ReadFromChannelAndRespond(readChan, numsToResp, reqByteBufferPool, w)
		return
	}

	//RawBytes
	if InputFormat == common.RawBytes {
		byteToStr, withDetails := selectByteToStr(OutputFormat)
		if byteToStr == nil {
			w.WriteHeader(http.StatusInternalServerError)
			common.RespondError(w, "Cannot convert '"+common.InputFormatMap[InputFormat]+"' to '"+common.OutputFormatMap[OutputFormat]+"'")
			return
		}
		readChan := make(chan []byte)
		go util.FileToNums[byte](file, reqByteBufferPool, readChan)
		numsToResp := &util.BytesToResp{
			NumToStr:       byteToStr,
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

func selectStrToInt64(InputFormat common.NumType) (strToInt64 myutil.Int64Util) {
	switch InputFormat {
	case common.Hex:
		strToInt64 = myutil.HexUtil{}
	case common.Dec:
		strToInt64 = myutil.DecUtil{}
	case common.Bin:
		strToInt64 = myutil.BinUtil{}
	default:
		strToInt64 = nil
	}
	return strToInt64
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

func selectInt64ToStr(OutputFormat common.NumType) (int64ToStr myutil.Int64Util) {
	switch OutputFormat {
	case common.Hex:
		int64ToStr = myutil.HexUtil{}
	case common.Dec:
		int64ToStr = myutil.DecUtil{}
	case common.Bin:
		int64ToStr = myutil.BinUtil{}
	default:
		int64ToStr = nil
	}
	return int64ToStr
}

func selectStrToByte(InputFormat common.NumType) (strToByte myutil.ByteUtil) {
	switch InputFormat {
	case common.HexByte:
		strToByte = myutil.Hex8Util{}
	case common.DecByte:
		strToByte = myutil.Byte8Util{}
	case common.DecInt8:
		strToByte = myutil.Int8Util{}
	default:
		strToByte = nil
	}
	return strToByte
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

func selectByteToStr(OutputFormat common.NumType) (byteToStr myutil.ByteUtil, withDetails bool) {
	withDetails = false
	switch OutputFormat {
	case common.HexByte:
		byteToStr = myutil.Hex8Util{}
	case common.DecByte:
		byteToStr = myutil.Byte8Util{}
	case common.DecInt8:
		byteToStr = myutil.Int8Util{}
	case common.HexByteFormatted:
		withDetails = true
		byteToStr = myutil.Hex8Util{}
	case common.DecByteFormatted:
		withDetails = true
		byteToStr = myutil.Byte8Util{}
	case common.DecInt8Formatted:
		withDetails = true
		byteToStr = myutil.Int8Util{}
	case common.RawBytes:
		byteToStr = myutil.RawBytesUtil{}
	default:
		byteToStr = nil
	}
	return byteToStr, withDetails
}
