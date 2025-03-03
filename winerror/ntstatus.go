package winerror

type NTStatus uint32

func (status NTStatus) Sev() uint8 {
	return uint8((status & 0xC0000000) >> 30)
}

func (status NTStatus) C() bool {
	return status&0x20000000 != 0
}

func (status NTStatus) N() bool {
	return status&0x10000000 != 0
}

func (status NTStatus) Facility() uint16 {
	return uint16((status & 0x0FFF0000) >> 16)
}

func (status NTStatus) Code() uint16 {
	return uint16(status & 0x0000FFFF)
}
