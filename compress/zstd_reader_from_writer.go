package compress

import (
	"io"

	"github.com/DataDog/zstd"
)

type ZstdReaderFromWriter struct {
	zstd.Writer
}

func NewZstdReaderFromWriter(dst io.Writer, level int) *ZstdReaderFromWriter {
	zstdWriter := zstd.NewWriterLevel(dst, level)
	return &ZstdReaderFromWriter{Writer: *zstdWriter}
}

func (writer *ZstdReaderFromWriter) ReadFrom(reader io.Reader) (n int64, err error) {
	n, err = fastCopy(writer, reader)
	return
}
