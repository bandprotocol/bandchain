package db

import (
	"time"
)

func (b *BandDB) AddBlock(
	height int64,
	timestamp time.Time,
	proposer string,
	blockHash []byte,
) error {
	return b.tx.Create(&Block{
		Height:    height,
		Timestamp: timestamp.UnixNano() / 1000000, // millisecond
		Proposer:  proposer,
		BlockHash: blockHash,
	}).Error
}
