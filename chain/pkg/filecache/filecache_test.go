package filecache_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/GeoDB-Limited/odincore/chain/pkg/filecache"
)

func TestAddFile(t *testing.T) {
	dir, err := ioutil.TempDir("", "filecache")
	if err != nil {
		panic(err)
	}
	defer func() {
		err := os.RemoveAll(dir)
		if err != nil {
			panic(err)
		}
	}()

	f := filecache.New(dir)
	filename := f.AddFile([]byte("HELLO_WORLD"))
	require.Equal(t, filename, "6f9b514093848217355d76365df1f54f42bdfd5f4e5f54a654c46b493d162c39")

	content, err := ioutil.ReadFile(filepath.Join(dir, filename))
	require.NoError(t, err)
	require.Equal(t, content, []byte("HELLO_WORLD"))

	filename2 := f.AddFile([]byte("HELLO_WORLD"))
	require.Equal(t, filename2, "6f9b514093848217355d76365df1f54f42bdfd5f4e5f54a654c46b493d162c39")
}

func TestMustGetFileOK(t *testing.T) {
	dir, err := ioutil.TempDir("", "filecache")
	if err != nil {
		panic(err)
	}
	defer func() {
		err := os.RemoveAll(dir)
		if err != nil {
			panic(err)
		}
	}()

	f := filecache.New(dir)
	filename := f.AddFile([]byte("BAND"))
	require.Equal(t, filename, "52f1b54ce34b64a02f9946b29f670a12933152b1122514ea969a91c211aa32fc")

	content := f.MustGetFile(filename)
	require.Equal(t, content, []byte("BAND"))
}

func TestGetFileOK(t *testing.T) {
	dir, err := ioutil.TempDir("", "filecache")
	if err != nil {
		panic(err)
	}
	defer func() {
		err := os.RemoveAll(dir)
		if err != nil {
			panic(err)
		}
	}()

	f := filecache.New(dir)
	filename := f.AddFile([]byte("BAND"))
	require.Equal(t, filename, "52f1b54ce34b64a02f9946b29f670a12933152b1122514ea969a91c211aa32fc")

	content, err := f.GetFile(filename)
	require.NoError(t, err)
	require.Equal(t, content, []byte("BAND"))
}

func TestMustGetFileNotExist(t *testing.T) {
	dir, err := ioutil.TempDir("", "filecache")
	if err != nil {
		panic(err)
	}
	defer func() {
		err := os.RemoveAll(dir)
		if err != nil {
			panic(err)
		}
	}()

	f := filecache.New(dir)
	require.Panics(t, func() {
		_ = f.MustGetFile("52f1b54ce34b64a02f9946b29f670a12933152b1122514ea969a91c211aa32fc")
	})
}

func TestGetFileNotExist(t *testing.T) {
	dir, err := ioutil.TempDir("", "filecache")
	if err != nil {
		panic(err)
	}
	defer func() {
		err := os.RemoveAll(dir)
		if err != nil {
			panic(err)
		}
	}()

	f := filecache.New(dir)
	_, err = f.GetFile("52f1b54ce34b64a02f9946b29f670a12933152b1122514ea969a91c211aa32fc")
	require.Error(t, err)
}

func TestMustGetFileGoodContent(t *testing.T) {
	dir, err := ioutil.TempDir("", "filecache")
	if err != nil {
		panic(err)
	}
	defer func() {
		err := os.RemoveAll(dir)
		if err != nil {
			panic(err)
		}
	}()

	f := filecache.New(dir)
	filename := "b20727a9b7cc4198d8785b0ef1fa4c774eb9a360e1563dd4f095ddc7af02bd55" // Correct
	filepath := filepath.Join(dir, filename)
	err = ioutil.WriteFile(filepath, []byte("NOT_LIKE_THIS"), 0666)
	require.NoError(t, err)

	content := f.MustGetFile(filename)
	require.Equal(t, content, []byte("NOT_LIKE_THIS"))
}

func TestGetFileGoodContent(t *testing.T) {
	dir, err := ioutil.TempDir("", "filecache")
	if err != nil {
		panic(err)
	}
	defer func() {
		err := os.RemoveAll(dir)
		if err != nil {
			panic(err)
		}
	}()

	f := filecache.New(dir)
	filename := "b20727a9b7cc4198d8785b0ef1fa4c774eb9a360e1563dd4f095ddc7af02bd55" // Correct
	filepath := filepath.Join(dir, filename)
	err = ioutil.WriteFile(filepath, []byte("NOT_LIKE_THIS"), 0666)
	require.NoError(t, err)

	content, err := f.GetFile(filename)
	require.NoError(t, err)
	require.Equal(t, content, []byte("NOT_LIKE_THIS"))
}

func TestMustGetFileBadContent(t *testing.T) {
	dir, err := ioutil.TempDir("", "filecache")
	if err != nil {
		panic(err)
	}
	defer func() {
		err := os.RemoveAll(dir)
		if err != nil {
			panic(err)
		}
	}()

	f := filecache.New(dir)
	filename := "b20727a9b7cc4198d8785b0ef1fa4c774eb9a360e1563dd4f095ddc7af02bd56" // Not correct
	filepath := filepath.Join(dir, filename)
	err = ioutil.WriteFile(filepath, []byte("NOT_LIKE_THIS"), 0666)
	require.NoError(t, err)

	require.Panics(t, func() {
		_ = f.MustGetFile(filename)
	})
}

func TesGetFileBadContent(t *testing.T) {
	dir, err := ioutil.TempDir("", "filecache")
	if err != nil {
		panic(err)
	}
	defer func() {
		err := os.RemoveAll(dir)
		if err != nil {
			panic(err)
		}
	}()

	f := filecache.New(dir)
	filename := "b20727a9b7cc4198d8785b0ef1fa4c774eb9a360e1563dd4f095ddc7af02bd56" // Not correct
	filepath := filepath.Join(dir, filename)
	err = ioutil.WriteFile(filepath, []byte("NOT_LIKE_THIS"), 0666)
	require.NoError(t, err)

	_, err = f.GetFile(filename)
	require.Error(t, err)
}

func TestMustGetFileInconsistentContent(t *testing.T) {
	dir, err := ioutil.TempDir("", "filecache")
	if err != nil {
		panic(err)
	}
	defer func() {
		err := os.RemoveAll(dir)
		if err != nil {
			panic(err)
		}
	}()

	f := filecache.New(dir)
	filename := "b20727a9b7cc4198d8785b0ef1fa4c774eb9a360e1563dd4f095ddc7af02bd55"
	filepath := filepath.Join(dir, filename)
	err = ioutil.WriteFile(filepath, []byte("INCONSISTENT"), 0666) // Not consistent with name
	require.NoError(t, err)
	require.Panics(t, func() {
		_ = f.MustGetFile(filename)
	})
}

func TestGetFileInconsistentContent(t *testing.T) {
	dir, err := ioutil.TempDir("", "filecache")
	if err != nil {
		panic(err)
	}
	defer func() {
		err := os.RemoveAll(dir)
		if err != nil {
			panic(err)
		}
	}()

	f := filecache.New(dir)
	filename := "b20727a9b7cc4198d8785b0ef1fa4c774eb9a360e1563dd4f095ddc7af02bd55"
	filepath := filepath.Join(dir, filename)
	err = ioutil.WriteFile(filepath, []byte("INCONSISTENT"), 0666) // Not consistent with name
	_, err = f.GetFile(filename)
	require.Error(t, err)
}
