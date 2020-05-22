package gzip

import (
	"bytes"
	gz "compress/gzip"
	"errors"
	"io"
	"io/ioutil"
)

// magic bytes to identify gzip.
// See https://www.ietf.org/rfc/rfc1952.txt
// and https://github.com/golang/go/blob/master/src/net/http/sniff.go#L186
var gzipIdent = []byte("\x1F\x8B\x08")

// IsGzipped returns true if the input is gzipped. (remove file)
func IsGzipped(src []byte) bool {
	return bytes.Equal(gzipIdent, src[0:3])
}

// Uncompress returns gzip uncompressed.
// If the file is not gzipped file this function will return an error.
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
		return uncompressed, errors.New("uncompressed file exceed maxSize")
	}

	return uncompressed, nil
}
