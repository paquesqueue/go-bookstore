package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/paquesqueue/bookstore/common"
	"github.com/sirupsen/logrus"
)

func InitMiddleware(e *echo.Echo, log *logrus.Logger, config common.Config) {

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authToken := c.Request().Header.Values("Authorization")
			if authToken != nil && authToken[0] == config.AccessToken {
				return next(c)
			}
			return echo.NewHTTPError(http.StatusUnauthorized, "Valid credential not provided")
		}
	})

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:      true,
		LogStatus:   true,
		LogRemoteIP: true,
		LogMethod:   true,
		LogHeaders:  []string{"Content-Type", "Authorization"},
		LogLatency:  true,
		LogError:    true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			if values.Error == nil {
				log.WithFields(logrus.Fields{
					"URI":     values.URI,
					"status":  values.Status,
					"method":  values.Method,
					"headers": values.Headers,
					"latency": values.Latency,
				}).Info("request")
			} else {
				log.WithFields(logrus.Fields{
					"URI":     values.URI,
					"status":  values.Status,
					"method":  values.Method,
					"headers": values.Headers,
					"latency": values.Latency,
					"error":   values.Error,
				}).Error("request error")
			}
			return nil
		},
	}))
	e.Use(middleware.Recover())
}
