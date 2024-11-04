package controller

import (
	"github.com/edward/scp-294/logger"
	"strconv"
	"sync"
	"sync/atomic"
)

var byteBufferCount int32
var reqByteBufferPool = &sync.Pool{
	New: func() interface{} {
		atomic.AddInt32(&byteBufferCount, 1)
		logger.Logger.Log("Main", "reqByteBufferPool: Count of new buffer: "+strconv.Itoa(int(byteBufferCount)))
		return make([]byte, 4096)
	},
}

var int64BufferCount int32
var reqInt64BufferPool = &sync.Pool{
	New: func() interface{} {
		atomic.AddInt32(&int64BufferCount, 1)
		logger.Logger.Log("Main", "reqInt64BufferPool: Count of new buffer: "+strconv.Itoa(int(int64BufferCount)))
		return make([]int64, 4096>>3)
	},
}
