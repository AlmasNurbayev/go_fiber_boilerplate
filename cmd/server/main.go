package main

import (
	"flag"
	"fmt"

	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/config"
	"github.com/gofiber/fiber/v3"
)

func main() {
	var configEnv string
	flag.StringVar(&configEnv, "configEnv", "", "Path to env-file")
	flag.Parse()

	// ключевые сообщения дублируем и в консоль и в логгер (он может писать в файл)
	fmt.Println("============ start main ============")
	cfg := config.Mustload(configEnv)
	app := fiber.New()

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Listen(":" + cfg.HTTP_PORT)
}
