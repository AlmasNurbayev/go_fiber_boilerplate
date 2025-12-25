package httpApp

import (
	"context"
	"log/slog"

	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/config"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/db/storage"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/httpApp/middleware"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/lib"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

type structValidator struct {
	validate *validator.Validate
}

type HttpApp struct {
	Log     *slog.Logger
	Server  *fiber.App
	Storage *storage.Storage
	Cfg     *config.Config
}

func (v *structValidator) Validate(out any) error {
	return v.validate.Struct(out)
}

func NewHttpApp(
	log *slog.Logger,
	cfg *config.Config,
	prometheus lib.PrometheusType,
) (*HttpApp, error) {

	dsn := "postgres://" + cfg.POSTGRES_USER + ":" + cfg.POSTGRES_PASSWORD + "@" + cfg.POSTGRES_HOST + ":" + cfg.POSTGRES_INT_PORT + "/" + cfg.POSTGRES_DB + "?sslmode=disable"

	ctxDB, cancel := context.WithTimeout(context.Background(), cfg.POSTGRES_TIMEOUT)
	defer cancel()

	storage, err := storage.NewStorage(ctxDB, dsn, log)
	if err != nil {
		log.Error("not init main storage")
		return nil, err
	}

	server := fiber.New(fiber.Config{
		StructValidator: &structValidator{validate: validator.New()},
		ReadTimeout:     cfg.HTTP_TIMEOUT,
		WriteTimeout:    cfg.HTTP_TIMEOUT,
		IdleTimeout:     cfg.HTTP_TIMEOUT,
	})

	if cfg.ENV != "prod" {
		//server.Use(middleware.RequestTracingMiddleware(log))
	}

	server.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.HTTP_CORS_ALLOW_ORIGINS,
		AllowCredentials: cfg.HTTP_CORS_ALLOW_CREDENTIALS,
		AllowHeaders:     cfg.HTTP_CORS_ALLOW_HEADERS,
	}))

	server.Use(middleware.PrometheusMiddleware(prometheus.CounterVec, prometheus.HistogramVec))

	RegisterMainRoutes(server, storage, log, cfg)

	server.Get("/healthz", func(c fiber.Ctx) error {
		return c.Status(200).SendString("OK")
	})

	return &HttpApp{
		Log:     log,
		Server:  server,
		Storage: storage,
		Cfg:     cfg,
	}, nil
}

func (a *HttpApp) Run() {
	err := a.Server.Listen(":"+a.Cfg.HTTP_PORT, fiber.ListenConfig{
		EnablePrefork:   a.Cfg.HTTP_PREFORK,
		ShutdownTimeout: a.Cfg.HTTP_TIMEOUT,
	})
	if err != nil {
		a.Log.Error("not start server: ", slog.String("err", err.Error()))
		panic(err)
	}
}

func (a *HttpApp) Stop() {
	err := a.Server.Shutdown()
	a.Storage.Close()
	if err != nil {
		a.Log.Error("error on stop server: ", slog.String("err", err.Error()))
		panic(err)
	}
}
