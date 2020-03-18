package db

import (
	"github.com/jinzhu/gorm"
)

type Event struct {
	gorm.Model
	Name string
}
