package ossz

import (
	"bytes"
	"io"
)

// ReadCloser *os.File  本身就具备  Seek 和 Close
type ReadCloser interface {
	io.ReadCloser
	io.Seeker
}

type BytesReader struct {
	reader *bytes.Reader
}

func NewBytesReader(r *bytes.Reader) ReadCloser {
	return BytesReader{
		reader: r,
	}
}
func (r BytesReader) Read(p []byte) (n int, err error) {
	return r.reader.Read(p)
}

func (r BytesReader) Seek(offset int64, whence int) (ret int64, err error) {
	return r.reader.Seek(offset, whence)
}

func (r BytesReader) Close() error {
	return nil
}
