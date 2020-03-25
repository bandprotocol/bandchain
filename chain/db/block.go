package db

import (
	"time"

	"github.com/tendermint/tendermint/libs/common"
)

func (b *BandDB) AddBlock(
	height int64,
	timestamp time.Time,
	proposer common.HexBytes,
	blockHash []byte,
) error {
	return b.tx.Create(&Block{
		Height:    height,
		Timestamp: timestamp.UnixNano() / int64(time.Millisecond),
		Proposer:  proposer.String(),
		BlockHash: blockHash,
	}).Error
}
