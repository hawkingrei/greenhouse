package compress

import "io"

const (
	Lz4AlgorithmName  = "lz4"
	LzmaAlgorithmName = "lzma"
	ZstdAlgorithmName = "zstd"

	Lz4FileExtension  = "lz4"
	LzmaFileExtension = "lzma"
	ZstdFileExtension = "zst"
	LzoFileExtension  = "lzo"

	CompressedBlockMaxSize = 20 << 20
)

type Compressor interface {
	NewWriter(writer io.Writer, level int) ReaderFromWriteCloser
	FileExtension() string
}

type Decompressor interface {
	Decompress(dst io.Writer, src io.Reader) error
	FileExtension() string
}

var Compressors = map[string]Compressor{
	ZstdAlgorithmName: ZstdCompressor{},
}

var Decompressors = map[string]Decompressor{
	ZstdAlgorithmName: ZstdDecompressor{},
}
