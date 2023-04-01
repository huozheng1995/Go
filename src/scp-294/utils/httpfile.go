package utils

import (
	"github.com/edward/scp-294/logger"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

func FileStreamToRawBytes(file multipart.File, bufferPool *sync.Pool) (exitChan chan struct{}, readChan chan []byte) {
	exitChan = make(chan struct{})
	readChan = make(chan []byte)
	go func() {
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
					bufferPool.Put(buffer)
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

func FileStreamToPageBytes(file multipart.File, bufferPool *sync.Pool, funcStrToByte StrToByte) (exitChan chan struct{}, readChan chan []byte) {
	exitChan = make(chan struct{})
	readChan = make(chan []byte)
	go func() {
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

				page := CreateEmptyPage(pageNum, pageBuffer)
				pageNum++
				var err error
				err, preBuffer, preOff, preLen, tempCell = FillPage(&page, preBuffer, preOff, preLen, tempCell, file, funcStrToByte)
				if err != nil {
					if err != io.EOF {
						logger.Log("Failed to read file stream, error: " + err.Error())
						bufferPool.Put(pageBuffer)
					} else {
						logger.Log("File stream read done")
						if page.Index > 0 {
							readChan <- page.Buffer[:page.Index]
						} else {
							bufferPool.Put(pageBuffer)
						}
					}
					return
				}
				readChan <- page.Buffer[:page.Index]
			}
		}
	}()
	return
}

func ReadBytesAndSendBody(w http.ResponseWriter, readChan <-chan []byte, funcByteToStr ByteToStr, withDetails bool, bufferPool *sync.Pool) {
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
		rowsBytes := ByteArrayToRowBytes(buffer, &globalRowIndex, funcByteToStr, withDetails)

		bufferPool.Put(buffer)
		w.Write(rowsBytes)
		readSize += len(buffer)
		writeSize += len(rowsBytes)
		//logger.Log("Read stream size: " + strconv.Itoa(readSize) + "Byte")
	}
}

func FileStreamToInt64Array(file multipart.File, bufferPool *sync.Pool, funcStrToInt64 StrToInt64) (exitChan chan struct{}, readChan chan []int64) {
	exitChan = make(chan struct{})
	readChan = make(chan []int64)
	go func() {
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
				pageBuffer := bufferPool.Get().([]int64)
				page := CreateEmptyInt64Page(pageNum, pageBuffer)
				pageNum++
				var err error
				err, preBuffer, preOff, preLen, tempCell = FillInt64Page(&page, preBuffer, preOff, preLen, tempCell, file, funcStrToInt64)
				if err != nil {
					if err != io.EOF {
						logger.Log("Failed to read file stream, error: " + err.Error())
						bufferPool.Put(pageBuffer)
					} else {
						logger.Log("File stream read done")
						if page.Index > 0 {
							readChan <- page.Buffer[:page.Index]
						} else {
							bufferPool.Put(pageBuffer)
						}
					}
					return
				}
				readChan <- page.Buffer[:page.Index]
			}
		}
	}()
	return
}

func ReadInt64ArrayAndSendBody(w http.ResponseWriter, readChan <-chan []int64, funcInt64ToStr Int64ToStr, bufferPool *sync.Pool) {
	readSize := 0
	writeSize := 0
	for {
		buffer, ok := <-readChan
		if !ok || len(buffer) <= 0 {
			logger.Log("Read channel done, total size: " + strconv.Itoa(readSize) + "Byte")
			logger.Log("Write stream done, total size: " + strconv.Itoa(writeSize) + "Byte")
			return
		}
		rowsBytes := Int64ArrayToRowBytes(buffer, funcInt64ToStr)

		bufferPool.Put(buffer)
		w.Write(rowsBytes)
		readSize += len(buffer)
		writeSize += len(rowsBytes)
		//logger.Log("Read stream size: " + strconv.Itoa(readSize) + "Int64")
	}
}
