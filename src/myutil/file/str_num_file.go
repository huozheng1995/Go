package file

import (
	"io"
	"myutil"
	"strconv"
	"strings"
)

type StrNumFile[T any] struct {
	buf          []byte
	bufPos       int
	bufSize      int
	file         io.ReadCloser
	funcStrToNum myutil.StrToNum[T]
}

func (h *StrNumFile[T]) innerRead() (err error) {
	h.bufSize, err = h.file.Read(h.buf)
	h.bufPos = 0
	if myutil.Logger != nil {
		myutil.Logger.Log("StrNumFile", "Read bytes "+strconv.Itoa(h.bufSize))
	}
	return err
}

func (h *StrNumFile[T]) Read(p []T) (int, error) {
	pPos := 0
	var builder strings.Builder
	for pPos < cap(p) {
		//read data
		if h.bufPos >= h.bufSize {
			err := h.innerRead()
			if err != nil {
				if builder.Len() > 0 {
					p[pPos] = h.funcStrToNum.ToNum(builder.String())
					pPos++
					builder.Reset()
				}
				if myutil.Logger != nil {
					myutil.Logger.Log("StrNumFile", "Return objects "+strconv.Itoa(pPos))
				}

				return pPos, err
			}
		}

		//parse data
		val := h.buf[h.bufPos]
		h.bufPos++
		if (val >= '0' && val <= '9') || (val >= 'a' && val <= 'f') || (val >= 'A' && val <= 'F') || val == '-' {
			builder.WriteByte(val)
			continue
		}
		if builder.Len() > 0 {
			p[pPos] = h.funcStrToNum.ToNum(builder.String())
			pPos++
			builder.Reset()
		}
	}
	if myutil.Logger != nil {
		myutil.Logger.Log("StrNumFile", "Return objects "+strconv.Itoa(pPos))
	}

	return pPos, nil
}

func (h *StrNumFile[T]) Close() error {
	return h.file.Close()
}
