package main

import (
	"context"
	"delivery-tracker/auth/internal/config"
	"delivery-tracker/auth/internal/handler"
	"delivery-tracker/auth/internal/repository"
	"delivery-tracker/auth/internal/service"
	"delivery-tracker/common/logger"
	"delivery-tracker/common/postgres"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	cfg, _ := config.New()

	ctx := context.Background()
	ctx, _ = logger.New(ctx)

	pool, _ := postgres.New(ctx, cfg.Postgres)

	err := postgres.Migrate(ctx, cfg.Postgres, cfg.Auth.MigrationPath)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "migration failed", zap.Error(err))
	}

	router := gin.Default()
	authRepo := repository.NewAuthRepository(pool)
	authService := service.NewAuthService(authRepo, cfg.Auth.SecretKey)
	authHandler := handler.NewAuthHandler(authService, router)

	authHandler.RegisterRoutes()

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Auth.Port), router); err != nil {
			logger.GetLoggerFromCtx(ctx).Fatal(ctx, "auth service failed to start", zap.Error(err))
		}
	}()

	select {}
}
