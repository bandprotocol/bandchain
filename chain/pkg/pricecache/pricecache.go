package pricecache

import (
	"fmt"

	"github.com/peterbourgon/diskv"

	"github.com/bandprotocol/bandchain/chain/pkg/obi"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

type Cache struct {
	priceCache *diskv.Diskv
}

type Price struct {
	Symbol      string          `json:"symbol"`
	Multiplier  uint64          `json:"multiplier"`
	Px          uint64          `json:"px"`
	RequestID   types.RequestID `json:"request_id"`
	ResolveTime int64           `json:"resolve_time"`
}

func NewPrice(symbol string, multiplier uint64, px uint64, reqID types.RequestID, resolveTime int64) Price {
	return Price{
		Symbol:      symbol,
		Multiplier:  multiplier,
		Px:          px,
		RequestID:   reqID,
		ResolveTime: resolveTime,
	}
}

// New creates and returns a new file-backed data caching instance.
func New(basePath string) Cache {
	return Cache{
		priceCache: diskv.New(diskv.Options{
			BasePath:     basePath,
			Transform:    func(s string) []string { return []string{} },
			CacheSizeMax: 32 * 1024 * 1024, // 32MB TODO: Make this configurable
		}),
	}
}

// GetFilename returns filename format as symbol,minCount,askCount.
func GetFilename(symbol string, minCount uint64, askCount uint64) string {
	return fmt.Sprintf("%s,%v,%v", symbol, minCount, askCount)
}

// SetPrice saves the given data to a file in HOME/prices directory.
func (c Cache) SetPrice(filename string, price Price) error {
	return c.priceCache.Write(filename, obi.MustEncode(price))
}

// GetPrice loads the file from the file storage. Returns error if the file does not exist.
func (c Cache) GetPrice(filename string) (Price, error) {
	bz, err := c.priceCache.Read(filename)
	if err != nil {
		return Price{}, err
	}
	var price Price
	obi.MustDecode(bz, &price)
	return price, nil
}
