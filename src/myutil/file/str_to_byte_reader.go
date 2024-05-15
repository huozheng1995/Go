package file

import (
	"io"
	"myutil"
	"strconv"
	"strings"
)

type StrToByteReader struct {
	buf           []byte
	bufPos        int
	bufSize       int
	reader        io.ReadCloser
	funcStrToByte myutil.StrToByte
}

func (h *StrToByteReader) innerRead() (err error) {
	h.bufSize, err = h.reader.Read(h.buf)
	h.bufPos = 0
	if myutil.Logger != nil {
		myutil.Logger.Log("StrToByteReader", "Read bytes "+strconv.Itoa(h.bufSize))
	}
	return err
}

func (h *StrToByteReader) Read(p []byte) (int, error) {
	pPos := 0
	var builder strings.Builder
	for pPos < cap(p) {
		//read data
		if h.bufPos >= h.bufSize {
			err := h.innerRead()
			if err != nil {
				if err == io.EOF && builder.Len() > 0 {
					p[pPos] = h.funcStrToByte(builder.String())
					pPos++
					builder.Reset()
				}
				if myutil.Logger != nil {
					myutil.Logger.Log("StrToByteReader", "Return bytes "+strconv.Itoa(pPos))
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
				p[pPos] = h.funcStrToByte(builder.String())
				pPos++
				builder.Reset()
			}
		}
	}
	if myutil.Logger != nil {
		myutil.Logger.Log("StrToByteReader", "Return bytes "+strconv.Itoa(pPos))
	}

	return pPos, nil
}

func (h *StrToByteReader) Close() error {
	return h.reader.Close()
}
