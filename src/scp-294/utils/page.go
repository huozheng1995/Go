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
	buffer       []T
	pageSize     int
	index        int
	funcStrToNum func(string) T
}

func (page *Page[T]) AppendValue(str string) {
	page.buffer[page.index] = page.funcStrToNum(str)
	page.index++
}

func (page *Page[T]) IsEmpty() bool {
	return page.index <= 0
}

func (page *Page[T]) IsFull() bool {
	return page.index >= page.pageSize
}

func (page *Page[T]) GetBuffer() []T {
	return page.buffer[:page.index]
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
			if (val >= '0' && val <= '9') || (val >= 'a' && val <= 'f') || (val >= 'A' && val <= 'F') || val == '-' {
				tempBuffer.cell.WriteByte(val)
			} else {
				if tempBuffer.cell.Len() > 0 {
					page.AppendValue(tempBuffer.cell.String())
					tempBuffer.cell.Reset()
					if page.IsFull() {
						tempBuffer.off++
						return nil
					}
				}
			}

			tempBuffer.off++
		}

		tempBuffer.off = 0
		tempBuffer.len = 0
	}
}

func CreateEmptyPage[T any](pageNum int, buffer []T, funcStrToNum func(string) T) Page[T] {
	return Page[T]{
		pageNum:      pageNum,
		buffer:       buffer,
		pageSize:     len(buffer),
		index:        0,
		funcStrToNum: funcStrToNum,
	}
}
