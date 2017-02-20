package hton

func U64(buffer []byte, index int, value uint64) []byte {
	buffer[index] = byte(0xff & (value >> 56))
	buffer[index+1] = byte(0xff & (value >> 48))
	buffer[index+2] = byte(0xff & (value >> 40))
	buffer[index+3] = byte(0xff & (value >> 32))
	buffer[index+4] = byte(0xff & (value >> 24))
	buffer[index+5] = byte(0xff & (value >> 16))
	buffer[index+6] = byte(0xff & (value >> 8))
	buffer[index+7] = byte(0xff & (value))
	return buffer
}

func U32(buffer []byte, index int, value uint32) []byte {
	buffer[index] = byte(0xff & (value >> 24))
	buffer[index+1] = byte(0xff & (value >> 16))
	buffer[index+2] = byte(0xff & (value >> 8))
	buffer[index+3] = byte(0xff & (value))
	return buffer
}

func U16(buffer []byte, index int, value uint16) []byte {
	buffer[index] = byte(0xff & (value >> 8))
	buffer[index+1] = byte(0xff & (value))
	return buffer
}

func U8(buffer []byte, index int, value uint8) []byte {
	buffer[index] = value
	return buffer
}
