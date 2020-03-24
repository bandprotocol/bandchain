package db

import (
	"time"

	"github.com/tendermint/tendermint/libs/common"
)

func (b *BandDB) AddBlock(
	height int64,
	timestamp time.Time,
	proposer string,
	blockHash common.HexBytes,
) error {
	return b.tx.Create(&Block{
		Height:    height,
		Timestamp: timestamp.Unix(),
		Proposer:  proposer,
		BlockHash: blockHash.String(),
	}).Error
}
