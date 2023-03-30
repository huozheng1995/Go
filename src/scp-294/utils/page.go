package utils

import (
	"github.com/edward/scp-294/logger"
	"io"
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

func FillPage(page *Page, preBuffer []byte, preOff, preLen int, preCell strings.Builder, file multipart.File,
	funcStrToDecByte StrToDecByte) (currentBuffer []byte, currentOff, currentLen int, currentCell strings.Builder) {
	for {
		if preLen == 0 {
			preOff = 0
			var err error
			preLen, err = file.Read(preBuffer)
			if err != nil {
				if err != io.EOF {
					logger.Log("Failed to read file stream, error: " + err.Error())
				} else {
					if preCell.Len() > 0 {
						page.Buffer[page.Index] = funcStrToDecByte(preCell.String())
						page.Index++
						preCell.Reset()
					}
					logger.Log("File stream read done")
				}
				return
			}
		}

		var val byte
		endIndex := preOff + preLen
		for i := preOff; i < endIndex; i++ {
			val = preBuffer[i]
			if (val >= '0' && val <= '9') || (val >= 'a' && val <= 'f') || (val >= 'A' && val <= 'F') || val == '-' {
				preCell.WriteByte(val)
			} else {
				if preCell.Len() > 0 {
					page.Buffer[page.Index] = funcStrToDecByte(preCell.String())
					page.Index++
					preCell.Reset()
					if page.Index == page.PageSize {
						return preBuffer, i + 1, endIndex - (i + 1), preCell
					}
				}
			}
		}
		preLen = 0
	}
}
