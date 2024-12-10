package handler

import (
	"Go-IM/internal/model"
	"Go-IM/pkg/captcha"
	"Go-IM/pkg/common/request"
	"Go-IM/pkg/common/resp"
	"Go-IM/pkg/err"
	"Go-IM/pkg/giutils"
	"Go-IM/pkg/validates"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetCaptcha godoc
//
//	@Summary		Get Captcha
//	@Description	get captcha
//	@Tags			account
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	customtypes.CaptDataBase64
//	@Router			/api/account/getcaptcha [get]
func (h *Handlers) GetCaptcha(c *fiber.Ctx) error {
	capt := captcha.New(h.db)
	dataBase64, e := capt.Generate()
	if e != nil {
		return err.CheckCode
	}
	return c.JSON(dataBase64)
}

// RegisterHandler godoc
//
//	@Summary		用户注册
//	@Description	用户注册接口，接收用户注册信息并返回注册结果
//	@Tags			账户管理
//	@Accept			json
//	@Produce		json
//	@Param			register	body		request.Register	true	"注册信息"
//	@Success		200			{object}	resp.Response	"注册成功通知"
//	@Failure		1003	{object}	error		"注册失败"
//	@Router			/api/account/register [get]
func (h *Handlers) RegisterHandler(c *fiber.Ctx) error {
	var register request.Register
	if e := c.BodyParser(&register); e != nil {
		return err.BadRequest
	}
	e := validates.Validatec.Validate(&register)
	if e != nil {
		return e
	}
	capt := captcha.New(h.db)
	e = capt.Verify(register.CheckCodeKey, register.CheckCode, true)
	if e != nil {
		return e
	}
	var user model.User
	e = h.db.GetDB().Model(&user).Where("username = ? or email = ?", register.UserName, register.Email).First(&user).Error
	if e == nil {
		return err.UserExists
	}
	user.Uuid = uuid.New().String()
	user.Email = register.Email
	user.Username = register.UserName
	user.Password = giutils.GenerateHashPass(register.Password)
	e = h.db.GetDB().Create(&user).Error
	if e != nil {
		return e
	}
	return c.JSON(resp.Success(0, "注册成功", nil))
}

func (h *Handlers) LoginHandler(c *fiber.Ctx) error {
	var login request.Login
	if e := c.BodyParser(&login); e != nil {
		return err.BadRequest
	}
	if e := validates.Validatec.Validate(&login); e != nil {
		return e
	}
	capt := captcha.New(h.db)
	if e := capt.Verify(login.CheckCodeKey, login.CheckCode, true); e != nil {
		return e
	}

	return nil
}