package compress

import (
	"io"
)

type ZstdCompressor struct{}

func (compressor ZstdCompressor) NewWriter(writer io.Writer, level int) ReaderFromWriteCloser {
	return NewZstdReaderFromWriter(writer, level)
}

func (compressor ZstdCompressor) FileExtension() string {
	return ZstdFileExtension
}
