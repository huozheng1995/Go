package utils

import (
	"bytes"
	"github.com/edward/scp-294/logger"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
)

func FileToRawBytes(file multipart.File, reqBufferPool *sync.Pool, exitChan chan struct{}) (readChan chan *Page[byte]) {
	readChan = make(chan *Page[byte])
	go func() {
		defer close(readChan)
		pageNum := 1
		for {
			select {
			case <-exitChan:
				logger.Log("Exit channel is closed")
				return
			default:
				pageBuf := reqBufferPool.Get().([]byte)
				var page *Page[byte]
				page = &Page[byte]{
					pageNum:  pageNum,
					buffer:   &pageBuf,
					pageSize: cap(pageBuf),
					length:   0,
				}
				pageNum++

				var err error
				page.length, err = file.Read(pageBuf)
				if err != nil {
					if err != io.EOF {
						logger.Log("Failed to read file stream, error: " + err.Error())
					} else {
						logger.Log("File stream read done")
					}
					reqBufferPool.Put(pageBuf)
					return
				}

				readChan <- page
			}
		}
	}()
	return
}

func FileToPageBuffer[T any](file multipart.File, reqBufferPool *sync.Pool, funcStrToNum func(string) T, exitChan chan struct{}) (readChan chan *Page[T]) {
	readChan = make(chan *Page[T])
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
				pageBuf := reqBufferPool.Get().([]T)
				var page *Page[T]
				page = CreateEmptyPage(pageNum, &pageBuf, funcStrToNum)
				pageNum++

				var err error
				err = page.AppendData(&tempBuffer, file)
				if err != nil {
					if err != io.EOF {
						logger.Log("Failed to read file stream, error: " + err.Error())
					} else {
						logger.Log("File stream read done")
					}
					reqBufferPool.Put(pageBuf)
					return
				}

				readChan <- page
			}
		}
	}()
	return
}

func ReadByteArrayAndResponse(readChan <-chan *Page[byte], funcByteToStr ByteToStr, withDetails bool, reqBufferPool *sync.Pool, w http.ResponseWriter) {
	readSize := 0
	writeSize := 0
	globalRowIndex := 0
	for {
		var page *Page[byte]
		page, ok := <-readChan
		if !ok || page.length <= 0 {
			logger.Log("Read channel done, total size: " + strconv.Itoa(readSize) + "Byte(" + strconv.Itoa(readSize>>10) + "KB)")
			logger.Log("Write stream done, total size: " + strconv.Itoa(writeSize) + "Byte(" + strconv.Itoa(writeSize>>10) + "KB)")
			return
		}
		pageBuf := *page.buffer
		resPageBuf := ResBufferPool.Get().(*bytes.Buffer)
		resPageBuf.Reset()
		ByteArrayToResponse(pageBuf[0:page.length], &globalRowIndex, funcByteToStr, withDetails, resPageBuf)
		response := resPageBuf.Bytes()
		resLen := len(response)

		reqBufferPool.Put(pageBuf)
		readSize += page.length
		w.Write(response)
		ResBufferPool.Put(resPageBuf)
		writeSize += resLen
		//logger.Log("Read stream size: " + strconv.Itoa(readSize) + "Byte")
	}
}

func ReadInt64ArrayAndResponse(readChan <-chan *Page[int64], funcInt64ToStr Int64ToStr, reqBufferPool *sync.Pool, w http.ResponseWriter) {
	readSize := 0
	writeSize := 0
	for {
		var page *Page[int64]
		page, ok := <-readChan
		if !ok || page.length <= 0 {
			logger.Log("Read channel done, total size: " + strconv.Itoa(readSize<<3) + "Byte")
			logger.Log("Write stream done, total size: " + strconv.Itoa(writeSize<<3) + "Byte")
			return
		}
		pageBuf := *page.buffer
		resPageBuf := ResBufferPool.Get().(*bytes.Buffer)
		resPageBuf.Reset()
		Int64ArrayToResponse(pageBuf[0:page.length], funcInt64ToStr, resPageBuf)
		response := resPageBuf.Bytes()
		resLen := len(response)

		reqBufferPool.Put(pageBuf)
		readSize += page.length
		w.Write(response)
		ResBufferPool.Put(resPageBuf)
		writeSize += resLen
		//logger.Log("Read stream size: " + strconv.Itoa(readSize << 3) + "Byte")
	}
}

var bufferCount int32
var ResBufferPool = sync.Pool{
	New: func() interface{} {
		atomic.AddInt32(&bufferCount, 1)
		logger.Log("ResBufferPool: Count of new buffer: " + strconv.Itoa(int(bufferCount)))
		return new(bytes.Buffer)
	},
}
