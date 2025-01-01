package err

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JwtErrorHandler(c *fiber.Ctx, err error) error {
	switch err {
	case jwt.ErrTokenExpired:
		return c.Status(fiber.StatusOK).JSON(Timeout)
	case jwt.ErrTokenInvalidAudience, jwt.ErrTokenInvalidIssuer, jwt.ErrTokenInvalidSubject, jwt.ErrTokenInvalidId, jwt.ErrTokenInvalidClaims:
		return c.Status(fiber.StatusOK).JSON(TokenInvalid)
	default:
		return c.Status(fiber.StatusOK).JSON(fiber.ErrInternalServerError)
	}
}
