package utils

import (
	"mime/multipart"
	"strings"
)

type Page struct {
	PageNum  int
	Buffer   []byte
	PageSize int
	Index    int
}

func CreateEmptyPage(pageNum int, buffer []byte) Page {
	return Page{
		PageNum:  pageNum,
		Buffer:   buffer,
		PageSize: len(buffer),
		Index:    0,
	}
}

func FillPage(page *Page, preBuffer []byte, preOff, preLen int, tempCell strings.Builder, file multipart.File,
	funcStrToDecByte StrToByte) (err error, newPreBuffer []byte, newPreOff, newPreLen int, newTempCell strings.Builder) {
	for {
		if preLen == 0 {
			preOff = 0
			preLen, err = file.Read(preBuffer)
			if err != nil {
				if tempCell.Len() > 0 {
					page.Buffer[page.Index] = funcStrToDecByte(tempCell.String())
					page.Index++
					tempCell.Reset()
				}
				return err, preBuffer, preOff, preLen, tempCell
			}
		}

		var val byte
		for i := preOff; i < preOff+preLen; i++ {
			val = preBuffer[i]
			if (val >= '0' && val <= '9') || (val >= 'a' && val <= 'f') || (val >= 'A' && val <= 'F') || val == '-' {
				tempCell.WriteByte(val)
			} else {
				if tempCell.Len() > 0 {
					page.Buffer[page.Index] = funcStrToDecByte(tempCell.String())
					page.Index++
					tempCell.Reset()
					if page.Index == page.PageSize {
						i++
						return nil, preBuffer, i, preOff + preLen - i, tempCell
					}
				}
			}
		}
		preLen = 0
	}
}
