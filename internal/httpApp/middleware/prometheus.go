package middleware

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/prometheus/client_golang/prometheus"
)

func PrometheusMiddleware(httpRequestCounter *prometheus.CounterVec, httpRequestDuration *prometheus.HistogramVec) fiber.Handler {
	return func(c fiber.Ctx) error {

		start := time.Now()

		// Выполняем следующий обработчик
		err := c.Next()
		if err != nil {
			return err
		}

		// Засекаем время выполнения
		duration := float64(time.Since(start).Milliseconds())

		routePath := "unknown"
		if c.Route() != nil {
			routePath = c.Route().Path
		}
		statusCode := c.Response().StatusCode()
		if statusCode == 0 {
			statusCode = 200 // Значение по умолчанию
		}

		//fmt.Println("httpRequestDuration address", httpRequestDuration)
		// сменили полный originalUrl на частичный routePath
		//httpRequestCounter.WithLabelValues(c.Method(), routePath, strconv.Itoa(statusCode), c.OriginalURL()).Inc()
		//httpRequestDuration.WithLabelValues(c.Method(), routePath, strconv.Itoa(statusCode), c.OriginalURL()).Observe(duration)
		httpRequestCounter.WithLabelValues(c.Method(), routePath, strconv.Itoa(statusCode), routePath).Inc()
		httpRequestDuration.WithLabelValues(c.Method(), routePath, strconv.Itoa(statusCode), routePath).Observe(duration)

		return nil
	}

}
