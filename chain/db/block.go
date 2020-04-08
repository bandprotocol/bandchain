package db

import (
	"time"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

func (b *BandDB) AddBlock(
	height int64,
	timestamp time.Time,
	proposer tmbytes.HexBytes,
	blockHash []byte,
) error {
	return b.tx.Create(&Block{
		Height:    height,
		Timestamp: timestamp.UnixNano() / int64(time.Millisecond),
		Proposer:  proposer.String(),
		BlockHash: blockHash,
	}).Error
}
