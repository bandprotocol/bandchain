package db

import (
	"database/sql"
	"encoding/json"
)

type Metadata struct {
	Key   string `gorm:"primary_key"`
	Value string `gorm:"not null"`
}

type Block struct {
	Height    int64  `gorm:"primary_key;auto_increment:false"`
	Timestamp int64  `gorm:"not null"`
	Proposer  string `gorm:"not null"`
	BlockHash []byte `gorm:"not null"`
}

type Transaction struct {
	TxHash      []byte          `gorm:"primary_key"`
	Timestamp   int64           `gorm:"not null"`
	GasUsed     int64           `gorm:"not null"`
	GasLimit    uint64          `gorm:"not null"`
	GasFee      string          `gorm:"not null"`
	Sender      string          `gorm:"not null"`
	Success     bool            `gorm:"not null"`
	BlockHeight int64           `gorm:"not null"`
	Messages    json.RawMessage `sql:"json;not null"`
}

type Account struct {
	Address       string `gorm:"primary_key"`
	Balance       string `gorm:"not null"`
	UpdatedHeight int64  `gorm:"not null"`
}

type Validator struct {
	OperatorAddress     string `gorm:"primary_key"`
	ConsensusAddress    string `gorm:"unique;not null"`
	ElectedCount        uint   `gorm:"not null"`
	VotedCount          uint   `gorm:"not null"`
	MissedCount         uint   `gorm:"not null"`
	Moniker             string `gorm:"not null"`
	Identity            string `gorm:"not null"`
	Website             string `gorm:"not null"`
	Details             string `gorm:"not null"`
	CommissionRate      string `gorm:"not null"`
	CommissionMaxRate   string `gorm:"not null"`
	CommissionMaxChange string `gorm:"not null"`
	MinSelfDelegation   string `gorm:"not null"`
	Jailed              bool   `gorm:"not null"`
	Tokens              string `gorm:"not null"`
	DelegatorShares     string `gorm:"not null"`
}

type ValidatorVote struct {
	ConsensusAddress string `gorm:"primary_key"`
	BlockHeight      int64  `gorm:"primary_key;auto_increment:false"`
	Voted            bool   `gorm:"not null"`
}

type Delegation struct {
	DelegatorAddress string `gorm:"primary_key"`
	ValidatorAddress string `gorm:"primary_key"`
	Shares           string `gorm:"not null"`
}

type DataSource struct {
	ID          int64  `gorm:"primary_key;auto_increment:false"`
	Name        string `gorm:"not null"`
	Description string `gorm:"not null"`
	Owner       string `gorm:"not null"`
	Executable  []byte `gorm:"not null"`
	Fee         string `gorm:"not null"`
	LastUpdated int64  `gorm:"not null"`
}

type DataSourceRevision struct {
	DataSourceID   int64  `gorm:"primary_key;auto_increment:false"`
	RevisionNumber int64  `gorm:"primary_key"`
	Name           string `gorm:"not null"`
	Timestamp      int64  `gorm:"not null"`
	BlockHeight    int64  `gorm:"not null"`
	TxHash         []byte `sql:"default:null"`
}

type OracleScript struct {
	ID          int64  `gorm:"primary_key;auto_increment:false"`
	Name        string `gorm:"not null"`
	Description string `gorm:"not null"`
	Owner       string `gorm:"not null"`
	CodeHash    []byte `gorm:"not null"`
	LastUpdated int64  `gorm:"not null"`
}

type OracleScriptRevision struct {
	OracleScriptID int64  `gorm:"primary_key;auto_increment:false"`
	RevisionNumber int64  `gorm:"primary_key"`
	Name           string `gorm:"not null"`
	Timestamp      int64  `gorm:"not null"`
	BlockHeight    int64  `gorm:"not null"`
	TxHash         []byte `sql:"default:null"`
}

type OracleScriptCode struct {
	CodeHash []byte         `gorm:"primary_key"`
	CodeText sql.NullString `sql:"default:null"`
	Schema   sql.NullString `sql:"default:null"`
}

type RelatedDataSources struct {
	DataSourceID   int64 `gorm:"primary_key;auto_increment:false"`
	OracleScriptID int64 `gorm:"primary_key;auto_increment:false"`
}

type Request struct {
	ID                       int64  `gorm:"primary_key;auto_increment:false"`
	OracleScriptID           int64  `gorm:"not null"`
	Calldata                 []byte `gorm:"not null"`
	SufficientValidatorCount int64  `gorm:"not null"`
	ExpirationHeight         int64  `gorm:"not null"`
	ResolveStatus            string `gorm:"not null"`
	Requester                string `gorm:"not null"`
	TxHash                   []byte `gorm:"not null"`
	Result                   []byte `sql:"default:null"`
}

type RequestedValidator struct {
	RequestID        int64  `gorm:"primary_key;auto_increment:false"`
	ValidatorAddress string `gorm:"primary_key"`
}

type RawDataRequests struct {
	RequestID    int64  `gorm:"primary_key;auto_increment:false"`
	ExternalID   int64  `gorm:"primary_key;auto_increment:false"`
	DataSourceID int64  `gorm:"not null"`
	Calldata     []byte `gorm:"not null"`
}

type Report struct {
	RequestID int64  `gorm:"primary_key;auto_increment:false"`
	Validator string `gorm:"primary_key"`
	TxHash    []byte `gorm:"not null"`
	Reporter  string `gorm:"not null"`
}

type ReportDetail struct {
	RequestID    int64  `gorm:"primary_key;auto_increment:false"`
	Validator    string `gorm:"primary_key"`
	ExternalID   int64  `gorm:"primary_key;auto_increment:false"`
	DataSourceID int64  `gorm:"not null"`
	Data         []byte `gorm:"not null"`
	Exitcode     uint8  `gorm:"not null"`
}
