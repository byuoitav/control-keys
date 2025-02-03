package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/byuoitav/control-keys/codemap"
	"github.com/byuoitav/control-keys/handlers"
	"github.com/byuoitav/control-keys/middleware"
	"github.com/byuoitav/control-keys/opa"
	"github.com/byuoitav/control-keys/wso2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	var (
		logLevel int8
		opaURL   string
		opaToken string
	)

	pflag.Int8VarP(&logLevel, "log-level", "L", 0, "Level to log at. Provided by zap logger: https://godoc.org/go.uber.org/zap/zapcore")
	pflag.StringVarP(&opaURL, "opa-address", "a", "", "OPA Address (Full URL)")
	pflag.StringVarP(&opaToken, "opa-token", "t", "", "OPA Token")

	port := ":8029"
	router := gin.Default()

	// Build out the Logger
	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(zapcore.Level(logLevel)),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "@",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	plain, err := config.Build()
	if err != nil {
		fmt.Printf("unable to build logger: %s", err)
		os.Exit(1)
	}

	logger := plain.Sugar()

	// WSO2 Create Client
	client := wso2.New("", "", "https://api.byu.edu", "")

	o := opa.Client{
		URL:   opaURL,
		Token: opaToken,
	}

	// build the main group and pass the middleware of WSO2
	api := router.Group("/api/v1")
	api.Use(func(c *gin.Context) {
		client.JWTValidationMiddleware()
		c.Next()
	})
	api.Use(func(c *gin.Context) {
		if middleware.Authenticated(c.Request) {
			c.Next()
			return
		}
		logger.Info("WSO2 Authentication Failed")
		logger.Debug("Output of JWT: %s", c.Request)
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		c.Abort()
	})
	api.Use(func(c *gin.Context) {
		o.Authorize()
		c.Next()
	})

	c := codemap.New()
	c.Start()
	h := handlers.New(c)

	// Functionality Endpoints
	router.GET("/:controlKey/getPreset", func(ctx *gin.Context) {
		h.GetPresetHandler(ctx)
	})
	router.GET("/:preset/getControlKey", func(ctx *gin.Context) {
		h.GetControlKeyHandler(ctx)
	})
	router.GET("/:room/refresh", func(ctx *gin.Context) {
		h.RefreshPresetKey(ctx)
	})
	router.GET("/status", func(ctx *gin.Context) {
		h.HealthCheck(ctx)
	})

	server := &http.Server{
		Addr:           port,
		Handler:        router,
		MaxHeaderBytes: 1024 * 10,
	}

	logger.Info("Starting Service.....")
	_ = server.ListenAndServe()
}
