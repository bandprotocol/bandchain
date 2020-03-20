package db

import (
	"github.com/jinzhu/gorm"
)

type Metadata struct {
	Key   string `gorm:"primary_key"`
	Value string
}

type Event struct {
	gorm.Model
	Name string
}

type ValidatorStatus struct {
	OperatorAddress  string `gorm:"primary_key"`
	ConsensusAddress string `gorm:"unique;not null"`
	ElectedCount     uint
	VotedCount       uint
	MissedCount      uint

	ValidatorVotes []ValidatorVote `gorm:"foreignkey:ConsensusAddress;association_foreignkey:ConsensusAddress"`
}

type ValidatorVote struct {
	ConsensusAddress string `gorm:"primary_key"`
	BlockHeight      int64  `gorm:"primary_key;auto_increment:false"`
	Voted            bool
}
