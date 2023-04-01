package utils

import (
	"github.com/edward/scp-294/logger"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"sync"
)

func FileToRawBytes(file multipart.File, bufferPool *sync.Pool, exitChan chan struct{}) (readChan chan []byte) {
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

func FileToPageBuffer[T any](file multipart.File, bufferPool *sync.Pool, funcStrToNum func(string) T, exitChan chan struct{}) (readChan chan []T) {
	readChan = make(chan []T)
	go func() {
		defer close(readChan)
		pageNum := 1
		tempBuffer := CreateTempBuffer()
		for {
			select {
			case <-exitChan:
				logger.Log("Exit channel is closed")
				return
			default:
				pageBuffer := bufferPool.Get().([]T)
				page := CreateEmptyPage(pageNum, pageBuffer, funcStrToNum)
				pageNum++
				err := page.AppendData(&tempBuffer, file)
				if err != nil {
					if err != io.EOF {
						logger.Log("Failed to read file stream, error: " + err.Error())
						bufferPool.Put(pageBuffer)
					} else {
						logger.Log("File stream read done")
						if !page.IsEmpty() {
							readChan <- page.GetBuffer()
						} else {
							bufferPool.Put(pageBuffer)
						}
					}
					return
				}
				readChan <- page.GetBuffer()
			}
		}
	}()
	return
}

func ReadBytesAndResponse(readChan <-chan []byte, funcByteToStr ByteToStr, withDetails bool, bufferPool *sync.Pool, w http.ResponseWriter) {
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
		rowsBytes := ByteArrayToOutput(buffer, &globalRowIndex, funcByteToStr, withDetails)

		bufferPool.Put(buffer)
		w.Write(rowsBytes)
		readSize += len(buffer)
		writeSize += len(rowsBytes)
		//logger.Log("Read stream size: " + strconv.Itoa(readSize) + "Byte")
	}
}

func ReadInt64ArrayAndResponse(readChan <-chan []int64, funcInt64ToStr Int64ToStr, bufferPool *sync.Pool, w http.ResponseWriter) {
	readSize := 0
	writeSize := 0
	for {
		buffer, ok := <-readChan
		if !ok || len(buffer) <= 0 {
			logger.Log("Read channel done, total size: " + strconv.Itoa(readSize) + "Byte")
			logger.Log("Write stream done, total size: " + strconv.Itoa(writeSize) + "Byte")
			return
		}
		rowsBytes := Int64ArrayToOutput(buffer, funcInt64ToStr)

		bufferPool.Put(buffer)
		w.Write(rowsBytes)
		readSize += len(buffer)
		writeSize += len(rowsBytes)
		//logger.Log("Read stream size: " + strconv.Itoa(readSize) + "Int64")
	}
}
