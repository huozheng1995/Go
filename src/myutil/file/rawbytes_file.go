package file

import (
	"mime/multipart"
)

type RawBytesNumFile struct {
	File multipart.File
}

func (h *RawBytesNumFile) Read(p []byte) (int, error) {
	return h.File.Read(p)
}
