package oracle

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
)

// magic bytes to identify gzip.
// See https://www.ietf.org/rfc/rfc1952.txt
// and https://github.com/golang/go/blob/master/src/net/http/sniff.go#L186
var gzipIdent = []byte("\x1F\x8B\x08")

func IsGzip(src []byte) bool {
	return bytes.Equal(gzipIdent, src[0:3])
}

// uncompress returns gzip uncompressed content or given src when not gzip.
func Uncompress(src []byte, maxSize int64) ([]byte, error) {
	zr, err := gzip.NewReader(bytes.NewReader(src))
	if err != nil {
		return nil, err
	}
	zr.Multistream(false)

	uncompressFile, err := ioutil.ReadAll(io.LimitReader(zr, maxSize+1))
	if err != nil {
		return nil, err
	}
	if len(uncompressFile) > int(maxSize) {
		return nil, errors.New("uncompress file exeed maxsize")
	}

	return uncompressFile, nil
}
