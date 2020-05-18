package oracle

import (
	"bytes"
	"compress/gzip"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUncompress(t *testing.T) {
	file1 := []byte("file")
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	zw.Write(file1)
	zw.Close()
	gzipFile := buf.Bytes()

	accFile, err := Uncompress(gzipFile, 10)
	require.Nil(t, err)
	require.Equal(t, file1, accFile)

	accFile, err = Uncompress(gzipFile, 2)
	require.NotNil(t, err)

	_, err = Uncompress(file1, 999)
	require.NotNil(t, err)
}

func TestIsGzip(t *testing.T) {
	file1 := []byte("file")
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	zw.Write(file1)
	zw.Close()
	gzipFile := buf.Bytes()

	require.True(t, IsGzip(gzipFile))
	require.False(t, IsGzip(file1))

}
