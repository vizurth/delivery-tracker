package main

import (
	"context"
	"delivery-tracker/common/kafka"
	"delivery-tracker/common/logger"
	"delivery-tracker/notifications/internal/service"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//cfg, _ := config.New()

	topic := "orders-events"
	groupID := "notification-group"

	consumer := kafka.NewConsumer([]string{"localhost:9092"}, topic, groupID)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx, _ = logger.New(ctx)

	serv := service.NewNotificationService()

	go consumer.Consume(ctx, serv.SendNotification)

	logger.GetLoggerFromCtx(ctx).Info(ctx, "starting notification consumer")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	logger.GetLoggerFromCtx(ctx).Info(ctx, "received shutdown signal")

	consumer.Close()
}
