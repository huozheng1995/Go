package file

import (
	"mime/multipart"
	"myutil"
)

type HexStrMultipartFile struct {
	*StrToInt64Reader
}

func NewHexStrMultipartFile(file multipart.File) *HexStrMultipartFile {
	return &HexStrMultipartFile{
		StrToInt64Reader: &StrToInt64Reader{
			buf:            make([]byte, 64*1024),
			bufPos:         0,
			bufSize:        0,
			reader:         file,
			funcStrToInt64: myutil.HexStrToInt64,
		},
	}
}

type DecStrMultipartFile struct {
	*StrToInt64Reader
}

func NewDecStrMultipartFile(file multipart.File) *DecStrMultipartFile {
	return &DecStrMultipartFile{
		StrToInt64Reader: &StrToInt64Reader{
			buf:            make([]byte, 64*1024),
			bufPos:         0,
			bufSize:        0,
			reader:         file,
			funcStrToInt64: myutil.Int64StrToInt64,
		},
	}
}

type BinStrMultipartFile struct {
	*StrToInt64Reader
}

func NewBinStrMultipartFile(file multipart.File) *BinStrMultipartFile {
	return &BinStrMultipartFile{
		StrToInt64Reader: &StrToInt64Reader{
			buf:            make([]byte, 64*1024),
			bufPos:         0,
			bufSize:        0,
			reader:         file,
			funcStrToInt64: myutil.BinStrToInt64,
		},
	}
}
