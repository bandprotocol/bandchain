package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type BandDB struct {
	db *gorm.DB
}

type Event struct {
	gorm.Model
	Name string
}

func NewDB(path string) (*BandDB, error) {
	db, err := gorm.Open("sqlite3", path)
	db.CreateTable(Event{})

	if err != nil {
		return nil, err
	}
	return &BandDB{db: db}, nil
}

func (b *BandDB) HandleEvent(eventName string) {
	b.db.Create(&Event{Name: eventName})
}
