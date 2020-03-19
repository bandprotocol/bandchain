package db

import (
	"errors"
	"fmt"
	"strconv"
)

const (
	KeyChainID             = "chain_id"
	KeyLastProcessedHeight = "last_processed_height"
)

func (b *BandDB) SaveChainID(chainID string) {
	var chainIDRow Metadata
	b.tx.Where(Metadata{Key: KeyChainID}).Assign(Metadata{Value: chainID}).FirstOrCreate(&chainIDRow)
}

func (b *BandDB) ValidateChainID(chainID string) error {
	var chainIDRow Metadata
	b.tx.Where(Metadata{Key: KeyChainID}).First(&chainIDRow)
	if chainIDRow.Value != chainID {
		return errors.New("Chain id not match")
	}
	return nil
}

func (b *BandDB) GetBlockHeight() int64 {
	var heightRow Metadata
	b.tx.Where(
		Metadata{Key: KeyLastProcessedHeight},
	).First(&heightRow)

	height, err := strconv.ParseInt(heightRow.Value, 10, 64)
	if err != nil {
		panic(err)
	}
	return height
}

func (b *BandDB) SetBlockHeight(height int64) {
	var heightRow Metadata
	b.tx.Where(
		Metadata{Key: KeyLastProcessedHeight},
	).Assign(
		Metadata{Value: fmt.Sprintf("%d", height)},
	).FirstOrCreate(&heightRow)
}
