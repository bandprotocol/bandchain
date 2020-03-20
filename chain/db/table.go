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
