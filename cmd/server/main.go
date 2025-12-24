package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/config"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/httpApp"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/logger"
)

func main() {
	var configEnv string
	flag.StringVar(&configEnv, "configEnv", "", "Path to env-file")
	flag.Parse()

	// ключевые сообщения дублируем и в консоль и в логгер (он может писать в файл)
	fmt.Println("============ start main ============")
	cfg := config.Mustload(configEnv)
	Log := logger.InitLogger(cfg.ENV, cfg.LOG_ERROR_PATH)

	http, err := httpApp.NewHttpApp(Log, cfg)
	if err != nil {
		Log.Error("error create http app", slog.String("err", err.Error()))
		panic(err)
	}

	go func() {
		http.Run()
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	signalString := <-done
	Log.Info("received signal " + signalString.String())
	fmt.Println("received signal " + signalString.String())

	http.Stop()
	Log.Info("http server stopped")
	fmt.Println("============ http server stopped ============")

}
