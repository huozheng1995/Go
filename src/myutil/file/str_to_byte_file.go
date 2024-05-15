package file

import (
	"io"
	"mime/multipart"
	"myutil"
	"os"
)

type Hex2StrMultipartFile struct {
	*StrToByteReader
}

func NewHex2StrMultipartFile(file multipart.File) *Hex2StrMultipartFile {
	return &Hex2StrMultipartFile{
		StrToByteReader: &StrToByteReader{
			buf:           make([]byte, 64*1024),
			bufPos:        0,
			bufSize:       0,
			reader:        file,
			funcStrToByte: myutil.Hex2StrToByte,
		},
	}
}

type Int8StrMultipartFile struct {
	*StrToByteReader
}

func NewInt8StrMultipartFile(file multipart.File) *Int8StrMultipartFile {
	return &Int8StrMultipartFile{
		StrToByteReader: &StrToByteReader{
			buf:           make([]byte, 64*1024),
			bufPos:        0,
			bufSize:       0,
			reader:        file,
			funcStrToByte: myutil.Int8StrToByte,
		},
	}
}

type ByteStrMultipartFile struct {
	*StrToByteReader
}

func NewByteStrMultipartFile(file multipart.File) *ByteStrMultipartFile {
	return &ByteStrMultipartFile{
		StrToByteReader: &StrToByteReader{
			buf:           make([]byte, 64*1024),
			bufPos:        0,
			bufSize:       0,
			reader:        file,
			funcStrToByte: myutil.ByteStrToByte,
		},
	}
}

type Hex2StrOSFile struct {
	*StrToByteReader
}

func NewHex2ByteOSFile(fileUri string) (*Hex2StrOSFile, error) {
	file, err := os.Open(fileUri)
	if err != nil {
		return nil, err
	}

	return &Hex2StrOSFile{
		StrToByteReader: &StrToByteReader{
			buf:           make([]byte, 64*1024),
			bufPos:        0,
			bufSize:       0,
			reader:        file,
			funcStrToByte: myutil.Hex2StrToByte,
		},
	}, nil
}

func (h *Hex2StrOSFile) ReadAll() ([]byte, error) {
	var size int
	if info, err := h.StrToByteReader.reader.(*os.File).Stat(); err == nil {
		size64 := info.Size()
		if int64(int(size64)) == size64 {
			size = int(size64)
		}
	}
	size++ // one byte for final read at EOF
	if size < 512 {
		size = 512
	}

	result := make([]byte, size>>1)
	n, err := h.Read(result)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return result[0:n], nil
}
