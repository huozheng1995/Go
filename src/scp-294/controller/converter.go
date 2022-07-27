package controller

import (
	"encoding/json"
	"fmt"
	"github.com/edward/scp-294/common"
	"github.com/edward/scp-294/converter"
	"github.com/edward/scp-294/logger"
	"io"
	"mime/multipart"
	"net/http"
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
		if InputType == "File" {
			var file multipart.File
			var fileName string
			for key, files := range form.File {
				if key == "InputFile" {
					fileName = files[0].Filename
					file, _ = files[0].Open()
				}
			}
			exitChan, dataChan := receiveFile(file)
			sendBody(w, dataChan)
			println(fileName)
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
				outputData = converter.HexByteArrayToString(strings)
			case "ByteArray":
				outputData = converter.ByteArrayToString(converter.HexByteArrayToDecByteArray(strings))
			case "Int8Array":
				outputData = converter.Int8ArrayToString(converter.HexByteArrayToInt8Array(strings))
			default:
				common.ResponseError(w, "Cannot convert '"+InputType+"' to '"+OutputType+"'")
				return
			}
		case "ByteArray":
			switch OutputType {
			case "HexByteArray":
				outputData = converter.HexByteArrayToString(converter.DecByteArrayToHexByteArray(strings))
			case "ByteArray":
				outputData = converter.ByteArrayToString(converter.DecByteArrayToDecByteArray(strings))
			case "Int8Array":
				outputData = converter.Int8ArrayToString(converter.DecByteArrayToInt8Array(strings))
			default:
				common.ResponseError(w, "Cannot convert '"+InputType+"' to '"+OutputType+"'")
				return
			}
		case "Int8Array":
			switch OutputType {
			case "HexByteArray":
				outputData = converter.HexByteArrayToString(converter.Int8ArrayToHexByteArray(strings))
			case "ByteArray":
				outputData = converter.ByteArrayToString(converter.Int8ArrayToDecByteArray(strings))
			case "Int8Array":
				outputData = converter.Int8ArrayToString(converter.Int8ArrayToInt8Array(strings))
			default:
				common.ResponseError(w, "Cannot convert '"+InputType+"' to '"+OutputType+"'")
				return
			}
		default:
			common.ResponseError(w, "Unknown input type: '"+InputType+"'")
			return
		}

		enc := json.NewEncoder(w)
		convertRes := common.ConvertRes{
			OutputData: outputData,
		}
		resData := common.ResData{
			Success: true,
			Message: "Data was converted!",
			Data:    convertRes,
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

func receiveFile(file multipart.File) (exitChan chan struct{}, dataChan chan []byte) {
	exitChan = make(chan struct{})
	dataChan = make(chan []byte, 10)
	go func() {
		defer close(dataChan)
		defer file.Close()
		for {
			select {
			case <-exitChan:
				fmt.Println("exit receive data!")
				return
			default:
				buf := make([]byte, 4096)
				n, err := file.Read(buf)
				if err != nil && err != io.EOF {
					fmt.Println("receive data error!")
					return
				} else {
					dataChan <- buf[:n]
				}
			}
		}
	}()
	return
}

func sendBody(w http.ResponseWriter, dataChan <-chan []byte) {
	for {
		data, ok := <-dataChan
		if !ok {
			fmt.Println("receive data error!")
			break
		}
		n, err := w.Write(data)
		if err != nil || n == 0 {
			fmt.Println("send data error!")
			break
		}
	}
}

func convertOld(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		dec := json.NewDecoder(r.Body)
		convertReq := common.ConvertReq{}
		err := dec.Decode(&convertReq)
		if err != nil {
			logger.Log(err.Error())
			common.ResponseError(w, "Failed to decode data, error: "+err.Error())
			return
		}

		var outputData string
		var strings = converter.SplitInputString(convertReq.InputData)
		switch convertReq.InputType {
		case "Hex":
			switch convertReq.OutputType {
			case "Hex":
				outputData = converter.HexArrayToString(strings)
			case "Dec":
				outputData = converter.DecArrayToString(converter.HexArrayToDecArray(strings))
			case "Bin":
				outputData = converter.BinArrayToString(converter.HexArrayToBinArray(strings))
			default:
				common.ResponseError(w, "Cannot convert '"+convertReq.InputType+"' to '"+convertReq.OutputType+"'")
				return
			}
		case "Dec":
			switch convertReq.OutputType {
			case "Hex":
				outputData = converter.HexArrayToString(converter.DecArrayToHexArray(strings))
			case "Dec":
				outputData = converter.DecArrayToString(converter.DecArrayToDecArray(strings))
			case "Bin":
				outputData = converter.BinArrayToString(converter.DecArrayToBinArray(strings))
			default:
				common.ResponseError(w, "Cannot convert '"+convertReq.InputType+"' to '"+convertReq.OutputType+"'")
				return
			}
		case "Bin":
			switch convertReq.OutputType {
			case "Hex":
				outputData = converter.HexArrayToString(converter.BinArrayToHexArray(strings))
			case "Dec":
				outputData = converter.DecArrayToString(converter.BinArrayToDecArray(strings))
			case "Bin":
				outputData = converter.BinArrayToString(strings)
			default:
				common.ResponseError(w, "Cannot convert '"+convertReq.InputType+"' to '"+convertReq.OutputType+"'")
				return
			}
		case "HexByteArray":
			switch convertReq.OutputType {
			case "HexByteArray":
				outputData = converter.HexByteArrayToString(strings)
			case "ByteArray":
				outputData = converter.ByteArrayToString(converter.HexByteArrayToDecByteArray(strings))
			case "Int8Array":
				outputData = converter.Int8ArrayToString(converter.HexByteArrayToInt8Array(strings))
			default:
				common.ResponseError(w, "Cannot convert '"+convertReq.InputType+"' to '"+convertReq.OutputType+"'")
				return
			}
		case "ByteArray":
			switch convertReq.OutputType {
			case "HexByteArray":
				outputData = converter.HexByteArrayToString(converter.DecByteArrayToHexByteArray(strings))
			case "ByteArray":
				outputData = converter.ByteArrayToString(converter.DecByteArrayToDecByteArray(strings))
			case "Int8Array":
				outputData = converter.Int8ArrayToString(converter.DecByteArrayToInt8Array(strings))
			default:
				common.ResponseError(w, "Cannot convert '"+convertReq.InputType+"' to '"+convertReq.OutputType+"'")
				return
			}
		case "Int8Array":
			switch convertReq.OutputType {
			case "HexByteArray":
				outputData = converter.HexByteArrayToString(converter.Int8ArrayToHexByteArray(strings))
			case "ByteArray":
				outputData = converter.ByteArrayToString(converter.Int8ArrayToDecByteArray(strings))
			case "Int8Array":
				outputData = converter.Int8ArrayToString(converter.Int8ArrayToInt8Array(strings))
			default:
				common.ResponseError(w, "Cannot convert '"+convertReq.InputType+"' to '"+convertReq.OutputType+"'")
				return
			}
		default:
			common.ResponseError(w, "Unknown input type: '"+convertReq.InputType+"'")
			return
		}

		enc := json.NewEncoder(w)
		convertRes := common.ConvertRes{
			OutputData: outputData,
		}
		resData := common.ResData{
			Success: true,
			Message: "Data was converted!",
			Data:    convertRes,
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
