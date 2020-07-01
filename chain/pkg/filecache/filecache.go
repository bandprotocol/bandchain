package filecache

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/peterbourgon/diskv"
)

type Cache struct {
	fileCache *diskv.Diskv
}

// New creates and returns a new file-backed data caching instance.
func New(basePath string) Cache {
	return Cache{
		fileCache: diskv.New(diskv.Options{
			BasePath:     basePath,
			Transform:    func(s string) []string { return []string{} },
			CacheSizeMax: 32 * 1024 * 1024, // 32MB TODO: Make this configurable
		}),
	}
}

func getFilename(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// AddFile saves the given data to a file in HOME/files directory using sha256 sum as filename.
func (c Cache) AddFile(data []byte) string {
	filename := getFilename(data)
	if !c.fileCache.Has(filename) {
		c.fileCache.Write(filename, data)
	}
	return filename
}

// GetFile loads the file from the file storage. Returns error if the file does not exist.
func (c Cache) GetFile(filename string) ([]byte, error) {
	data, err := c.fileCache.Read(filename)
	if err != nil {
		return nil, err
	}
	if getFilename(data) != filename { // Perform integrity check for safety. NEVER EXPECT TO HIT.
		return nil, errors.New("Inconsistent filecache content")
	}
	return data, nil
}

// MustGetFile loads the file from the file storage. Panics if the file does not exist.
func (c Cache) MustGetFile(filename string) []byte {
	data, err := c.GetFile(filename)
	if err != nil {
		panic(err)
	}
	return data
}
