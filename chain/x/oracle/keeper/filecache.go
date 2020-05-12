package keeper

import (
	"crypto/sha256"
	"encoding/hex"
)

// AddFile saves the given data to a file in HOME/files directory using sha256 sum as filename.
func (k Keeper) AddFile(data []byte) string {
	hash := sha256.Sum256(data)
	fileName := hex.EncodeToString(hash[:])
	if !k.fileCache.Has(fileName) {
		k.fileCache.Write(fileName, data)
	}
	return fileName
}

// GetFile loads the file from the file storage. Panics if the file does not exist.
func (k Keeper) GetFile(fileName string) []byte {
	data, err := k.fileCache.Read(fileName)
	if err != nil {
		panic(err)
	}
	// We do not perform integrity check here for performance optimization.
	return data
}
