package service

import (
	"context"
	"delivery-tracker/common/kafka"
	"delivery-tracker/order/internal/models"
	"delivery-tracker/order/internal/repository"
	"encoding/json"
	"errors"
	kafka2 "github.com/segmentio/kafka-go"
	"strconv"
)

type OrderService struct {
	repo     *repository.OrderRepository
	producer *kafka.KafkaProducer
}

func NewOrderService(repo *repository.OrderRepository, producer *kafka.KafkaProducer) *OrderService {
	return &OrderService{
		repo:     repo,
		producer: producer,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, req models.OrderCreateRequest) error {
	if err := s.repo.CreateOrder(ctx, req); err != nil {
		return err
	}

	event := kafka.OrderCreate{
		EventType: "order_create",
		OrderId:   req.ID,
		UserId:    req.ID,
		Status:    req.Status,
	}

	value, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := &kafka2.Message{
		Key:   []byte(strconv.Itoa(req.ID)),
		Value: value,
	}

	return s.producer.Send(ctx, msg)
}

func (s *OrderService) UpdateOrder(ctx context.Context, req models.OrderUpdateRequest, orderId int) error {
	if err := s.repo.UpdateOrder(ctx, req, orderId); err != nil {
		return err
	}

	event := kafka.OrderStatusUpdate{
		EventType: "order_status_update",
		OrderID:   orderId,
		Status:    req.Status,
	}

	value, err := json.Marshal(event)
	if err != nil {
		return err
	}
	msg := &kafka2.Message{
		Key:   []byte(strconv.Itoa(orderId)),
		Value: value,
	}

	return s.producer.Send(ctx, msg)
}

func (s *OrderService) GetOrder(ctx context.Context, order *models.OrderGet, orderId int) error {
	if err := s.repo.GetOrderByID(ctx, order, orderId); err != nil {
		return err
	}
	return nil
}

func (s *OrderService) GetUserOrder(ctx context.Context, orders *[]models.OrderGet, userId int) error {
	if err := s.repo.GetOrdersByUserID(ctx, userId, orders); err != nil {
		return err
	}
	return nil
}

func (s *OrderService) DeleteOrder(ctx context.Context, orderId int) error {
	if err := s.repo.DeleteOrder(ctx, orderId); err != nil {
		return err
	}
	return nil
}

func (s *OrderService) AcceptOrder(ctx context.Context, orderId, courierId int) error {
	var order models.OrderGet
	err := s.repo.GetOrderByID(ctx, &order, orderId)
	if err != nil {
		return err
	}
	if order.Status != "created" {
		return errors.New("order already accepted or completed")
	}

	err = s.repo.AssignCourierAndUpdateStatus(ctx, orderId, courierId, "accepted")
	if err != nil {
		return err
	}

	event := kafka.OrderStatusUpdate{
		EventType: "order_accept",
		OrderID:   order.ID,
		Status:    order.Status,
		CourierId: courierId,
	}

	value, err := json.Marshal(event)
	if err != nil {
		return err
	}
	msg := &kafka2.Message{
		Key:   []byte(strconv.Itoa(order.ID)),
		Value: value,
	}

	return s.producer.Send(ctx, msg)
}
