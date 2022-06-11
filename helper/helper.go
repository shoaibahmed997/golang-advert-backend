package helper

import (
	"go-ecom/config"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

var Validator = validator.New()

type ValidationError struct {
	Field string
	Tag   string
	Value string
}

func ParseToken(c *fiber.Ctx) error {

	headers := c.GetReqHeaders()
	tokstring := headers["Token"]
	claims := jwt.MapClaims{}

	_, tokenError := jwt.ParseWithClaims(tokstring, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.Config("SECRETKEY")), nil
	})
	if tokenError != nil {
		return c.Status(503).JSON(fiber.Map{"success": false, "error": tokenError.Error()})
	}
	email := claims["email"]
	c.Locals("email", email)

	return c.Next()

}
