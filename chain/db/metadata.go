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

func (b *BandDB) SaveChainID(chainID string) error {
	return b.tx.Where(Metadata{Key: KeyChainID}).
		Assign(Metadata{Value: chainID}).
		FirstOrCreate(&Metadata{}).Error
}

func (b *BandDB) ValidateChainID(chainID string) error {
	var chainIDRow Metadata
	err := b.tx.Where(Metadata{Key: KeyChainID}).First(&chainIDRow).Error
	if err != nil {
		return err
	}
	if chainIDRow.Value != chainID {
		return errors.New("Chain id not match")
	}
	return nil
}

func (b *BandDB) SetLastProcessedHeight(height int64) error {
	return b.tx.Where(Metadata{Key: KeyLastProcessedHeight}).
		Assign(Metadata{Value: fmt.Sprintf("%d", height)}).
		FirstOrCreate(&Metadata{}).Error
}

func (b *BandDB) GetLastProcessedHeight() (int64, error) {
	var heightRow Metadata
	err := b.tx.Where(Metadata{Key: KeyLastProcessedHeight}).
		First(&heightRow).Error
	if err != nil {
		return 0, err
	}

	height, err := strconv.ParseInt(heightRow.Value, 10, 64)
	if err != nil {
		return 0, err
	}
	return height, nil
}
