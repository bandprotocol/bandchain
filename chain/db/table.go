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
	gorm.Model
	ValidatorAddress string
	ElectedCount     uint
	VotedCount       uint
	MissedCount      uint
}

type ValidatorVote struct {
	ValidatorAddress string `gorm:"primary_key"`
	Block            int64  `gorm:"primary_key;auto_increment:false"`
	Voted            bool
}
