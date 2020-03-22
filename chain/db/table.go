package db

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Metadata struct {
	Key   string `gorm:"primary_key"`
	Value string `gorm:"not null"`
}

type Event struct {
	gorm.Model
	Name string
}

type Validator struct {
	OperatorAddress  string `gorm:"primary_key"`
	ConsensusAddress string `gorm:"unique;not null"`
	ElectedCount     uint   `gorm:"not null"`
	VotedCount       uint   `gorm:"not null"`
	MissedCount      uint   `gorm:"not null"`
}

type ValidatorVote struct {
	ConsensusAddress string `gorm:"primary_key"`
	BlockHeight      int64  `gorm:"primary_key;auto_increment:false"`
	Voted            bool   `gorm:"not null"`
}

type DataSource struct {
	ID          int64 `gorm:"primary_key;auto_increment:false"`
	Name        string
	Description string
	Owner       string
	Executable  []byte
	Fee         string
	LastUpdated time.Time
}

type DataSourceRevision struct {
	DataSourceID   int64 `gorm:"primary_key;auto_increment:false"`
	RevisionNumber int64 `gorm:"primary_key"`
	Name           string
	Timestamp      time.Time
	BlockHeight    int64
	TxHash         []byte
}
