package util

import (
	"github.com/edward/scp-294/logger"
	"io"
	"myutil"
	myfile "myutil/file"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// text to num array

func TextToNums[T any](text string, strToNum myutil.NumUtil[T]) []T {
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
			continue
		}
		if builder.Len() > 0 {
			result = append(result, strToNum.ToNum(builder.String()))
			builder.Reset()
		}
	}

	return result
}

// file to num array

func FileToNums[T any](file myfile.INumFile[T], bufferPool *sync.Pool, readChan chan []T) {
	defer close(readChan)
	for {
		buf := bufferPool.Get().([]T)
		if len(buf) < cap(buf) {
			logger.Logger.Log("Main", "Resize the buffer from "+strconv.Itoa(len(buf))+" to "+strconv.Itoa(cap(buf)))
			buf = buf[0:cap(buf)]
		}
		len, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				logger.Logger.Log("Main", "File stream read done")
			} else {
				logger.Logger.Log("Main", "Failed to read file stream, error: "+err.Error())
			}
			if len > 0 {
				readChan <- buf[0:len]
			} else {
				bufferPool.Put(buf)
			}
			return
		}

		readChan <- buf[0:len]
	}
}

func ReadFromChannelAndRespond[T any](readChan <-chan []T, numsToResp NumsToResp[T], bufferPool *sync.Pool, w http.ResponseWriter) {
	readSize := 0
	writeSize := 0
	for {
		buf, ok := <-readChan
		bufLen := len(buf)
		if !ok || len(buf) <= 0 {
			readKB := strconv.Itoa((readSize * numsToResp.GetByteNum()) >> 10)
			writeKB := strconv.Itoa((writeSize * numsToResp.GetByteNum()) >> 10)
			logger.Logger.Log("Main", "Read channel done, total size: "+strconv.Itoa(readSize)+"("+readKB+"KB)")
			logger.Logger.Log("Main", "Write stream done, total size: "+strconv.Itoa(writeSize)+"("+writeKB+"KB)")
			return
		}

		resp := numsToResp.ToResp(buf)
		respLen := len(resp)

		bufferPool.Put(buf)
		w.Write(resp)

		readSize += bufLen
		writeSize += respLen
		//logger.Logger.Log("Main", "Read stream size: "+strconv.Itoa((readSize*numsToResp.GetByteNum()))+"Byte")
	}
}
