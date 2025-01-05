package handler

import (
	"Go-IM/internal/model"
	"Go-IM/pkg/captcha"
	"Go-IM/pkg/common/customtypes"
	"Go-IM/pkg/common/defines"
	"Go-IM/pkg/common/request"
	"Go-IM/pkg/common/resp"
	"Go-IM/pkg/err"
	"Go-IM/pkg/giutils"
	"Go-IM/pkg/validates"
	"Go-IM/pkg/zaplog"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"os"
	"time"
)

// GetCaptcha godoc
//
//	@Summary		Get Captcha
//	@Description	get captcha
//	@Tags			账户管理
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
//	@Failure		1003	{object}	resp.Response		"注册失败"
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
	e = h.db.GetDB(c).Model(&user).Where("username = ? or email = ?", register.UserName, register.Email).First(&user).Error
	if e == nil {
		return err.UserExists
	}
	user.Uuid = uuid.New().String()
	user.Email = register.Email
	user.Username = register.UserName
	user.Password = giutils.GenerateHashPass(register.Password)
	e = h.db.GetDB(c).Create(&user).Error
	if e != nil {
		return e
	}
	return c.JSON(resp.Success(0, "注册成功", nil))
}

// LoginHandler godoc
// @Summary		用户登录
// @Description	用户登录接口，接收用户登录信息并返回登录结果
// @Tags			账户管理
// @Accept			json
// @Produce		json
// @Param			login	body		request.Login	true	"登录信息"
// @Success		200	{object}	resp.Response	"登录成功通知"
// @Failure		1003	{object}	resp.Response		"登录失败"
// @Router			/api/account/login [get]
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
	var user model.User
	if h.db.GetDB(c).Model(&user).Where("email = ?", login.Email).First(&user).Error != nil {
		return err.NotFound
	}
	if !giutils.CompareHashPassword(user.Password, login.Password) {
		return err.PassError
	}
	claims := customtypes.GIClaims{
		UserId: user.Uuid,
		Admin:  false,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(time.Hour * defines.TWOKEN_EXPIRE),
			},
		},
	}
	if giutils.IsAdmin(user.Email) {
		claims.Admin = true
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, e := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if e != nil {
		return e
	}
	if e := h.db.SetAndTime(defines.USER_TOKEN_KEY+user.Uuid, tokenString, defines.USER_TOKEN); e != nil {
		return e
	}
	return c.JSON(resp.Success(0, "登录成功", fiber.Map{"token": tokenString}))
}

// GetUserInfoHandler godoc
// @Summary		获取用户信息
// @Description	获取用户信息接口，返回用户信息
// @Tags			账户管理
// @Accept			json
// @Produce		json
// @Success		200	{object}	resp.Response	"获取用户信息成功"
// @Failure		1003	{object}	resp.Response		"获取用户信息失败"
// @Router			/api/account/getuserinfo [get]
func (h *Handlers) GetUserInfoHandler(c *fiber.Ctx) error {
	user := c.Locals("UserInfo").(*jwt.Token)
	claims := user.Claims.(*customtypes.GIClaims)
	zaplog.Logger.Info("GetUserInfoHandler", zap.Any("userId", claims.UserId))
	if value := h.db.GetValue(defines.USER_TOKEN_KEY + claims.UserId); value == "" {
		return err.TokenInvalid
	}
	var userInfo model.User
	if e := h.db.GetDB(c).Model(&userInfo).Where("uuid = ?", claims.UserId).First(&userInfo).Error; e != nil {
		return err.NotFound
	}
	return c.JSON(resp.Success(0, "获取用户信息成功", userInfo))
}
