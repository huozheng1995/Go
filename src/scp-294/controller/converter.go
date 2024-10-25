package controller

import (
	"encoding/json"
	"github.com/edward/scp-294/common"
	"github.com/edward/scp-294/logger"
	"github.com/edward/scp-294/processor"
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
	int64Util := selectInt64Util(InputFormat)
	if int64Util != nil {
		int64ToStr := selectInt64Util(OutputFormat)
		if int64ToStr == nil {
			w.WriteHeader(http.StatusInternalServerError)
			common.RespondError(w, "Cannot convert '"+common.InputFormatMap[InputFormat]+"' to '"+common.OutputFormatMap[OutputFormat]+"'")
			return
		}
		int64Array := util.TextToNums[int64](InputData, int64Util)
		numsToResp := &util.Int64sToResp{
			NumToStr: int64ToStr,
		}
		resp := numsToResp.ToResp(int64Array)
		writeResponse(w, string(resp))
		return
	}

	byteUtil, _ := selectByteUtil(InputFormat)
	if byteUtil != nil {
		byteToStr, withDetails := selectByteUtil(OutputFormat)
		if byteToStr == nil {
			w.WriteHeader(http.StatusInternalServerError)
			common.RespondError(w, "Cannot convert '"+common.InputFormatMap[InputFormat]+"' to '"+common.OutputFormatMap[OutputFormat]+"'")
			return
		}
		byteArray := util.TextToNums[byte](InputData, byteUtil)

		processors := []string{processor.IbmEBCDICDecoder}
		byteArray, err := processByProcessors(processors, byteArray)
		if err != nil {
			common.RespondError(w, err.Error())
			return
		}

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

var registry = processor.NewProcessorRegistry()

func processByProcessors(processors []string, input []byte) ([]byte, error) {
	chain := processor.NewProcessorChain()
	for _, name := range processors {
		processor, err := registry.GetProcessor(name)
		if err != nil {
			return nil, err
		}
		chain.AddProcessor(processor)
	}

	result, err := chain.Process(input)
	if err != nil {
		return nil, err
	}

	return result, nil
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
	int64File := selectInt64File(file, InputFormat)
	if int64File != nil {
		int64ToStr := selectInt64Util(OutputFormat)
		if int64ToStr == nil {
			w.WriteHeader(http.StatusInternalServerError)
			common.RespondError(w, "Cannot convert '"+common.InputFormatMap[InputFormat]+"' to '"+common.OutputFormatMap[OutputFormat]+"'")
			return
		}
		readChan := make(chan []int64)
		go util.FileToNums[int64](int64File, reqInt64BufferPool, readChan)
		numsToResp := &util.Int64sToResp{
			NumToStr: int64ToStr,
		}
		util.ReadFromChannelAndRespond(readChan, numsToResp, reqInt64BufferPool, w)
		return
	}

	byteFile := selectByteFile(file, InputFormat)
	if byteFile != nil {
		byteToStr, withDetails := selectByteUtil(OutputFormat)
		if byteToStr == nil {
			w.WriteHeader(http.StatusInternalServerError)
			common.RespondError(w, "Cannot convert '"+common.InputFormatMap[InputFormat]+"' to '"+common.OutputFormatMap[OutputFormat]+"'")
			return
		}
		readChan := make(chan []byte)
		go util.FileToNums[byte](byteFile, reqByteBufferPool, readChan)
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

func selectInt64Util(format common.NumType) (strToInt64 myutil.Int64Util) {
	switch format {
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

func selectInt64File(file multipart.File, format common.NumType) (newFile *myfile.StrNumFile[int64]) {
	switch format {
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

func selectByteUtil(format common.NumType) (byteToStr myutil.ByteUtil, withDetails bool) {
	withDetails = false
	switch format {
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

func selectByteFile(file multipart.File, format common.NumType) (newFile myfile.INumFile[byte]) {
	switch format {
	case common.HexByte:
		newFile = myfile.NewStrHex8File(file)
	case common.DecByte:
		newFile = myfile.NewStrByte8File(file)
	case common.DecInt8:
		newFile = myfile.NewStrInt8File(file)
	case common.RawBytes:
		newFile = file
	default:
		newFile = nil
	}
	return newFile
}
