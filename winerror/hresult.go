package winerror

type HResult uint32

func (hr HResult) S() bool {
	return hr&0x80000000 != 0
}

func (hr HResult) R() bool {
	return hr&0x40000000 != 0
}

func (hr HResult) C() bool {
	return hr&0x20000000 != 0
}

func (hr HResult) N() bool {
	return hr&0x10000000 != 0
}

func (hr HResult) X() bool {
	return hr&0x08000000 != 0
}

func (hr HResult) Facility() uint16 {
	return uint16((hr & 0x07FF0000) >> 16)
}

func (hr HResult) Code() uint16 {
	return uint16(hr & 0xFFFF)
}
