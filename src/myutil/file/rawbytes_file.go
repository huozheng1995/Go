package file

import (
	"mime/multipart"
	"unsafe"
)

type RawBytesNumFile[T any] struct {
	File multipart.File
}

func (h *RawBytesNumFile[T]) Read(p []T) (int, error) {
	bytes := *(*[]byte)(unsafe.Pointer(&p))
	return h.File.Read(bytes)
}
