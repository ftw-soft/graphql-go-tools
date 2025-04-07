package websocket

import (
	"bytes"
	"sync"
)

// The byte buffer may be returned to the pool via Put after the use
// in order to minimize GC overhead.
// nolint:gochecknoglobals // should be global
var bufferPool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

func newBuffer() *bytes.Buffer {
	if ret := bufferPool.Get(); ret != nil {
		return ret.(*bytes.Buffer)
	} else {
		return new(bytes.Buffer)
	}
}

func freeBuffer(b *bytes.Buffer) {
	b.Reset()
	bufferPool.Put(b)
}
