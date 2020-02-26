package models

import "github.com/jinzhu/gorm"

type Mode int

const (
	LIVE = iota + 1
	LATEST
)

type Watcher struct {
	gorm.Model `json:"-"`
	App        string `json:"app"`
	Host       string `json:"host" gorm:"unique"`
	Mode       int    `json:"mode"`
	IsActive   bool   `json:"is_active"`
}
