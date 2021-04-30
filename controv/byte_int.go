package controv

import (
	"bytes"
	"encoding/binary"
)

func IntToBytes(n int, order binary.ByteOrder) []byte {
	l := 8
	if (32 << (^uint(0) >> 63)) == 32 {
		l = 4
	}

	bytebuf := bytes.NewBuffer(make([]byte, l))
	binary.Write(bytebuf, order, n)
	return bytebuf.Bytes()
}
