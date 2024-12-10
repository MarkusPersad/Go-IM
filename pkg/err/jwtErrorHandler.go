package err

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JwtErrorHandler(c *fiber.Ctx, err error) error {
	switch err {
	case jwt.ErrTokenExpired:
		return Timeout
	case jwt.ErrTokenInvalidAudience, jwt.ErrTokenInvalidIssuer, jwt.ErrTokenInvalidSubject, jwt.ErrTokenInvalidId, jwt.ErrTokenInvalidClaims:
		return TokenInvalid
	default:
		return fiber.ErrInternalServerError
	}
}
