package crc

const (
	polynom = 0xEDB88320
)

// ComputeCRC32 вычисляет контрольную сумму CRC-32 для заданного массива байтов
func ComputeCRC32(data []byte) uint32 {
	var crc uint32 = 0xFFFFFFFF

	for _, b := range data {
		crc ^= uint32(b) // XOR с текущим байтом
		for i := 0; i < 8; i++ {
			if crc&1 != 0 {
				crc = (crc >> 1) ^ polynom
			} else {
				crc >>= 1
			}
		}
	}

	return ^crc
}
