package util

import "encoding/binary"

func SerialzeInt64(value int64) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(value))
	return buf
}
func DeserializeInt64(data []byte) int64 {
	return int64(binary.LittleEndian.Uint64(data))
}
