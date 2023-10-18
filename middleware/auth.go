package middleware

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
)

type HTTPError struct {
	Message string `json:"message"`
}

// Protected protect routes
func Authentication(key string) func(*fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(key),
		ErrorHandler: jwtError,
	})
}

func CtxUserID(c *fiber.Ctx) (id int, err error) {
	cl := claims(c)
	userID := cl["userid"]
	if userID == "" {
		return 0, fiber.ErrNotFound
	}
	fid, ok := userID.(float64)
	if !ok {
		return 0, fiber.ErrNotFound
	}
	return int(fid), nil
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusUnauthorized).JSON(HTTPError{Message: ErrUnauthorized.Error()})
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(HTTPError{Message: err.Error()})
	}
}

func claims(c *fiber.Ctx) jwt.MapClaims {
	userSession := c.Locals("user").(*jwt.Token)
	claims := userSession.Claims.(jwt.MapClaims)
	return claims
}
