package utils

import (
	"mime/multipart"
	"strings"
)

type TempBuffer struct {
	data []byte
	off  int
	len  int
	cell strings.Builder
}

func CreateTempBuffer() TempBuffer {
	return TempBuffer{
		data: make([]byte, 1024),
		off:  0,
		len:  0,
	}
}

type Page[T any] struct {
	pageNum      int
	buffer       *[]T
	pageSize     int
	length       int
	funcStrToNum func(string) T
}

func (page *Page[T]) AppendValue(str string) {
	(*page.buffer)[page.length] = page.funcStrToNum(str)
	page.length++
}

func (page *Page[T]) IsEmpty() bool {
	return page.length <= 0
}

func (page *Page[T]) IsFull() bool {
	return page.length >= page.pageSize
}

func (page *Page[T]) AppendData(tempBuffer *TempBuffer, file multipart.File) (err error) {
	for {
		if tempBuffer.len == 0 {
			tempBuffer.len, err = file.Read(tempBuffer.data)
			if err != nil {
				if tempBuffer.cell.Len() > 0 {
					page.AppendValue(tempBuffer.cell.String())
					tempBuffer.cell.Reset()
				}
				return err
			}
		}

		var val byte
		endOff := tempBuffer.off + tempBuffer.len
		for tempBuffer.off < endOff {
			val = tempBuffer.data[tempBuffer.off]
			tempBuffer.off++
			if (val >= '0' && val <= '9') || (val >= 'a' && val <= 'f') || (val >= 'A' && val <= 'F') || val == '-' {
				tempBuffer.cell.WriteByte(val)
			} else {
				if tempBuffer.cell.Len() > 0 {
					page.AppendValue(tempBuffer.cell.String())
					tempBuffer.cell.Reset()
					if page.IsFull() {
						tempBuffer.len = endOff - tempBuffer.off
						if tempBuffer.len == 0 {
							tempBuffer.off = 0
						}
						return nil
					}
				}
			}
		}

		tempBuffer.off = 0
		tempBuffer.len = 0
	}
}

func CreateEmptyPage[T any](pageNum int, buffer *[]T, funcStrToNum func(string) T) *Page[T] {
	return &Page[T]{
		pageNum:      pageNum,
		buffer:       buffer,
		pageSize:     cap(*buffer),
		length:       0,
		funcStrToNum: funcStrToNum,
	}
}
