package models

import (
	"github.com/jinzhu/gorm"
)

type Event int

const (
	DISCONNECT Event = iota + 1
	CONNECT
	ERROR
)

type WatcherEvent struct {
	gorm.Model
	Watcher   Watcher `json:"watcher"`
	WatcherID uint    `json:"-"`
	Event     Event   `json:"event"`
	Message   string  `json:"message"`
}
