package utils

import (
	"mime/multipart"
	"strings"
)

type Page[T any] struct {
	PageNum      int
	Buffer       []T
	PageSize     int
	Index        int
	funcStrToNum func(string) T
}

func (page *Page[T]) AppendValue(str string) {
	page.Buffer[page.Index] = page.funcStrToNum(str)
	page.Index++
}

func (page *Page[T]) IsEOF() bool {
	return page.Index >= page.PageSize
}

func CreateEmptyPage[T any](pageNum int, buffer []T, funcStrToNum func(string) T) Page[T] {
	return Page[T]{
		PageNum:      pageNum,
		Buffer:       buffer,
		PageSize:     len(buffer),
		Index:        0,
		funcStrToNum: funcStrToNum,
	}
}

func FillPage[T any](page *Page[T], preBuffer []byte, preOff, preLen int, tempCell strings.Builder, file multipart.File,
) (err error, newPreBuffer []byte, newPreOff, newPreLen int, newTempCell strings.Builder) {

	for {
		if preLen == 0 {
			preOff = 0
			preLen, err = file.Read(preBuffer)
			if err != nil {
				if tempCell.Len() > 0 {
					page.AppendValue(tempCell.String())
					tempCell.Reset()
				}
				return err, preBuffer, preOff, preLen, tempCell
			}
		}

		var val byte
		for i := 0; i < preLen; i++ {
			val = preBuffer[preOff+i]
			if (val >= '0' && val <= '9') || (val >= 'a' && val <= 'f') || (val >= 'A' && val <= 'F') || val == '-' {
				tempCell.WriteByte(val)
			} else {
				if tempCell.Len() > 0 {
					page.AppendValue(tempCell.String())
					tempCell.Reset()
					if page.IsEOF() {
						i++
						return nil, preBuffer, preOff + i, preLen - i, tempCell
					}
				}
			}
		}
		preLen = 0
	}
}
