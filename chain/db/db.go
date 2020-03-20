package db

import (
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

	db.AutoMigrate(
		&Metadata{},
		&Event{},
		&ValidatorStatus{},
		&ValidatorVote{},
	)

	return &BandDB{db: db}, nil
}

func (b *BandDB) BeginTransaction() {
	if b.tx != nil {
		panic("BeginTransaction: Cannot begin a new transaction without closing the pending one.")
	}
	b.tx = b.db.Begin()
	if b.tx.Error != nil {
		panic(b.tx.Error)
	}
}

func (b *BandDB) Commit() {
	err := b.tx.Commit().Error
	if err != nil {
		panic(err)
	}
	b.tx = nil
}

func (b *BandDB) RollBack() {
	err := b.tx.Rollback()
	if err != nil {
		panic(err)
	}
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
