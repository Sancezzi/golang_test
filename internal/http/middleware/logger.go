package middleware

import (
	"test-api/internal/http/responses"
	logs "test-api/internal/log"

	"github.com/gofiber/fiber/v3"
)

func LoggerMiddleware(logger *logs.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		defer func() {
			if err := recover(); err != nil {
				logger.Error(err)
			}
		}()

		err := c.Next()
		if err != nil {
			logger.Info(err)
		}

		return c.JSON(responses.ErrorResponse(err))
	}
}
