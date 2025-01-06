package handler

import (
	"Go-IM/internal/model"
	"Go-IM/pkg/common/customtypes"
	"Go-IM/pkg/common/defines"
	"Go-IM/pkg/common/request"
	"Go-IM/pkg/common/resp"
	"Go-IM/pkg/err"
	"Go-IM/pkg/validates"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func (h *Handlers) AddFriend(c *fiber.Ctx) error {
	token := c.Locals("UserInfo").(*jwt.Token)
	claims := token.Claims.(*customtypes.GIClaims)
	if h.db.GetValue(defines.USER_TOKEN_KEY+claims.UserId) == "" {
		return err.Timeout
	}
	var addFriend request.AddFriend
	if e := c.BodyParser(&addFriend); e != nil {
		return err.BadRequest
	}
	if e := validates.New().Validate(&addFriend); e != nil {
		return e
	}
	if e := h.db.GetDB(c).Transaction(func(tx *gorm.DB) error {
		var friend model.User
		if e := tx.Model(&friend).Where("username = ?", addFriend.FriendInfo).Or("email = ?", addFriend.FriendInfo).First(&friend).Error; e != nil {
			return err.NotFound
		}
		if e := createFriendShip(tx, claims.UserId, friend.Uuid); e != nil {
			return e
		}
		if e := createFriendShip(tx, friend.Uuid, claims.UserId); e != nil {
			return e
		}
		//TODO 发送打招呼消息
		return nil
	}); e != nil {
		return e
	}
	return c.JSON(resp.Success(0, "添加成功", nil))
}

func createFriendShip(tx *gorm.DB, userId, friendId string) error {
	var userFriend model.UserFriend
	userFriend.UserId = userId
	userFriend.FriendId = friendId
	if e := tx.Model(&userFriend).Where("userid = ? AND friendid = ?", userFriend.UserId, userFriend.FriendId).First(&userFriend).Error; e == nil {
		return err.AlreadyExists
	}
	if e := tx.Create(&userFriend).Error; e != nil {
		return e
	}
	return nil
}
func sendMessage() error {
	//TODO 发送消息
	return nil
}
