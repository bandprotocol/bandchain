package db

import (
	"github.com/jinzhu/gorm"
)

type Metadata struct {
	Key   string `gorm:"primary_key"`
	Value string
}

func (Metadata) TableName() string {
	return "metatdata"
}

type Event struct {
	gorm.Model
	Name string
}
