package charset

import (
	"bytes"
	"myutil/codepage"
)

type IbmEBCDIC struct {
	Charset
}

func (c IbmEBCDIC) GetName() string {
	return "IbmEBCDIC"
}

func (c IbmEBCDIC) Encode(arr []byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	for _, val := range arr {
		buf.WriteByte(codepage.IbmEBCDIC[val])
	}

	return buf.Bytes(), nil
}

func (c IbmEBCDIC) Decode(arr []byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	for _, val := range arr {
		for idx, pageVal := range codepage.IbmEBCDIC {
			if val == pageVal {
				buf.WriteByte(byte(idx))
				break
			}
		}
	}

	return buf.Bytes(), nil
}
