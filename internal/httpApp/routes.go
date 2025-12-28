package httpApp

import (
	"log/slog"

	_ "github.com/AlmasNurbayev/go_fiber_boilerplate/docs"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/config"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/db/storage"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/httpApp/handlers"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/httpApp/services"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/swagger/v2"
)

func RegisterMainRoutes(app *fiber.App, storage *storage.Storage, log *slog.Logger, cfg *config.Config) {
	cp := "registerRoutes"
	log = log.With(slog.String("cp", cp))
	log.Info("Register routes:")

	app.Get("/swagger/*", swagger.HandlerDefault) // default

	log.Info("/api")
	api := app.Group("/api")
	RegisterUserRoutes(api, storage, log, cfg)
	RegisterAuthRoutes(api, storage, log, cfg)
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

func RegisterAuthRoutes(api fiber.Router, storage *storage.Storage, log *slog.Logger, cfg *config.Config) {

	authService := services.NewAuthService(log, storage, cfg)
	authHandler := handlers.NewAuthHandler(log, authService)

	log.Info("POST /api/auth/register")
	api.Post("/auth/register", authHandler.AuthRegister)
	log.Info("POST /api/auth/login")
	api.Post("/auth/login", authHandler.AuthLogin)
	log.Info("GET /api/auth/hello")
	api.Get("/auth/hello", authHandler.AuthHello)

}
