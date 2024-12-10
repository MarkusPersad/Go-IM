package err

import (
	"github.com/gofiber/fiber/v2"
)

// HandleError 是一个用于处理错误的函数，它接受一个fiber上下文和一个错误对象作为参数。
// 该函数的主要作用是根据错误的类型，返回不同的HTTP状态码和错误信息。
// 参数:
//
//	ctx *fiber.Ctx: fiber的上下文，用于管理HTTP请求和响应。
//	err error: 发生的错误对象。
//
// 返回值:
//
//	error: 返回错误处理的结果，通常是一个HTTP响应。
func HandleError(ctx *fiber.Ctx, err error) error {
	// 检查错误是否为*PersonalError类型，如果是，则返回该错误的详细信息。
	if e, ok := err.(*PersonalError); ok {
		return ctx.Status(fiber.StatusOK).JSON(PersonalError{
			Code:    e.Code,
			Message: e.Message,
		})
	}
	// 如果错误不是*PersonalError类型，则返回一个默认的内部服务器错误响应。
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.ErrInternalServerError.Code,
		"message": fiber.ErrInternalServerError.Message,
	})
}
