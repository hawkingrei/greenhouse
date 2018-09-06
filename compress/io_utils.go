package compress

import "io"

type ReaderFromWriteCloser interface {
	io.ReaderFrom
	io.WriteCloser
}

type UntilEOFReader struct {
	underlying io.Reader
	isEOF      bool
}

func NewUntilEofReader(underlying io.Reader) *UntilEOFReader {
	return &UntilEOFReader{underlying, false}
}

func (reader *UntilEOFReader) Read(p []byte) (n int, err error) {
	if reader.isEOF {
		return 0, io.EOF
	}
	n, err = reader.underlying.Read(p)
	if err == io.EOF {
		reader.isEOF = true
	}
	return
}

func fastCopy(dst io.Writer, src io.Reader) (int64, error) {
	n := int64(0)
	buf := make([]byte, CompressedBlockMaxSize)
	for {
		m, readingErr := src.Read(buf)
		if readingErr != nil && readingErr != io.EOF {
			return n, readingErr
		}
		m, writingErr := dst.Write(buf[:m])
		n += int64(m)
		if writingErr != nil || readingErr == io.EOF {
			return n, writingErr
		}
	}
}
