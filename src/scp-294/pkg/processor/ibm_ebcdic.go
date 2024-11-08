package processor

import mycharset "myutil/charset"

type IbmEBCDIC struct {
	name     string
	procType ProcType
	charset  mycharset.ICharset
}

func NewIbmEBCDIC(name string, procType ProcType) *IbmEBCDIC {
	instance := &IbmEBCDIC{
		name:     name,
		procType: procType,
		charset:  mycharset.IbmEBCDIC{},
	}
	return instance
}

func (p *IbmEBCDIC) Process(arr []byte) ([]byte, error) {
	if p.procType == Encode {
		return p.charset.Encode(arr)
	} else if p.procType == Decode {
		return p.charset.Decode(arr)
	}
	return arr, nil
}

func (p *IbmEBCDIC) GetName() string {
	return p.name
}

func (p *IbmEBCDIC) GetType() ProcType {
	return p.procType
}
