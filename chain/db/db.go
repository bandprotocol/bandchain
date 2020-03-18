package db

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type BandDB struct {
	db *gorm.DB
	tx *gorm.DB
}

func NewDB(dialect, path string) (*BandDB, error) {
	db, err := gorm.Open(dialect, path)
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&Metadata{}, &Event{})

	return &BandDB{db: db}, nil
}

func (b *BandDB) SaveChainID(chainID string) {
	chainIDRow := Metadata{Key: "chain-id", Value: chainID}
	b.db.Where(Metadata{Key: "chain-id"}).Assign(Metadata{Value: chainID}).FirstOrCreate(&chainIDRow)
}

func (b *BandDB) ValidateChainID(chainID string) error {
	var chainIDRow Metadata
	b.db.Where("key = ?", "chain-id").First(&chainIDRow)
	if chainIDRow.Value != chainID {
		return errors.New("Chain id not match")
	}
	return nil
}

func (b *BandDB) OpenTransaction() {
	if b.tx != nil {
		panic("There is an transaction that left open")
	}
	b.tx = b.db.Begin()
}

func (b *BandDB) Commit() {
	b.tx.Commit()
	b.tx = nil
}

func (b *BandDB) RollBack() {
	b.tx.Rollback()
	b.tx = nil
}

func (b *BandDB) HandleEvent(eventName string, attributes map[string]string) {
	switch eventName {
	// Just proof of concept
	case "message":
		{
			// Event message split events on report event eg.
			// message map[action:report]
			// message map[sender:band17xpfvakm2amg962yls6f84z3kell8c5lfkrzn4]
			action, ok := attributes["action"]
			if ok {
				b.handleMessageEvent(action)
			}
		}
	default:
		// TODO: Better logging
		fmt.Println("There isn't event handler for this type")
	}
}
