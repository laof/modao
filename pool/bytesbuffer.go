package pool

import (
	"bytes"
	"sync"
)

var bytesBufPool = sync.Pool{
	New: func() interface{} { return &bytes.Buffer{} },
}

// GetBytesBuffer returns a bytes.buffer from pool.
func GetBytesBuffer() *bytes.Buffer {
	return bytesBufPool.Get().(*bytes.Buffer)
}

// PutBytesBuffer puts a bytes.buffer into pool.
func PutBytesBuffer(buf *bytes.Buffer) {
	if buf.Cap() > 64<<10 {
		return
	}

	buf.Reset()
	bytesBufPool.Put(buf)
}
