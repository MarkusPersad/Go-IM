package model

import (
	"gorm.io/gorm"
	"gorm.io/plugin/optimisticlock"
)

type Group struct {
	gorm.Model
	Uuid    string `json:"uuid" gorm:"type:varchar(150);not null;unique_index:idx_uuid;comment:'uuid'"`
	UserId  int32  `json:"userId" gorm:"index;comment:'群主ID'"`
	Name    string `json:"name" gorm:"type:varchar(150);comment:'群名称"`
	Notice  string `json:"notice" gorm:"type:varchar(350);comment:'群公告"`
	Version optimisticlock.Version
}
