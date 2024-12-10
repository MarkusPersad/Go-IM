package model

import (
	"gorm.io/gorm"
	"gorm.io/plugin/optimisticlock"
)

type GroupMember struct {
	gorm.Model
	UserId   int32  `json:"userId" gorm:"index;comment:'用户ID'"`
	GroupId  int32  `json:"groupId" gorm:"index;comment:'群组ID'"`
	Nickname string `json:"nickname" gorm:"type:varchar(350);comment:'昵称"`
	Mute     int16  `json:"mute" gorm:"comment:'是否禁言'"`
	Version  optimisticlock.Version
}
