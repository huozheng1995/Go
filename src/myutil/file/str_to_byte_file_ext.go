package file

import (
	"io"
	"myutil"
	"os"
)

type Hex2StrOSFile struct {
	*StrToNumFile[byte]
}

func NewHex2ByteOSFile(fileUri string) (*Hex2StrOSFile, error) {
	file, err := os.Open(fileUri)
	if err != nil {
		return nil, err
	}

	return &Hex2StrOSFile{
		StrToNumFile: &StrToNumFile[byte]{
			buf:          make([]byte, 64*1024),
			bufPos:       0,
			bufSize:      0,
			file:         file,
			funcStrToNum: myutil.Hex2StrToByte,
		},
	}, nil
}

func (h *Hex2StrOSFile) ReadAll() ([]byte, error) {
	var size int
	if info, err := h.StrToNumFile.file.(*os.File).Stat(); err == nil {
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
