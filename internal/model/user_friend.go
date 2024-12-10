package model

import (
	"gorm.io/gorm"
	"gorm.io/plugin/optimisticlock"
)

type UserFriend struct {
	gorm.Model
	UserId   int32 `json:"userId" gorm:"index;comment:'用户ID'"`
	FriendId int32 `json:"friendId" gorm:"index;comment:'好友ID'"`
	Version  optimisticlock.Version
}
