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
	KeyInflationRate          = "inflation_rate"
)

func (b *BandDB) GetMetadataValue(key string) (string, error) {
	var data Metadata
	err := b.tx.Where(Metadata{Key: key}).First(&data).Error
	if err != nil {
		return "", err
	}
	return data.Value, nil
}

func (b *BandDB) SetMetadataValue(key, value string) error {
	return b.tx.Where(Metadata{Key: key}).
		Assign(Metadata{Value: value}).
		FirstOrCreate(&Metadata{}).Error
}

func (b *BandDB) SetMetadataValueInt64(key string, value int64) error {
	return b.SetMetadataValue(key, fmt.Sprintf("%d", value))
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
	return b.SetMetadataValueInt64(KeyLastProcessedHeight, height)
}

func (b *BandDB) GetLastProcessedHeight() (int64, error) {
	return b.GetMetadataValueInt64(KeyLastProcessedHeight)
}

func (b *BandDB) SetUptimeLookBackDuration(duration int64) error {
	return b.SetMetadataValueInt64(KeyUptimeLookBackDuration, duration)
}

func (b *BandDB) GetUptimeLookBackDuration() (int64, error) {
	return b.GetMetadataValueInt64(KeyUptimeLookBackDuration)
}

func (b *BandDB) SetInflationRate(inflationRate string) error {
	return b.SetMetadataValue(KeyInflationRate, inflationRate)
}

func (b *BandDB) GetInflationRate() (string, error) {
	rawString, err := b.GetMetadataValue(KeyInflationRate)
	if err != nil {
		return "", err
	}
	return rawString, nil
}
