package models

import "github.com/jinzhu/gorm"

type Event int

const (
	DISCONNECT Event = iota + 1
	CONNECT
	ERROR
)

type WatcherEvent struct {
	gorm.Model
	Watcher   Watcher
	WatcherID uint
	Event     Event
	Message   string
}
