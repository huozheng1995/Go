package utils

import (
	"github.com/edward/scp-294/logger"
	"io"
	"mime/multipart"
	"myutil"
	myfile "myutil/file"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// text to num array

func TextToNums[T any](text string, strToNum myutil.StrToNum[T]) []T {
	result := make([]T, 0, len(text))
	var val byte
	var builder strings.Builder
	for i := 0; i < len(text)+1; i++ {
		if i == len(text) {
			val = 0
		} else {
			val = text[i]
		}
		if (val >= '0' && val <= '9') || (val >= 'a' && val <= 'f') || (val >= 'A' && val <= 'F') || val == '-' {
			builder.WriteByte(val)
		} else {
			if builder.Len() > 0 {
				result = append(result, strToNum.ToNum(builder.String()))
				builder.Reset()
			}
		}
	}

	return result
}

// file to num array

func RawBytesFileToBytes(file multipart.File, reqBufferPool *sync.Pool, readChan chan []byte) {
	rawBytesNumFile := &myfile.RawBytesNumFile{
		File: file,
	}
	fileToNums[byte](rawBytesNumFile, reqBufferPool, readChan)
}

func StrNumFileToNums[T any](file *myfile.StrNumFile[T], reqBufferPool *sync.Pool, readChan chan []T) {
	fileToNums[T](file, reqBufferPool, readChan)
}

func fileToNums[T any](file myfile.INumFile[T], reqBufferPool *sync.Pool, readChan chan []T) {
	defer close(readChan)
	for {
		buf := reqBufferPool.Get().([]T)
		if len(buf) < cap(buf) {
			logger.Logger.Log("Main", "Resize the buf from "+strconv.Itoa(len(buf))+" to "+strconv.Itoa(cap(buf)))
			buf = buf[0:cap(buf)]
		}
		len, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				logger.Logger.Log("Main", "File stream read done")
				if len > 0 {
					readChan <- buf[0:len]
					return
				}
			} else {
				logger.Logger.Log("Main", "Failed to read file stream, error: "+err.Error())
			}
			reqBufferPool.Put(buf)
			return
		}

		readChan <- buf[0:len]
	}
}

func ReadFromChannelAndRespond[T any](readChan <-chan []T, numsToResp NumsToResp[T], reqBufferPool *sync.Pool, w http.ResponseWriter) {
	readSize := 0
	writeSize := 0
	for {
		buf, ok := <-readChan
		bufLen := len(buf)
		if !ok || len(buf) <= 0 {
			logger.Logger.Log("Main", "Read channel done, total size: "+strconv.Itoa(readSize)+"("+strconv.Itoa((readSize*numsToResp.GetBytes())>>10)+"KB)")
			logger.Logger.Log("Main", "Write stream done, total size: "+strconv.Itoa(writeSize)+"("+strconv.Itoa((writeSize*numsToResp.GetBytes())>>10)+"KB)")
			return
		}

		resPageBuf := numsToResp.ToResp(buf)
		response := resPageBuf.Bytes()
		resLen := len(response)

		reqBufferPool.Put(buf)
		w.Write(response)

		readSize += bufLen
		writeSize += resLen
		//logger.Logger.Log("Main", "Read stream size: " + strconv.Itoa(readSize << 3) + "Byte")
	}
}
