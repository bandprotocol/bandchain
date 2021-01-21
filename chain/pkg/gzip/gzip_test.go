package gzip_test

import (
	"bytes"
	gz "compress/gzip"
	"testing"

	"github.com/GeoDB-Limited/odincore/chain/pkg/gzip"
	"github.com/stretchr/testify/require"
)

func TestUncompress(t *testing.T) {
	file1 := []byte("file")
	var buf bytes.Buffer
	zw := gz.NewWriter(&buf)
	zw.Write(file1)
	zw.Close()
	gzipFile := buf.Bytes()

	accFile, err := gzip.Uncompress(gzipFile, 10)
	require.NoError(t, err)
	require.Equal(t, file1, accFile)

	accFile, err = gzip.Uncompress(gzipFile, 2)
	require.Error(t, err)

	_, err = gzip.Uncompress(file1, 999)
	require.Error(t, err)
}

func TestIsGzip(t *testing.T) {
	file1 := []byte("file")
	var buf bytes.Buffer
	zw := gz.NewWriter(&buf)
	zw.Write(file1)
	zw.Close()
	gzipFile := buf.Bytes()
	require.True(t, gzip.IsGzipped(gzipFile))
	require.False(t, gzip.IsGzipped(file1))
}
