package middleware

import (
	"log/slog"

	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/config"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/lib"
	"github.com/gofiber/fiber/v3"
)

func RequireAuth(log *slog.Logger, cfg *config.Config) fiber.Handler {
	return func(c fiber.Ctx) error {

		token := ""
		authHeader := c.Get("Authorization")
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		}
		claims, err := lib.GetClaimsFromAccessToken(token, cfg.AUTH_SECRET_KEY, cfg.SERVICE_NAME)
		if err != nil {
			log.Error("GetClaimsFromAccessToken error: ", slog.Any("err", err))
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		//log.Debug("Claims: ", slog.Any("claims", claims))
		c.Locals("user_id", claims.UserId)

		return c.Next()
	}
}
