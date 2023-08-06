package file

import (
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
	err = hexFileStream.innerRead()
	if err != nil {
		return nil, err
	}

	return &hexFileStream, err
}

func (h *HexFile) innerRead() (err error) {
	h.bufSize, err = h.file.Read(h.buf)
	h.bufPos = 0
	if err != nil {
		return err
	}
	return nil
}

func (h *HexFile) Read(p []byte) (int, error) {
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
			if myutil.Logger != nil {
				myutil.Logger.Log("HexFile", "Read bytes "+strconv.Itoa(h.bufSize))
			}
			if err != nil {
				if myutil.Logger != nil {
					myutil.Logger.Log("HexFile", "Return bytes "+strconv.Itoa(pPos))
				}
				return pPos, err
			}
		}
	}

	if myutil.Logger != nil {
		myutil.Logger.Log("HexFile", "Return bytes "+strconv.Itoa(pPos))
	}
	return pPos, nil
}

func (h *HexFile) Close() error {
	return h.file.Close()
}
