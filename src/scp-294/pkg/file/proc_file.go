package file

import (
	"github.com/edward/scp-294/pkg/processor"
	"io"
	"myutil"
	"strconv"
	"strings"
)

type ProcFile[T byte] struct {
	buf       []byte
	bufPos    int
	bufSize   int
	file      io.ReadCloser
	processor processor.IProcessor
}

func (obj *ProcFile[T]) innerRead() (err error) {
	obj.bufSize, err = obj.file.Read(obj.buf)
	obj.bufPos = 0
	if myutil.Logger != nil {
		myutil.Logger.Log("ProcFile", "Read bytes "+strconv.Itoa(obj.bufSize))
	}
	return err
}

func (obj *ProcFile[T]) Read(p []T) (int, error) {
	pPos := 0
	var builder strings.Builder
	for pPos < cap(p) {
		//read data
		if obj.bufPos >= obj.bufSize {
			err := obj.innerRead()
			if err != nil {
				if builder.Len() > 0 {
					p[pPos] = obj.numUtil.ToNum(builder.String())
					pPos++
					builder.Reset()
				}
				if myutil.Logger != nil {
					myutil.Logger.Log("ProcFile", "Return objects "+strconv.Itoa(pPos))
				}

				return pPos, err
			}
		}

		//parse data
		val := obj.buf[obj.bufPos]
		obj.bufPos++
		if (val >= '0' && val <= '9') || (val >= 'a' && val <= 'f') || (val >= 'A' && val <= 'F') || val == '-' {
			builder.WriteByte(val)
			continue
		}
		if builder.Len() > 0 {
			p[pPos] = obj.numUtil.ToNum(builder.String())
			pPos++
			builder.Reset()
		}
	}
	if myutil.Logger != nil {
		myutil.Logger.Log("ProcFile", "Return objects "+strconv.Itoa(pPos))
	}

	return pPos, nil
}

func (obj *ProcFile[T]) Close() error {
	return obj.file.Close()
}
