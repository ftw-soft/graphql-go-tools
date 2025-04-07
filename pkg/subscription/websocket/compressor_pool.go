package websocket

import (
	"compress/flate"
	"io"
	"sync"

	"github.com/gobwas/ws/wsflate"
)

// nolint:gochecknoglobals // should be global
var wsCompressorPool = sync.Pool{
	New: func() any {
		return wsflate.NewWriter(nil, func(w io.Writer) wsflate.Compressor {
			f, _ := flate.NewWriter(w, flate.DefaultCompression) //nolint:errcheck // ok
			return f
		})
	},
}

//nolint:gochecknoglobals // should be global
var wsDecompressorPool = sync.Pool{
	New: func() any {
		return wsflate.NewReader(nil, func(r io.Reader) wsflate.Decompressor {
			return newResetDecorator(flate.NewReader(r))
		})
	},
}

func newCompressor() *wsflate.Writer {
	return wsCompressorPool.Get().(*wsflate.Writer)
}

func freeCompressor(w *wsflate.Writer) {
	w.Close()
	wsCompressorPool.Put(w)
}

func newDecompressor() *wsflate.Reader {
	return wsDecompressorPool.Get().(*wsflate.Reader)
}

func freeDecompressor(r *wsflate.Reader) {
	r.Close()
	wsDecompressorPool.Put(r)
}

func newResetDecorator(rc io.ReadCloser) io.ReadCloser {
	return &resetDecorator{ReadCloser: rc}
}

type resetDecorator struct {
	io.ReadCloser
}

func (f *resetDecorator) Reset(r io.Reader) {
	if x, ok := f.ReadCloser.(DictResetter); ok {
		x.Reset(r, nil) //nolint:errcheck // ok
	}
}

type DictResetter interface {
	Reset(r io.Reader, dict []byte) error
}
