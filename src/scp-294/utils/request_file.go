package utils

import (
	"github.com/edward/scp-294/logger"
	"io"
	"mime/multipart"
	"myutil"
	myfile "myutil/file"
	"net/http"
	"strconv"
	"sync"
)

func RawBytesFileToBytes(file multipart.File, reqBufferPool *sync.Pool, readChan chan []byte) {
	rawBytesNumFile := &myfile.RawBytesNumFile[byte]{
		File: file,
	}
	fileToNum[byte](rawBytesNumFile, reqBufferPool, readChan)
}

func StrNumFileToNums[T any](strNumFile *myfile.StrNumFile[T], reqBufferPool *sync.Pool, readChan chan []T) {
	fileToNum[T](strNumFile, reqBufferPool, readChan)
}

func fileToNum[T any](file myfile.INumFile[T], reqBufferPool *sync.Pool, readChan chan []T) {
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

func ByteArrayChannelToResponse(readChan <-chan []byte, funcByteToStr myutil.ByteToStr, withDetails bool, reqBufferPool *sync.Pool, w http.ResponseWriter) {
	readSize := 0
	writeSize := 0
	globalRowIndex := 0
	for {
		buf, ok := <-readChan
		bufLen := len(buf)
		if !ok || len(buf) <= 0 {
			logger.Logger.Log("Main", "Read channel done, total size: "+strconv.Itoa(readSize)+"Byte("+strconv.Itoa(readSize>>10)+"KB)")
			logger.Logger.Log("Main", "Write stream done, total size: "+strconv.Itoa(writeSize)+"Byte("+strconv.Itoa(writeSize>>10)+"KB)")
			return
		}

		resPageBuf := ByteArrayToResponse(buf, &globalRowIndex, funcByteToStr, withDetails)
		response := resPageBuf.Bytes()
		resLen := len(response)

		reqBufferPool.Put(buf)
		w.Write(response)

		readSize += bufLen
		writeSize += resLen
		//logger.Logger.Log("Main", "Read stream size: " + strconv.Itoa(readSize) + "Byte")
	}
}

func Int64ArrayChannelToResponse(readChan <-chan []int64, funcInt64ToStr myutil.Int64ToStr, reqBufferPool *sync.Pool, w http.ResponseWriter) {
	readSize := 0
	writeSize := 0
	for {
		buf, ok := <-readChan
		bufLen := len(buf)
		if !ok || len(buf) <= 0 {
			logger.Logger.Log("Main", "Read channel done, total size: "+strconv.Itoa(readSize<<3)+"Byte")
			logger.Logger.Log("Main", "Write stream done, total size: "+strconv.Itoa(writeSize<<3)+"Byte")
			return
		}

		resPageBuf := Int64ArrayToResponse(buf, funcInt64ToStr)
		response := resPageBuf.Bytes()
		resLen := len(response)

		reqBufferPool.Put(buf)
		w.Write(response)

		readSize += bufLen
		writeSize += resLen
		//logger.Logger.Log("Main", "Read stream size: " + strconv.Itoa(readSize << 3) + "Byte")
	}
}
