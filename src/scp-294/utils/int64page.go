package utils

import (
	"mime/multipart"
	"strings"
)

type Int64Page struct {
	PageNum  int
	Buffer   []int64
	PageSize int
	Index    int
}

func CreateEmptyInt64Page(pageNum int, buffer []int64) Int64Page {
	return Int64Page{
		PageNum:  pageNum,
		Buffer:   buffer,
		PageSize: len(buffer),
		Index:    0,
	}
}

func FillInt64Page(page *Int64Page, preBuffer []byte, preOff, preLen int, tempCell strings.Builder, file multipart.File,
	funcStrToInt64 StrToInt64) (err error, newPreBuffer []byte, newPreOff, newPreLen int, newTempCell strings.Builder) {
	for {
		if preLen == 0 {
			preOff = 0
			preLen, err = file.Read(preBuffer)
			if err != nil {
				if tempCell.Len() > 0 {
					page.Buffer[page.Index] = funcStrToInt64(tempCell.String())
					page.Index++
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
					page.Buffer[page.Index] = funcStrToInt64(tempCell.String())
					page.Index++
					tempCell.Reset()
					if page.Index == page.PageSize {
						i++
						return nil, preBuffer, preOff + i, preLen - i, tempCell
					}
				}
			}
		}
		preLen = 0
	}
}
