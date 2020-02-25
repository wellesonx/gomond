package models

import "github.com/jinzhu/gorm"

type Log struct {
	gorm.Model
	AppName  string `json:"app_name"`
	Label    string `json:"label"`
	Level    int32  `json:"level"`
	Message  string `json:"message"`
	Hostname string `json:"hostname"`
	Payload  []byte `json:"payload"`
	Line     int32  `json:"line"`
	File     string `json:"file"`
}
