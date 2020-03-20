package db

import (
	"errors"
	"fmt"
	"strconv"
)

const (
	KeyChainID                = "chain_id"
	KeyLastProcessedHeight    = "last_processed_height"
	KeyUptimeLookBackDuration = "uptime_look_back_duration"
)

func (b *BandDB) GetMetadataValue(key string) (string, error) {
	var data Metadata
	err := b.tx.Where(Metadata{Key: key}).First(&data).Error
	if err != nil {
		return "", err
	}
	return data.Value, nil
}

func (b *BandDB) GetMetadataValueInt64(key string) (int64, error) {
	rawString, err := b.GetMetadataValue(key)
	if err != nil {
		return 0, err
	}

	value, err := strconv.ParseInt(rawString, 10, 64)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func (b *BandDB) SetMetadataValue(key, value string) error {
	return b.tx.Where(Metadata{Key: key}).
		Assign(Metadata{Value: value}).
		FirstOrCreate(&Metadata{}).Error
}

func (b *BandDB) SaveChainID(chainID string) error {
	return b.SetMetadataValue(KeyChainID, chainID)
}

func (b *BandDB) ValidateChainID(chainID string) error {
	chainIDDB, err := b.GetMetadataValue(KeyChainID)
	if err != nil {
		return err
	}
	if chainIDDB != chainID {
		return errors.New("Chain id not match")
	}
	return nil
}

func (b *BandDB) SetLastProcessedHeight(height int64) error {
	return b.SetMetadataValue(KeyLastProcessedHeight, fmt.Sprintf("%d", height))
}

func (b *BandDB) GetLastProcessedHeight() (int64, error) {
	return b.GetMetadataValueInt64(KeyLastProcessedHeight)
}

func (b *BandDB) SetUptimeLookBackDuration(duration int64) error {
	return b.SetMetadataValue(KeyUptimeLookBackDuration, fmt.Sprintf("%d", duration))
}

func (b *BandDB) GetUptimeLookBackDuration() (int64, error) {
	return b.GetMetadataValueInt64(KeyUptimeLookBackDuration)
}
