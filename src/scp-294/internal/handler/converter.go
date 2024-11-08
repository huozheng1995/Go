package handler

import (
	"encoding/json"
	"fmt"
	"github.com/edward/scp-294/internal/constants"
	"github.com/edward/scp-294/pkg/processor"
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
			respondError(w, "Failed to parse form data, error: "+err.Error())
			return
		}
		form := r.MultipartForm

		var InputType constants.TypeCode
		var InputFormat, OutputFormat constants.FormatCode
		processors := make([]string, 0)
		var InputData string
		for key, values := range form.Value {
			if key == "InputType" {
				intVal, _ := strconv.Atoi(values[0])
				InputType = constants.TypeCode(intVal)
			} else if key == "InputFormat" {
				intVal, _ := strconv.Atoi(values[0])
				InputFormat = constants.FormatCode(intVal)
			} else if key == "OutputFormat" {
				intVal, _ := strconv.Atoi(values[0])
				OutputFormat = constants.FormatCode(intVal)
			} else if key == "processor" {
				processors = append(processors, values[0])
			} else if key == "InputData" {
				InputData = values[0]
			}
		}

		if InputType == constants.File {
			var file multipart.File
			for key, files := range form.File {
				if key == "InputFile" {
					file, err = files[0].Open()
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						respondError(w, "Failed to open file, error: "+err.Error())
						return
					}
				}
			}
			if file != nil {
				defer file.Close()
			}
			convertFile(file, InputFormat, OutputFormat, processors, w)
		} else {
			convertText(InputData, InputFormat, OutputFormat, processors, w)
		}

	default:
		respondError(w, "Failed to convert data")
	}
}

func convertText(InputData string, InputFormat, OutputFormat constants.FormatCode, processors []string, w http.ResponseWriter) {
	int64Util := selectInt64Util(InputFormat)
	if int64Util != nil {
		int64ToStr := selectInt64Util(OutputFormat)
		int64Array := TextToNums[int64](InputData, int64Util)
		numsToResp := &Int64sToResp{
			NumToStr: int64ToStr,
		}
		resp := numsToResp.ToResp(int64Array)
		writeResponse(w, string(resp))
		return
	}

	byteUtil, _ := selectByteUtil(InputFormat)
	if byteUtil != nil {
		byteToStr, withDetails := selectByteUtil(OutputFormat)
		byteArray := TextToNums[byte](InputData, byteUtil)
		//process data before conversion
		byteArray, err := process(processors, byteArray)
		if err != nil {
			respondError(w, err.Error())
			return
		}

		numsToResp := &BytesToResp{
			NumToStr:       byteToStr,
			WithDetails:    withDetails,
			GlobalRowIndex: 0,
		}
		resp := numsToResp.ToResp(byteArray)
		writeResponse(w, string(resp))
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	respondError(w, fmt.Sprintln("Cannot convert ", InputFormat, " to ", OutputFormat))
}

func convertFile(file multipart.File, InputFormat, OutputFormat constants.FormatCode, processors []string, w http.ResponseWriter) {
	int64File := selectInt64File(file, InputFormat)
	if int64File != nil {
		int64ToStr := selectInt64Util(OutputFormat)
		readChan := make(chan []int64)
		go FileToNums[int64](int64File, reqInt64BufferPool, readChan)
		numsToResp := &Int64sToResp{
			NumToStr: int64ToStr,
		}
		ReadFromChannelAndRespond(readChan, numsToResp, reqInt64BufferPool, w)
		return
	}

	byteFile := selectByteFile(file, InputFormat)
	if byteFile != nil {
		byteToStr, withDetails := selectByteUtil(OutputFormat)
		readChan := make(chan []byte)
		go FileToNums[byte](byteFile, reqByteBufferPool, readChan)
		numsToResp := &BytesToResp{
			NumToStr:       byteToStr,
			WithDetails:    withDetails,
			GlobalRowIndex: 0,
		}
		ReadFromChannelAndRespond(readChan, numsToResp, reqByteBufferPool, w)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	respondError(w, fmt.Sprintln("Cannot convert ", InputFormat, " to ", OutputFormat))
}

func process(processors []string, arr []byte) ([]byte, error) {
	chain := processor.NewProcessorChain(processors)

	result, err := chain.Process(arr)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func writeResponse(w http.ResponseWriter, response string) {
	enc := json.NewEncoder(w)
	resData := resData{
		Success: true,
		Message: "Data was converted!",
		Data:    response,
	}
	err := enc.Encode(resData)
	if err != nil {
		respondError(w, "Failed to encode data, error: "+err.Error())
		return
	}
}

func selectInt64Util(format constants.FormatCode) (strToInt64 myutil.Int64Util) {
	switch format {
	case constants.Hex:
		strToInt64 = myutil.HexUtil{}
	case constants.Dec:
		strToInt64 = myutil.DecUtil{}
	case constants.Bin:
		strToInt64 = myutil.BinUtil{}
	default:
		strToInt64 = nil
	}
	return strToInt64
}

func selectInt64File(file multipart.File, format constants.FormatCode) (newFile *myfile.StrNumFile[int64]) {
	switch format {
	case constants.Hex:
		newFile = myfile.NewStrHexFile(file)
	case constants.Dec:
		newFile = myfile.NewStrDecFile(file)
	case constants.Bin:
		newFile = myfile.NewStrBinFile(file)
	default:
		newFile = nil
	}
	return newFile
}

func selectByteUtil(format constants.FormatCode) (byteToStr myutil.ByteUtil, withDetails bool) {
	withDetails = false
	switch format {
	case constants.HexByte:
		byteToStr = myutil.Hex8Util{}
	case constants.DecByte:
		byteToStr = myutil.Byte8Util{}
	case constants.DecInt8:
		byteToStr = myutil.Int8Util{}
	case constants.FormattedHexByte:
		withDetails = true
		byteToStr = myutil.Hex8Util{}
	case constants.FormattedDecByte:
		withDetails = true
		byteToStr = myutil.Byte8Util{}
	case constants.FormattedDecInt8:
		withDetails = true
		byteToStr = myutil.Int8Util{}
	case constants.RawBytes:
		byteToStr = myutil.RawBytesUtil{}
	default:
		byteToStr = nil
	}
	return byteToStr, withDetails
}

func selectByteFile(file multipart.File, format constants.FormatCode) (newFile myfile.INumFile[byte]) {
	switch format {
	case constants.HexByte:
		newFile = myfile.NewStrHex8File(file)
	case constants.DecByte:
		newFile = myfile.NewStrByte8File(file)
	case constants.DecInt8:
		newFile = myfile.NewStrInt8File(file)
	case constants.RawBytes:
		newFile = file
	default:
		newFile = nil
	}
	return newFile
}
