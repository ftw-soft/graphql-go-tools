package websocket

import (
	"bytes"
	"io"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsflate"
)

func decompressMessage(f ws.Frame) (ws.Frame, error) {
	var (
		compressed bool
		err        error
	)

	f.Header, compressed, err = wsflate.UnsetBit(f.Header)
	if err != nil {
		return f, err
	}

	if !compressed {
		return f, nil
	}

	decompressorReader := newDecompressor()
	defer freeDecompressor(decompressorReader)

	decompressorReader.Reset(bytes.NewReader(f.Payload))

	buf := newBuffer()
	defer freeBuffer(buf)

	if _, err = io.Copy(buf, decompressorReader); err != nil {
		return f, err
	}

	if err := decompressorReader.Close(); err != nil {
		return f, err
	}

	f.Payload = buf.Bytes()
	f.Header.Length = int64(len(f.Payload))

	return f, nil
}

func compressMessage(message ws.Frame) (ws.Frame, error) {
	var err error

	buf := newBuffer()
	defer freeBuffer(buf)

	compressorWriter := newCompressor()
	defer freeCompressor(compressorWriter)

	compressorWriter.Reset(buf)

	if _, err = compressorWriter.Write(message.Payload); err != nil {
		return message, err
	}

	if err := compressorWriter.Flush(); err != nil {
		return message, err
	}

	message.Payload = buf.Bytes()
	message.Header.Length = int64(len(message.Payload))
	message.Header, err = wsflate.SetBit(message.Header)

	if err != nil {
		return message, err
	}

	return message, nil
}
