package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/config"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/httpApp"
	"github.com/AlmasNurbayev/go_fiber_boilerplate/internal/lib"
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

	prometheus := lib.NewPromRegistry(Log)
	mux := http.NewServeMux()
	lib.RegisterMetricsHandlerWithRegistry(mux, prometheus.Registry)

	go func() {
		addr := ":" + cfg.PROMETHEUS_HTTP_PORT
		Log.Info("metrics server started", slog.String("addr", addr))
		if err := http.ListenAndServe(addr, mux); err != nil && err != http.ErrServerClosed {
			Log.Error("metrics server error", slog.String("err", err.Error()))
		}
	}()

	httpFiber, err := httpApp.NewHttpApp(Log, cfg, prometheus)
	if err != nil {
		Log.Error("error create http app", slog.String("err", err.Error()))
		panic(err)
	}

	go func() {
		httpFiber.Run()
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	signalString := <-done
	Log.Info("received signal " + signalString.String())
	fmt.Println("received signal " + signalString.String())

	httpFiber.Stop()
	Log.Info("http server stopped")
	fmt.Println("============ http server stopped ============")

}
