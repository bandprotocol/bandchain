package pricecache

import (
	"fmt"

	"github.com/peterbourgon/diskv"

	"github.com/bandprotocol/bandchain/chain/pkg/obi"
)

type Cache struct {
	priceCache *diskv.Diskv
}

type Price struct {
	Multiplier  uint64 `json:"multiplier"`
	Px          uint64 `json:"px"`
	ResolveTime int64  `json:"resolve_time"`
}

func NewPrice(multiplier uint64, px uint64, resolveTime int64) Price {
	return Price{
		Multiplier:  multiplier,
		Px:          px,
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
