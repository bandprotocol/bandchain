package requestcache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/GeoDB-Limited/odincore/chain/x/oracle/types"
	"github.com/peterbourgon/diskv"
)

type Cache struct {
	kv *diskv.Diskv
}

// New creates and returns a new file-backed data caching instance.
func New(basePath string) Cache {
	return Cache{
		kv: diskv.New(diskv.Options{
			BasePath:     basePath,
			Transform:    func(s string) []string { return []string{} },
			CacheSizeMax: 32 * 1024 * 1024, // 32MB TODO: Make this configurable
		}),
	}
}

func getFilename(oid types.OracleScriptID, calldata []byte, askCount uint64, minCount uint64) string {
	full := fmt.Sprintf("%d,%x,%d,%d", oid, calldata, askCount, minCount)
	hash := sha256.Sum256([]byte(full))
	return hex.EncodeToString(hash[:])
}

// SaveLatestRequest saves the latest request id to a file with key that combined from event attributes.
func (c Cache) SaveLatestRequest(oid types.OracleScriptID, calldata []byte, askCount uint64, minCount uint64, reqID types.RequestID) error {
	bz, err := json.Marshal(reqID)
	if err != nil {
		return err
	}
	return c.kv.Write(getFilename(oid, calldata, askCount, minCount), bz)
}

// GetLatestRequest loads the latest request from the file storage. Returns error if the file does not exist.
func (c Cache) GetLatestRequest(oid types.OracleScriptID, calldata []byte, askCount uint64, minCount uint64) (types.RequestID, error) {
	bz, err := c.kv.Read(getFilename(oid, calldata, askCount, minCount))
	if err != nil {
		return 0, err
	}
	var id types.RequestID
	err = json.Unmarshal(bz, &id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
