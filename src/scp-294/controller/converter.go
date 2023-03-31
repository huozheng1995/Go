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

func convertFile(w http.ResponseWriter, file multipart.File, InputFormat, OutputFormat common.NumType) {
	rawBytes := false
	var funcStrToInt64 utils.StrToInt64
	var funcStrToByte utils.StrToByte
	switch InputFormat {
	case common.RawBytes:
		rawBytes = true
	case common.Hex:
		funcStrToInt64 = utils.HexStrToInt64
	case common.Dec:
		funcStrToInt64 = utils.Int64StrToInt64
	case common.Bin:
		funcStrToInt64 = utils.BinStrToInt64
	case common.HexByte:
		funcStrToByte = utils.HexStrToByte
	case common.DecByte:
		funcStrToByte = utils.ByteStrToByte
	case common.DecInt8:
		funcStrToByte = utils.Int8StrToByte
	default:
		w.WriteHeader(http.StatusInternalServerError)
		common.ResponseError(w, "Unknown input format: '"+common.InputFormatMap[InputFormat]+"'")
		return
	}

	var funcInt64ToStr utils.Int64ToStr
	var funcInt64ArrayToRow utils.Int64ArrayToRow
	var funcByteToStr utils.ByteToStr
	var funcBytesToRow utils.ByteArrayToRow
	var withDetails = false
	if funcStrToInt64 != nil {
		switch OutputFormat {
		case common.Hex:
			funcInt64ToStr = utils.Int64ToHexStr
			funcBytesToRow = utils.Int64ArrayToHexRow
		case common.Dec:
			funcInt64ToStr = utils.Int64ToInt64Str
			funcBytesToRow = utils.Int64ArrayToDecRow
		case common.Bin:
			funcInt64ToStr = utils.Int64ToBinStr
			funcBytesToRow = utils.Int64ArrayToBinRow
		default:
			w.WriteHeader(http.StatusInternalServerError)
			common.ResponseError(w, "Cannot convert '"+common.InputFormatMap[InputFormat]+"' to '"+common.OutputFormatMap[OutputFormat]+"'")
			return
		}
	} else if rawBytes || (funcStrToByte != nil) {
		switch OutputFormat {
		case common.HexByte:
			funcByteToStr = utils.ByteToHexStr
			funcBytesToRow = utils.ByteArrayToHexByteRow
		case common.DecByte:
			funcByteToStr = utils.ByteToByteStr
			funcBytesToRow = utils.ByteArrayToByteRow
		case common.DecInt8:
			funcByteToStr = utils.ByteToInt8Str
			funcBytesToRow = utils.ByteArrayToInt8Row
		case common.HexByteFormatted:
			withDetails = true
			funcByteToStr = utils.ByteToHexStr
			funcBytesToRow = utils.ByteArrayToHexByteRow
		case common.DecByteFormatted:
			withDetails = true
			funcByteToStr = utils.ByteToByteStr
			funcBytesToRow = utils.ByteArrayToByteRow
		case common.DecInt8Formatted:
			withDetails = true
			funcByteToStr = utils.ByteToInt8Str
			funcBytesToRow = utils.ByteArrayToInt8Row
		default:
			w.WriteHeader(http.StatusInternalServerError)
			common.ResponseError(w, "Cannot convert '"+common.InputFormatMap[InputFormat]+"' to '"+common.OutputFormatMap[OutputFormat]+"'")
			return
		}
	} else {
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
	exitChannel, readChan := fileStreamToPage(file, pageBufferPool, funcStrToByte)
	readStreamAndSendBody(w, readChan, funcBytesToRow, withDetails, pageBufferPool)
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

				//If you're just reading data from the buffer in readStreamAndSendBody without modifying it,
				//you can use readChan <- buffer[:n] without any issues.
				//There's no need to use append([]byte(nil), buffer[:n]...) in this case.
				//In fact, using readChan <- buffer[:n] is more efficient because it avoids creating a new slice,
				//which would incur extra memory usage and GC overhead.
				readChan <- buffer[:n]
			}
		}
	}()
	return
}

func fileStreamToPage(file multipart.File, bufferPool *sync.Pool, funcStrToDecByte utils.StrToByte) (exitChan chan struct{}, readChan chan []byte) {
	exitChan = make(chan struct{})
	readChan = make(chan []byte)
	go func() {
		defer file.Close()
		defer close(readChan)
		pageNum := 1
		preBuffer := make([]byte, 1024)
		preOff := 0
		preLen := 0
		var tempCell strings.Builder
		for {
			select {
			case <-exitChan:
				logger.Log("Exit channel is closed")
				return
			default:
				pageBuffer := bufferPool.Get().([]byte)

				page := utils.CreateEmptyPage(pageNum, pageBuffer)
				pageNum++
				var err error
				err, preBuffer, preOff, preLen, tempCell = utils.FillPage(&page, preBuffer, preOff, preLen, tempCell, file, funcStrToDecByte)
				if err != nil {
					if err != io.EOF {
						logger.Log("Failed to read file stream, error: " + err.Error())
					} else {
						logger.Log("File stream read done")
					}
					return
				}
				readChan <- page.Buffer[:page.Index]
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
		buffer, ok := <-readChan
		if !ok || len(buffer) <= 0 {
			logger.Log("Read channel done, total size: " + strconv.Itoa(readSize) + "Byte")
			logger.Log("Write stream done, total size: " + strconv.Itoa(writeSize) + "Byte")
			return
		}
		rowsBytes := utils.StreamBytesToRowsBytes(buffer, &globalRowIndex, funcBytesToRow, withDetails)
		bufferPool.Put(buffer)
		w.Write(rowsBytes)
		readSize += len(buffer)
		writeSize += len(rowsBytes)
		//logger.Log("Read stream size: " + strconv.Itoa(readSize) + "Byte")
	}
}
