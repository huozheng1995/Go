package controller

import (
	"encoding/json"
	"fmt"
	"github.com/edward/scp-294/common"
	"github.com/edward/scp-294/inout"
	"github.com/edward/scp-294/processor"
	"github.com/edward/scp-294/util"
	"mime/multipart"
	"myutil"
	myfile "myutil/file"
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
			common.RespondError(w, "Failed to parse form data, error: "+err.Error())
			return
		}
		form := r.MultipartForm

		var InputType inout.TypeCode
		var InputFormat, OutputFormat inout.FormatCode
		var InputData string
		for key, values := range form.Value {
			if key == "InputType" {
				intVal, _ := strconv.Atoi(values[0])
				InputType = inout.TypeCode(intVal)
			} else if key == "InputFormat" {
				intVal, _ := strconv.Atoi(values[0])
				InputFormat = inout.FormatCode(intVal)
			} else if key == "OutputFormat" {
				intVal, _ := strconv.Atoi(values[0])
				OutputFormat = inout.FormatCode(intVal)
			} else if key == "InputData" {
				InputData = values[0]
			}
		}

		if InputType == inout.File {
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

func convertText(InputData string, InputFormat, OutputFormat inout.FormatCode, w http.ResponseWriter) {
	int64Util := selectInt64Util(InputFormat)
	if int64Util != nil {
		int64ToStr := selectInt64Util(OutputFormat)
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
		byteArray := util.TextToNums[byte](InputData, byteUtil)

		processors := []string{processor.IbmEBCDICDecoder}
		byteArray, err := process(processors, byteArray)
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
	common.RespondError(w, fmt.Sprintln("Cannot convert ", InputFormat, " to ", OutputFormat))
}

func convertFile(file multipart.File, InputFormat, OutputFormat inout.FormatCode, w http.ResponseWriter) {
	int64File := selectInt64File(file, InputFormat)
	if int64File != nil {
		int64ToStr := selectInt64Util(OutputFormat)
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
	common.RespondError(w, fmt.Sprintln("Cannot convert ", InputFormat, " to ", OutputFormat))
}

func process(processors []string, arr []byte) ([]byte, error) {
	chain, err := processor.NewProcessorChain(processors)
	if err != nil {
		return nil, err
	}

	result, err := chain.Process(arr)
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

func selectInt64Util(format inout.FormatCode) (strToInt64 myutil.Int64Util) {
	switch format {
	case inout.Hex:
		strToInt64 = myutil.HexUtil{}
	case inout.Dec:
		strToInt64 = myutil.DecUtil{}
	case inout.Bin:
		strToInt64 = myutil.BinUtil{}
	default:
		strToInt64 = nil
	}
	return strToInt64
}

func selectInt64File(file multipart.File, format inout.FormatCode) (newFile *myfile.StrNumFile[int64]) {
	switch format {
	case inout.Hex:
		newFile = myfile.NewStrHexFile(file)
	case inout.Dec:
		newFile = myfile.NewStrDecFile(file)
	case inout.Bin:
		newFile = myfile.NewStrBinFile(file)
	default:
		newFile = nil
	}
	return newFile
}

func selectByteUtil(format inout.FormatCode) (byteToStr myutil.ByteUtil, withDetails bool) {
	withDetails = false
	switch format {
	case inout.HexByte:
		byteToStr = myutil.Hex8Util{}
	case inout.DecByte:
		byteToStr = myutil.Byte8Util{}
	case inout.DecInt8:
		byteToStr = myutil.Int8Util{}
	case inout.FormattedHexByte:
		withDetails = true
		byteToStr = myutil.Hex8Util{}
	case inout.FormattedDecByte:
		withDetails = true
		byteToStr = myutil.Byte8Util{}
	case inout.FormattedDecInt8:
		withDetails = true
		byteToStr = myutil.Int8Util{}
	case inout.RawBytes:
		byteToStr = myutil.RawBytesUtil{}
	default:
		byteToStr = nil
	}
	return byteToStr, withDetails
}

func selectByteFile(file multipart.File, format inout.FormatCode) (newFile myfile.INumFile[byte]) {
	switch format {
	case inout.HexByte:
		newFile = myfile.NewStrHex8File(file)
	case inout.DecByte:
		newFile = myfile.NewStrByte8File(file)
	case inout.DecInt8:
		newFile = myfile.NewStrInt8File(file)
	case inout.RawBytes:
		newFile = file
	default:
		newFile = nil
	}
	return newFile
}
