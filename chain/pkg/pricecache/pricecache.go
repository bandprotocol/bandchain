package pricecache

import (
	"fmt"

	"github.com/bandprotocol/bandchain/chain/pkg/obi"
	"github.com/peterbourgon/diskv"
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

func GetFilename(symbol string, minCount uint64, askCount uint64) string {
	return fmt.Sprintf("%s,%v,%v", symbol, minCount, askCount)
}

// SetPrice saves the given data to a file in HOME/files directory using symbol,minCount,askCount format as filename.
func (c Cache) SetPrice(symbol string, minCount uint64, askCount uint64, price Price) error {
	data, err := obi.Encode(price)
	if err != nil {
		return err
	}
	return c.priceCache.Write(GetFilename(symbol, minCount, askCount), data)
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
