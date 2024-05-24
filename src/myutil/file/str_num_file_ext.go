package file

import (
	"io"
	"mime/multipart"
	"myutil"
	"os"
)

// String To Byte

func NewStrHex2File(file multipart.File) *StrNumFile[byte] {
	return &StrNumFile[byte]{
		buf:          make([]byte, 64*1024),
		bufPos:       0,
		bufSize:      0,
		file:         file,
		funcStrToNum: myutil.Hex2StrToByte{},
	}
}

func NewStrInt8File(file multipart.File) *StrNumFile[byte] {
	return &StrNumFile[byte]{
		buf:          make([]byte, 64*1024),
		bufPos:       0,
		bufSize:      0,
		file:         file,
		funcStrToNum: myutil.Int8StrToByte{},
	}
}

func NewStrByteFile(file multipart.File) *StrNumFile[byte] {
	return &StrNumFile[byte]{
		buf:          make([]byte, 64*1024),
		bufPos:       0,
		bufSize:      0,
		file:         file,
		funcStrToNum: myutil.ByteStrToByte{},
	}
}

// String To Int64

func NewStrHexFile(file multipart.File) *StrNumFile[int64] {
	return &StrNumFile[int64]{
		buf:          make([]byte, 64*1024),
		bufPos:       0,
		bufSize:      0,
		file:         file,
		funcStrToNum: myutil.HexStrToInt64{},
	}
}

func NewStrDecFile(file multipart.File) *StrNumFile[int64] {
	return &StrNumFile[int64]{
		buf:          make([]byte, 64*1024),
		bufPos:       0,
		bufSize:      0,
		file:         file,
		funcStrToNum: myutil.DecStrToInt64{},
	}
}

func NewStrBinFile(file multipart.File) *StrNumFile[int64] {
	return &StrNumFile[int64]{
		buf:          make([]byte, 64*1024),
		bufPos:       0,
		bufSize:      0,
		file:         file,
		funcStrToNum: myutil.BinStrToInt64{},
	}
}

// String To Byte, but with os.File

type StrHex2OSFile struct {
	*StrNumFile[byte]
}

func NewStrHex2OSFile(fileUri string) (*StrHex2OSFile, error) {
	file, err := os.Open(fileUri)
	if err != nil {
		return nil, err
	}

	return &StrHex2OSFile{
		StrNumFile: &StrNumFile[byte]{
			buf:          make([]byte, 64*1024),
			bufPos:       0,
			bufSize:      0,
			file:         file,
			funcStrToNum: myutil.Hex2StrToByte{},
		},
	}, nil
}

func (h *StrHex2OSFile) ReadAll() ([]byte, error) {
	var size int
	if info, err := h.StrNumFile.file.(*os.File).Stat(); err == nil {
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
