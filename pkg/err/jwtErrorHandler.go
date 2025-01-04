package err

import (
	"github.com/gofiber/fiber/v2"
)

func JwtErrorHandler(_ *fiber.Ctx, _ error) error {
	//switch err {
	//case jwt.ErrTokenExpired:
	//	return c.Status(fiber.StatusOK).JSON(Timeout)
	//case jwt.ErrTokenInvalidAudience, jwt.ErrTokenInvalidIssuer, jwt.ErrTokenInvalidSubject, jwt.ErrTokenInvalidId, jwt.ErrTokenInvalidClaims:
	//	return c.Status(fiber.StatusOK).JSON(TokenInvalid)
	//default:
	//	return c.Status(fiber.StatusOK).JSON(fiber.ErrInternalServerError)
	//}
	return TokenNull
}
