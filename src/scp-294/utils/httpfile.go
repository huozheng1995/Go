package utils

import (
	"github.com/edward/scp-294/logger"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"sync"
)

func FileToRawBytes(file multipart.File, bufferPool *sync.Pool, exitChan chan struct{}) (readChan chan Page[byte]) {
	readChan = make(chan Page[byte])
	go func() {
		defer close(readChan)
		pageNum := 1
		for {
			select {
			case <-exitChan:
				logger.Log("Exit channel is closed")
				return
			default:
				pageBuffer := bufferPool.Get().([]byte)
				var page Page[byte]
				page = Page[byte]{
					pageNum:  pageNum,
					buffer:   pageBuffer,
					pageSize: cap(pageBuffer),
					index:    0,
				}
				pageNum++

				var err error
				page.index, err = file.Read(pageBuffer)
				if err != nil {
					if err != io.EOF {
						logger.Log("Failed to read file stream, error: " + err.Error())
						bufferPool.Put(pageBuffer)
						return
					} else {
						logger.Log("File stream read done")
					}
				}

				//If you're just reading data from the buffer in readStreamAndSendBody without modifying it,
				//you can use readChan <- buffer[:n] without any issues.
				//There's no need to use append([]byte(nil), buffer[:n]...) in this case.
				//In fact, using readChan <- buffer[:n] is more efficient because it avoids creating a new slice,
				//which would incur extra memory usage and GC overhead.
				readChan <- page
			}
		}
	}()
	return
}

func FileToPageBuffer[T any](file multipart.File, bufferPool *sync.Pool, funcStrToNum func(string) T, exitChan chan struct{}) (readChan chan Page[T]) {
	readChan = make(chan Page[T])
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
				var page Page[T]
				page = CreateEmptyPage(pageNum, pageBuffer, funcStrToNum)
				pageNum++

				var err error
				err = page.AppendData(&tempBuffer, file)
				if err != nil {
					if err != io.EOF {
						logger.Log("Failed to read file stream, error: " + err.Error())
						bufferPool.Put(pageBuffer)
						return
					} else {
						logger.Log("File stream read done")
					}
				}
				readChan <- page
			}
		}
	}()
	return
}

func ReadBytesAndResponse(readChan <-chan Page[byte], funcByteToStr ByteToStr, withDetails bool, bufferPool *sync.Pool, w http.ResponseWriter) {
	readSize := 0
	writeSize := 0
	globalRowIndex := 0
	for {
		var page Page[byte]
		var buffer []byte
		page, ok := <-readChan
		buffer = page.buffer
		if !ok || page.index <= 0 {
			logger.Log("Read channel done, total size: " + strconv.Itoa(readSize) + "Byte")
			logger.Log("Write stream done, total size: " + strconv.Itoa(writeSize) + "Byte")
			return
		}
		rowsBytes := ByteArrayToOutput(buffer[0:page.index], &globalRowIndex, funcByteToStr, withDetails)

		bufferPool.Put(buffer)
		w.Write(rowsBytes)
		readSize += page.index
		writeSize += len(rowsBytes)
		//logger.Log("Read stream size: " + strconv.Itoa(readSize) + "Byte")
	}
}

func ReadInt64ArrayAndResponse(readChan <-chan Page[int64], funcInt64ToStr Int64ToStr, bufferPool *sync.Pool, w http.ResponseWriter) {
	readSize := 0
	writeSize := 0
	for {
		var page Page[int64]
		var buffer []int64
		page, ok := <-readChan
		buffer = page.buffer
		if !ok || len(buffer) <= 0 {
			logger.Log("Read channel done, total size: " + strconv.Itoa(readSize) + "Byte")
			logger.Log("Write stream done, total size: " + strconv.Itoa(writeSize) + "Byte")
			return
		}
		rowsBytes := Int64ArrayToOutput(buffer[0:page.index], funcInt64ToStr)

		bufferPool.Put(buffer)
		w.Write(rowsBytes)
		readSize += page.index
		writeSize += len(rowsBytes)
		//logger.Log("Read stream size: " + strconv.Itoa(readSize) + "Int64")
	}
}
