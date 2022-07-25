package controller

import (
	"encoding/json"
	"github.com/edward/scp-294/common"
	"github.com/edward/scp-294/converter"
	"github.com/edward/scp-294/logger"
	"net/http"
)

func registerConverterRoutes() {
	http.HandleFunc("/convert", convert)
}

func convert(w http.ResponseWriter, r *http.Request) {
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
