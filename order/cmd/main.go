package main

import (
	"context"
	"delivery-tracker/common/kafka"
	"delivery-tracker/common/logger"
	"delivery-tracker/common/postgres"
	"delivery-tracker/order/internal/config"
	"delivery-tracker/order/internal/handler"
	"delivery-tracker/order/internal/repository"
	"delivery-tracker/order/internal/service"
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
	err := postgres.Migrate(ctx, cfg.Postgres, cfg.Order.MigrationPath)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "migration failed", zap.Error(err))
	}

	producer := kafka.NewProducer([]string{"localhost:9092"}, "orders-events")

	router := gin.Default()
	orderRepo := repository.NewOrderRepository(pool)
	orderService := service.NewOrderService(orderRepo, producer)
	orderHandler := handler.NewOrderHandler(orderService, router)

	orderHandler.RegisterRoutes(cfg.Order.SecretKey)

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Order.Port), router); err != nil {
			logger.GetLoggerFromCtx(ctx).Info(ctx, "order service failed to start", zap.Error(err))
		}
	}()

	select {}
}
