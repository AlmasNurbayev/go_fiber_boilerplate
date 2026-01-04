package middleware

import (
	"log/slog"

	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/config"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/lib"
	"github.com/gofiber/fiber/v3"
)

func RequireAuth(log *slog.Logger, cfg *config.Config) fiber.Handler {
	return func(c fiber.Ctx) error {
		// sess := session.FromContext(c)
		// if sess == nil {
		// 	return c.SendStatus(fiber.StatusUnauthorized)
		// }
		// log.Debug("ID: " + sess.ID())
		// log.Debug("Authenticated: " + fmt.Sprintf("%v", sess.Get("authenticated")))
		// Check if user is authenticated
		// if sess.Get("authenticated") != true {
		// 	log.Error("User is not authenticated")
		// 	return c.SendStatus(fiber.StatusUnauthorized)
		// }

		token := ""
		authHeader := c.Get("Authorization")
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		}
		claims, err := lib.GetClaimsFromAccessToken(token, cfg.AUTH_SECRET_KEY, cfg.SERVICE_NAME)
		if err != nil {
			log.Error("GetClaimsFromAccessToken error: ", err)
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		// if claims.Jti != sess.ID() {
		// 	log.Error("Invalid session ID")
		// 	return c.SendStatus(fiber.StatusUnauthorized)
		// }
		log.Debug("Claims: ", slog.Any("claims", claims))

		return c.Next()
	}
}
