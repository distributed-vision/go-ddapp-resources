package ntoh

func U64(buffer []byte, index int) uint64 {
	return (uint64(buffer[index]) << 56) |
		(uint64(buffer[index+1]) << 48) |
		(uint64(buffer[index+2]) << 40) |
		uint64(buffer[index+3])<<32 |
		(uint64(buffer[index+4]) << 24) |
		(uint64(buffer[index+5]) << 16) |
		(uint64(buffer[index+6]) << 8) |
		uint64(buffer[index+7])
}

func U32(buffer []byte, index int) uint32 {
	return (uint32(buffer[index]) << 24) |
		(uint32(buffer[index+1]) << 16) |
		(uint32(buffer[index+2]) << 8) |
		uint32(buffer[index+3])
}

func U16(buffer []byte, index int) uint16 {
	return uint16(buffer[index])<<8 |
		uint16(buffer[index+1])
}

func U8(buffer []byte, index int) uint8 {
	return buffer[index]
}
