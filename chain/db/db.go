package db

import (
	"fmt"

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

func NewDB(dialect, path string) (*BandDB, error) {
	db, err := gorm.Open(dialect, path)
	db.CreateTable(Event{})

	if err != nil {
		return nil, err
	}
	return &BandDB{db: db}, nil
}

func (b *BandDB) HandleEvent(eventName string, attributes map[string]string) {
	switch eventName {
	case "message":
		{
			b.db.Create(&Event{Name: attributes["action"]})
		}
	default:
		// TODO: Better logging
		fmt.Println("There isn't event handler for this type")
	}
}
