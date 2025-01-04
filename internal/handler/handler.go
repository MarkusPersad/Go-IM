package handler

import (
	"Go-IM/internal/database"
	"Go-IM/pkg/zaplog"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Handlers struct {
	db database.Service
}

func NewHandlers() *Handlers {
	return &Handlers{db: database.New()}
}

// HealthHandler godoc
//
//	@Summary		检查服务健康状态
//	@Description	返回数据库健康状态
//	@Tags			健康检查
//	@Success		200	{object}	map[string]string	"返回健康状态"
//	@Router			/health [get]
func (h *Handlers) HealthHandler(c *fiber.Ctx) error {
	return c.JSON(h.db.Health())
}

func (h *Handlers) InitDBTables(tables ...interface{}) {
	if len(tables) == 0 {
		return
	}
	e := h.db.GetDB(nil).AutoMigrate(tables...)
	if e != nil {
		zaplog.Logger.Fatal("init db tables failed", zap.Error(e))
	}
}
