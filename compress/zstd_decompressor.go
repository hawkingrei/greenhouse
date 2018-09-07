package compress

import (
	"io"

	"github.com/DataDog/zstd"
)

type ZstdDecompressor struct{}

func (decompressor ZstdDecompressor) Decompress(dst io.Writer, src io.Reader) error {
	zstdReader := zstd.NewReaderDict(NewUntilEofReader(src), ZstdDict)
	defer zstdReader.Close()
	_, err := fastCopy(dst, zstdReader)
	return err
}

func (decompressor ZstdDecompressor) FileExtension() string {
	return ZstdFileExtension
}
