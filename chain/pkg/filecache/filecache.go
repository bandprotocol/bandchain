package filecache

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/peterbourgon/diskv"
)

type Cache struct {
	fileCache *diskv.Diskv
}

func New(basePath string) Cache {
	return Cache{
		fileCache: diskv.New(diskv.Options{
			BasePath:     basePath,
			Transform:    func(s string) []string { return []string{} },
			CacheSizeMax: 32 * 1024 * 1024, // 32MB TODO: Make this configurable
		}),
	}
}

func getFileName(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// AddFile saves the given data to a file in HOME/files directory using sha256 sum as filename.
func (c Cache) AddFile(data []byte) string {
	fileName := getFileName(data)
	if !c.fileCache.Has(fileName) {
		c.fileCache.Write(fileName, data)
	}
	return fileName
}

// GetFile loads the file from the file storage. Panics if the file does not exist.
func (c Cache) GetFile(fileName string) []byte {
	data, err := c.fileCache.Read(fileName)
	if err != nil {
		panic(err)
	}
	if getFileName(data) != fileName { // Perform integrity check for safety. NEVER EXPECT TO HIT.
		panic("Inconsistent filecache content")
	}
	return data
}
