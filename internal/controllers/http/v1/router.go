package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"github.com/timohahaa/executor/internal/service"
)

func NewRouter(commandService service.CommandService, logger *logrus.Logger) *echo.Echo {
	e := echo.New()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogMethod:   true,
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		LogRemoteIP: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.WithFields(logrus.Fields{
				"method": v.Method,
				"URI":    v.URI,
				"status": v.Status,
				"ip":     v.RemoteIP,
				"error":  v.Error,
			}).Info("request")

			return nil
		},
	}))

	e.GET("/health", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	v1 := e.Group("/api/v1")
	{
		newCommandRoutes(v1, commandService)
	}
	return e
}
