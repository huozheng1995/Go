package file

import (
	"io"
	"myutil"
	"os"
	"strconv"
	"strings"
)

type HexFile struct {
	file    *os.File
	buf     []byte
	bufPos  int
	bufSize int
}

func NewHexFile(fileUri string) (*HexFile, error) {
	file, err := os.Open(fileUri)
	if err != nil {
		return nil, err
	}

	hexFileStream := HexFile{
		file:    file,
		buf:     make([]byte, 64*1024),
		bufPos:  0,
		bufSize: 0,
	}

	return &hexFileStream, nil
}

func (h *HexFile) innerRead() (err error) {
	h.bufSize, err = h.file.Read(h.buf)
	h.bufPos = 0

	if myutil.Logger != nil {
		myutil.Logger.Log("HexFile", "Read bytes "+strconv.Itoa(h.bufSize))
	}

	return err
}

func (h *HexFile) Read(p []byte) (int, error) {
	pPos := 0
	var builder strings.Builder
	for pPos < cap(p) {
		//read data
		if h.bufPos >= h.bufSize {
			err := h.innerRead()
			if err != nil {
				if err == io.EOF && builder.Len() > 0 {
					intVal, _ := strconv.ParseInt(builder.String(), 16, 64)
					builder.Reset()
					p[pPos] = byte(intVal)
					pPos++
				}
				if myutil.Logger != nil {
					myutil.Logger.Log("HexFile", "Return bytes "+strconv.Itoa(pPos))
				}
				return pPos, err
			}
		}

		//parse data
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
	}

	if myutil.Logger != nil {
		myutil.Logger.Log("HexFile", "Return bytes "+strconv.Itoa(pPos))
	}
	return pPos, nil
}

func (h *HexFile) ReadAll() ([]byte, error) {
	var size int
	if info, err := h.file.Stat(); err == nil {
		size64 := info.Size()
		if int64(int(size64)) == size64 {
			size = int(size64)
		}
	}
	size++ // one byte for final read at EOF
	if size < 512 {
		size = 512
	}

	result := make([]byte, size>>1)
	n, err := h.Read(result)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return result[0:n], nil
}

func (h *HexFile) Close() error {
	return h.file.Close()
}
