package service

import (
	"context"
	"delivery-tracker/common/kafka"
	"delivery-tracker/common/logger"
	"encoding/json"
	"fmt"
)

type NotificationService struct{}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

func (s *NotificationService) SendNotification(ctx context.Context, msg []byte) error {
	var payload kafka.OrderStatusUpdate

	if err := json.Unmarshal(msg, &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	var notificationMsg string

	switch payload.Status {
	case "created":
		notificationMsg = fmt.Sprintf("Заказ #%d создан пользователем #%d", payload.OrderID)
	case "accepted":
		notificationMsg = fmt.Sprintf("Заказ #%d принят курьером #%d", payload.OrderID, payload.CourierId)
	case "delivered":
		notificationMsg = fmt.Sprintf("Заказ #%d доставлен пользователю #%d", payload.OrderID, payload.UserId)
	case "canceled":
		notificationMsg = fmt.Sprintf("Заказ #%d отменён", payload.OrderID)
	default:
		notificationMsg = fmt.Sprintf("Обновление по заказу #%d: статус %q", payload.OrderID, payload.Status)
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx, notificationMsg)
	return nil
}
