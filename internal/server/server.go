package server

import (
	"Go-IM/internal/handler"
	"Go-IM/internal/model"
	"Go-IM/pkg/common/customtypes"
	"Go-IM/pkg/err"
	"Go-IM/pkg/zaplog"
	"github.com/gofiber/contrib/fiberzap/v2"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type FiberServer struct {
	*fiber.App
	*handler.Handlers
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: os.Getenv("NAME"),
			AppName:      os.Getenv("NAME"),
			ErrorHandler: err.HandleError,
			Prefork:      true,
			ReadTimeout:  30000,
			WriteTimeout: 30000,
			IdleTimeout:  30000,
		}),
		Handlers: handler.NewHandlers(),
	}
	server.Use(recover2.New(recover2.Config{
		Next:              nil,
		EnableStackTrace:  false,
		StackTraceHandler: nil,
	}))
	server.Use(fiberzap.New(fiberzap.Config{
		Logger: zaplog.Logger,
		Levels: []zapcore.Level{zap.DebugLevel, zap.InfoLevel, zap.InfoLevel, zap.WarnLevel, zap.ErrorLevel},
	}))
	server.Use(swagger.New(swagger.Config{
		Next:     nil,
		BasePath: "/",
		FilePath: "./docs/swagger.json",
		Path:     "docs",
		Title:    "Go-IM",
		CacheAge: 3600,
	}))
	server.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: "HS256",
			Key:    []byte(os.Getenv("JWT_SECRET_KEY")),
		},
		ErrorHandler: err.JwtErrorHandler,
		Filter: func(ctx *fiber.Ctx) bool {
			return ctx.Path() == "/api/account/getcaptcha" || ctx.Path() == "/api/account/register" ||
				ctx.Path() == "/api/account/login"
		},
		ContextKey:  "UserInfo",
		TokenLookup: "header:Authorization",
		AuthScheme:  "Bearer",
		Claims:      &customtypes.GIClaims{},
	}))
	server.InitDBTables(&model.User{}, &model.Group{}, &model.GroupMember{}, &model.Message{}, &model.UserFriend{})
	return server
}
