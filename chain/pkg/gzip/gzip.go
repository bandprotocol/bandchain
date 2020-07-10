package gzip

import (
	"bytes"
	gz "compress/gzip"
	"errors"
	"io"
	"io/ioutil"
)

// Magic bytes to identify gzip. See https://www.ietf.org/rfc/rfc1952.txt section 2.3.1.
var gzipIdent = []byte("\x1F\x8B\x08")

// IsGzipped returns true iff the input is gzipped.
func IsGzipped(src []byte) bool {
	return bytes.Equal(gzipIdent, src[0:3])
}

// Uncompress performs gzip uncompression and returns the result. Returns error if the
// input file is not in gzipped format or if the result's size exceeds maxSize.
func Uncompress(src []byte, maxSize int64) ([]byte, error) {
	zr, err := gz.NewReader(bytes.NewReader(src))
	if err != nil {
		return nil, err
	}
	zr.Multistream(false)
	uncompressed, err := ioutil.ReadAll(io.LimitReader(zr, maxSize+1))
	if err != nil {
		return nil, err
	}
	if len(uncompressed) > int(maxSize) {
		return uncompressed, errors.New("uncompressed file exceeds maxSize")
	}
	return uncompressed, nil
}
