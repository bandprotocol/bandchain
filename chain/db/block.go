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
		Timestamp: timestamp.UnixNano() / int64(time.Millisecond), // millisecond
		Proposer:  proposer,
		BlockHash: blockHash,
	}).Error
}
