package stream

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type HexFileStream struct {
	file    *os.File
	buf     []byte
	bufPos  int
	bufSize int
}

func NewHexFileStream(fileUri string) (*HexFileStream, error) {
	file, err := os.Open(fileUri)
	if err != nil {
		return nil, err
	}

	hexFileStream := HexFileStream{
		file:    file,
		buf:     make([]byte, 64*1024),
		bufPos:  0,
		bufSize: 0,
	}
	err = hexFileStream.innerRead()
	if err != nil {
		return nil, err
	}

	return &hexFileStream, err
}

func (h *HexFileStream) innerRead() (err error) {
	h.bufSize, err = h.file.Read(h.buf)
	h.bufPos = 0
	if err != nil {
		return err
	}
	return nil
}

func (h *HexFileStream) Read(p []byte) (int, error) {
	pPos := 0
	var builder strings.Builder
	for pPos < len(p) {
		val := h.buf[h.bufPos]
		h.bufPos++
		if (val >= '0' && val <= '9') || (val >= 'a' && val <= 'f') || (val >= 'A' && val <= 'F') || val == '-' {
			builder.WriteByte(val)
		} else {
			if builder.Len() > 0 {
				intVal, _ := strconv.ParseInt(builder.String(), 16, 64)
				builder.Reset()
				p[pPos] = byte(intVal)
				pPos++
			}
		}

		if h.bufPos >= h.bufSize {
			err := h.innerRead()
			fmt.Println("read next..." + strconv.Itoa(h.bufSize))
			if err != nil {
				return pPos, err
			}
		}
	}

	return pPos, nil
}

func (h *HexFileStream) Close() error {
	return h.file.Close()
}
