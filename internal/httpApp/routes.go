package httpApp

import (
	"log/slog"

	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/config"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/db/storage"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/httpApp/handlers"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/httpApp/services"
	"github.com/gofiber/fiber/v3"
)

func RegisterMainRoutes(app *fiber.App, storage *storage.Storage, log *slog.Logger, cfg *config.Config) {
	cp := "registerRoutes"
	log = log.With(slog.String("cp", cp))
	log.Info("Register routes:")

	log.Info("/api")
	api := app.Group("/api")
	RegisterUserRoutes(api, storage, log, cfg)
}

func RegisterUserRoutes(api fiber.Router, storage *storage.Storage, log *slog.Logger, cfg *config.Config) {

	userService := services.NewUserService(log, storage, cfg)
	userHandler := handlers.NewUserHandler(log, userService)

	log.Info("GET /api/user")
	api.Get("/user/search/", userHandler.GetUserSearch)

	log.Info("GET /api/user/:id?")
	api.Get("/user/:id?", userHandler.GetUserById)

	log.Info("GET /api/user/search/")
	api.Get("/user/search/", userHandler.GetUserSearch)

}
