package controller

import (
	"encoding/json"
	"github.com/edward/scp-294/common"
	"github.com/edward/scp-294/logger"
	"github.com/edward/scp-294/utils"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
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

var pageBufferCount int32
var pageBufferPool = &sync.Pool{
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

		var InputType, InputFormat, OutputType, OutputFormat common.NumType
		var InputData string
		for key, values := range form.Value {
			if key == "InputType" {
				intVal, _ := strconv.Atoi(values[0])
				InputType = common.NumType(intVal)
			} else if key == "InputFormat" {
				intVal, _ := strconv.Atoi(values[0])
				InputFormat = common.NumType(intVal)
			} else if key == "OutputType" {
				intVal, _ := strconv.Atoi(values[0])
				OutputType = common.NumType(intVal)
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
			convertFile(w, file, OutputType, InputFormat, OutputFormat)
		} else {
			convertText(w, InputData, OutputType, InputFormat, OutputFormat)
		}

	default:
		common.ResponseError(w, "Failed to convert data")
	}
}

func convertText(w http.ResponseWriter, InputData string, OutputType, InputFormat, OutputFormat common.NumType) {
	var outputData string
	var strings = utils.SplitInputString(InputData)
	switch InputFormat {
	case common.Hex:
		switch OutputFormat {
		case common.Hex:
			outputData = utils.HexArrayToString(strings)
		case common.Dec:
			outputData = utils.DecArrayToString(utils.HexArrayToDecArray(strings))
		case common.Bin:
			outputData = utils.BinArrayToString(utils.HexArrayToBinArray(strings))
		default:
			common.ResponseError(w, "Cannot convert '"+common.InputFormatMap[InputFormat]+"' to '"+common.OutputFormatMap[OutputFormat]+"'")
			return
		}
	case common.Dec:
		switch OutputFormat {
		case common.Hex:
			outputData = utils.HexArrayToString(utils.DecArrayToHexArray(strings))
		case common.Dec:
			outputData = utils.DecArrayToString(utils.DecArrayToDecArray(strings))
		case common.Bin:
			outputData = utils.BinArrayToString(utils.DecArrayToBinArray(strings))
		default:
			common.ResponseError(w, "Cannot convert '"+common.InputFormatMap[InputFormat]+"' to '"+common.OutputFormatMap[OutputFormat]+"'")
			return
		}
	case common.Bin:
		switch OutputFormat {
		case common.Hex:
			outputData = utils.HexArrayToString(utils.BinArrayToHexArray(strings))
		case common.Dec:
			outputData = utils.DecArrayToString(utils.BinArrayToDecArray(strings))
		case common.Bin:
			outputData = utils.BinArrayToString(strings)
		default:
			common.ResponseError(w, "Cannot convert '"+common.InputFormatMap[InputFormat]+"' to '"+common.OutputFormatMap[OutputFormat]+"'")
			return
		}
	case common.HexByte:
		switch OutputFormat {
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
			common.ResponseError(w, "Cannot convert '"+common.InputFormatMap[InputFormat]+"' to '"+common.OutputFormatMap[OutputFormat]+"'")
			return
		}
	case common.DecByte:
		switch OutputFormat {
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
			common.ResponseError(w, "Cannot convert '"+common.InputFormatMap[InputFormat]+"' to '"+common.OutputFormatMap[OutputFormat]+"'")
			return
		}
	case common.DecInt8:
		switch OutputFormat {
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
			common.ResponseError(w, "Cannot convert '"+common.InputFormatMap[InputFormat]+"' to '"+common.OutputFormatMap[OutputFormat]+"'")
			return
		}
	default:
		common.ResponseError(w, "Unknown input format: '"+common.InputFormatMap[InputFormat]+"'")
		return
	}

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

func convertFile(w http.ResponseWriter, file multipart.File, OutputType, InputFormat, OutputFormat common.NumType) {
	var funcStrToDecByte utils.StrToDecByte
	switch InputFormat {
	case common.RawBytes:
	case common.Hex:
	case common.Dec:
	case common.Bin:
	case common.HexByte:
		funcStrToDecByte = utils.HexByteToDecByte
	case common.DecByte:
		funcStrToDecByte = utils.DecByteToDecByte
	case common.DecInt8:
		funcStrToDecByte = utils.DecInt8ToDecByte
	default:
		w.WriteHeader(http.StatusInternalServerError)
		common.ResponseError(w, "Cannot convert '"+common.InputFormatMap[InputFormat]+"' to '"+common.OutputFormatMap[OutputFormat]+"'")
		return
	}

	var funcBytesToRow utils.ByteArrayToRow
	var withDetails = false
	switch OutputFormat {
	case common.HexByte:
		funcBytesToRow = utils.ByteArrayToHexByteRow
	case common.DecByte:
		funcBytesToRow = utils.ByteArrayToByteRow
	case common.DecInt8:
		funcBytesToRow = utils.ByteArrayToInt8Row
	case common.HexByteFormatted:
		withDetails = true
		funcBytesToRow = utils.ByteArrayToHexByteRow
	case common.DecByteFormatted:
		withDetails = true
		funcBytesToRow = utils.ByteArrayToByteRow
	case common.DecInt8Formatted:
		withDetails = true
		funcBytesToRow = utils.ByteArrayToInt8Row
	default:
		w.WriteHeader(http.StatusInternalServerError)
		common.ResponseError(w, "Cannot convert '"+common.InputFormatMap[InputFormat]+"' to '"+common.OutputFormatMap[OutputFormat]+"'")
		return
	}

	//logger.Log("Begin parse, buffer count: " + strconv.Itoa(int(fileBufferCount)))
	//exitChannel, readChan := fileStreamToRawBytes(file, fileBufferPool)
	//readStreamAndSendBody(w, readChan, funcBytesToRow, withDetails, fileBufferPool)
	//logger.Log("End parse, buffer count: " + strconv.Itoa(int(fileBufferCount)))
	//close(exitChannel)

	logger.Log("Begin parse, buffer count: " + strconv.Itoa(int(pageBufferCount)))
	exitChannel, readChan := fileStreamToPage(file, pageBufferPool, funcStrToDecByte)
	readStreamAndSendBody(w, readChan, funcBytesToRow, withDetails, fileBufferPool)
	logger.Log("End parse, buffer count: " + strconv.Itoa(int(pageBufferCount)))
	close(exitChannel)
}

func fileStreamToRawBytes(file multipart.File, bufferPool *sync.Pool) (exitChan chan struct{}, readChan chan []byte) {
	exitChan = make(chan struct{})
	readChan = make(chan []byte)
	go func() {
		defer file.Close()
		defer close(readChan)
		for {
			select {
			case <-exitChan:
				logger.Log("Exit channel is closed")
				return
			default:
				buffer := bufferPool.Get().([]byte)
				n, err := file.Read(buffer)
				if err != nil {
					if err != io.EOF {
						logger.Log("Failed to read file stream, error: " + err.Error())
					} else {
						logger.Log("File stream read done")
					}
					return
				}
				readChan <- buffer[:n]
			}
		}
	}()
	return
}

func readStreamAndSendBody(w http.ResponseWriter, readChan <-chan []byte, funcBytesToRow utils.ByteArrayToRow, withDetails bool, bufferPool *sync.Pool) {
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
		rowsBytes := utils.StreamBytesToRowsBytes(data, &globalRowIndex, funcBytesToRow, withDetails)
		bufferPool.Put(data)
		w.Write(rowsBytes)
		readSize += len(data)
		writeSize += len(rowsBytes)
		//logger.Log("Read stream size: " + strconv.Itoa(readSize) + "Byte")
	}
}

func fileStreamToPage(file multipart.File, bufferPool *sync.Pool, funcStrToDecByte utils.StrToDecByte) (exitChan chan struct{}, readChan chan []byte) {
	exitChan = make(chan struct{})
	readChan = make(chan []byte)
	go func() {
		defer file.Close()
		defer close(readChan)
		pageNum := 1
		preBuffer := make([]byte, 1024)
		preOff := 0
		preLen := 0
		var preCell strings.Builder
		for {
			select {
			case <-exitChan:
				logger.Log("Exit channel is closed")
				return
			default:
				pageBuffer := bufferPool.Get().([]byte)

				page := utils.CreateEmptyPage(pageNum, pageBuffer)
				pageNum++
				preBuffer, preOff, preLen, preCell = utils.FillPage(&page, preBuffer, preOff, preLen, preCell, file, funcStrToDecByte)

				readChan <- page.Buffer[:page.Index]
			}
		}
	}()
	return
}
