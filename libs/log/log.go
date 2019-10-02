package log

import (
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// GetLogger returns a configured logger
func GetLogger() *logrus.Logger {
	logger := logrus.New()
	env := os.Getenv("ENV")

	if env == "production" || env == "prod" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	}

	return logger
}

// EchoLoggingMiddleware returns a logging middleware that handles selecting
// appropriate format for the current context, JSON in production otherwise
// colored and column formatting etc. Adds a bunch of Echo-specific logging
// fields automatically.
func EchoLoggingMiddleware() echo.MiddlewareFunc {
	return loggingMiddleware()
}

func loggingMiddleware() echo.MiddlewareFunc {
	logger := GetLogger()
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			req := c.Request()
			res := c.Response()
			start := time.Now()
			if err = next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()

			requestID := req.Header.Get(echo.HeaderXRequestID)
			if requestID == "" {
				requestID = res.Header().Get(echo.HeaderXRequestID)
			}
			latency := stop.Sub(start)

			requestPath := req.URL.Path
			if requestPath == "" {
				requestPath = "/"
			}

			requestError := ""
			if err != nil {
				b, _ := json.Marshal(err.Error())
				requestError = string(b[1 : len(b)-1])
			}

			// TODO: Add userID from auth token
			logger.WithFields(logrus.Fields{
				"id":            requestID,
				"remote_ip":     c.RealIP(),
				"user_agent":    req.UserAgent(),
				"status":        res.Status,
				"latency":       strconv.FormatInt(int64(latency), 10),
				"latency_human": stop.Sub(start).String(),
				"method":        req.Method,
				"path":          requestPath,
				"error":         requestError,
			}).Info(req.Method + " " + req.RequestURI)

			return nil
		}
	}
}
