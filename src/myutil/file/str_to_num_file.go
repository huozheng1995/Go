package file

import (
	"io"
	"mime/multipart"
	"myutil"
	"strconv"
	"strings"
)

// String To Byte

func NewHex2StrMultipartFile(file multipart.File) *StrToNumFile[byte] {
	return &StrToNumFile[byte]{
		buf:          make([]byte, 64*1024),
		bufPos:       0,
		bufSize:      0,
		file:         file,
		funcStrToNum: myutil.Hex2StrToByte,
	}
}

func NewInt8StrMultipartFile(file multipart.File) *StrToNumFile[byte] {
	return &StrToNumFile[byte]{
		buf:          make([]byte, 64*1024),
		bufPos:       0,
		bufSize:      0,
		file:         file,
		funcStrToNum: myutil.Int8StrToByte,
	}
}

func NewByteStrMultipartFile(file multipart.File) *StrToNumFile[byte] {
	return &StrToNumFile[byte]{
		buf:          make([]byte, 64*1024),
		bufPos:       0,
		bufSize:      0,
		file:         file,
		funcStrToNum: myutil.ByteStrToByte,
	}
}

// String To Int64

func NewHexStrMultipartFile(file multipart.File) *StrToNumFile[int64] {
	return &StrToNumFile[int64]{
		buf:          make([]byte, 64*1024),
		bufPos:       0,
		bufSize:      0,
		file:         file,
		funcStrToNum: myutil.HexStrToInt64,
	}
}

func NewDecStrMultipartFile(file multipart.File) *StrToNumFile[int64] {
	return &StrToNumFile[int64]{
		buf:          make([]byte, 64*1024),
		bufPos:       0,
		bufSize:      0,
		file:         file,
		funcStrToNum: myutil.DecStrToInt64,
	}
}

func NewBinStrMultipartFile(file multipart.File) *StrToNumFile[int64] {
	return &StrToNumFile[int64]{
		buf:          make([]byte, 64*1024),
		bufPos:       0,
		bufSize:      0,
		file:         file,
		funcStrToNum: myutil.BinStrToInt64,
	}
}

type StrToNumFile[T any] struct {
	buf          []byte
	bufPos       int
	bufSize      int
	file         io.ReadCloser
	funcStrToNum func(str string) T
}

func (h *StrToNumFile[T]) innerRead() (err error) {
	h.bufSize, err = h.file.Read(h.buf)
	h.bufPos = 0
	if myutil.Logger != nil {
		myutil.Logger.Log("StrToNumFile", "Read bytes "+strconv.Itoa(h.bufSize))
	}
	return err
}

func (h *StrToNumFile[T]) Read(p []T) (int, error) {
	pPos := 0
	var builder strings.Builder
	for pPos < cap(p) {
		//read data
		if h.bufPos >= h.bufSize {
			err := h.innerRead()
			if err != nil {
				if err == io.EOF && builder.Len() > 0 {
					p[pPos] = h.funcStrToNum(builder.String())
					pPos++
					builder.Reset()
				}
				if myutil.Logger != nil {
					myutil.Logger.Log("StrToNumFile", "Return objects "+strconv.Itoa(pPos))
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
				p[pPos] = h.funcStrToNum(builder.String())
				pPos++
				builder.Reset()
			}
		}
	}
	if myutil.Logger != nil {
		myutil.Logger.Log("StrToNumFile", "Return objects "+strconv.Itoa(pPos))
	}

	return pPos, nil
}

func (h *StrToNumFile[T]) Close() error {
	return h.file.Close()
}
