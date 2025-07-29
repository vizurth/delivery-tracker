package repository

import (
	"context"
	"database/sql"
	"delivery-tracker/order/internal/models"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, req models.OrderCreateRequest) error {
	_, err := r.db.Exec(ctx,
		"INSERT INTO schema_name.orders (user_id, status) VALUES ($1, $2)", req.UserId, req.Status)
	if err != nil {
		return errors.New("error inserting order")
	}

	return nil
}

func (r *OrderRepository) UpdateOrder(ctx context.Context, req models.OrderUpdateRequest, orderId int) error {
	_, err := r.db.Exec(ctx,
		"UPDATE INTO schema_name.orders SET status = $1 WHERE id = $2", req.Status, orderId,
	)
	if err != nil {
		return errors.New("error updating order")
	}
	return nil
}

func (r *OrderRepository) GetOrderByID(ctx context.Context, order *models.OrderGet, orderId int) error {
	// Используем QueryRow, потому что ожидаем только одну строку
	row := r.db.QueryRow(
		ctx,
		"SELECT * FROM schema_name.orders WHERE id = $1",
		orderId,
	)

	// Заполняем структуру order
	err := row.Scan(
		&order.ID,
		&order.UserId,
		&order.Status,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("order with ID %d not found", orderId)
		}
		return fmt.Errorf("failed to get order: %w", err)
	}

	return nil
}

func (r *OrderRepository) GetOrdersByUserID(ctx context.Context, userId int, orders *[]models.OrderGet) error {
	rows, _ := r.db.Query(
		ctx,
		"SELECT * FROM schema_name.orders WHERE user_id = $1",
		userId,
	)

	for rows.Next() {
		temp := models.OrderGet{}
		err := rows.Scan(
			&temp.ID,
			&temp.UserId,
			&temp.Status,
		)
		if err != nil {
			return fmt.Errorf("failed to get order: %w", err)
		}
		*orders = append(*orders, temp)
	}
	return nil
}

func (r *OrderRepository) DeleteOrder(ctx context.Context, orderId int) error {
	_, err := r.db.Exec(
		ctx,
		"DELETE FROM schema_name.orders WHERE id = $1",
	)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}
	return nil
}

func (r *OrderRepository) AssignCourierAndUpdateStatus(ctx context.Context, orderId, courierId int, status string) error {
	row := r.db.QueryRow(ctx, "SELECT * FROM schema_name.orders WHERE id = $1", orderId)

	temp := models.OrderAcceptRequest{}
	err := row.Scan(
		&temp.ID,
		&temp.UserId,
		&temp.CourierId,
		&temp.Status,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("order with ID %d not found", orderId)
		}
	}
	if temp.CourierId != "-1" {
		return errors.New("diff courier must accept")
	}
	_, err = r.db.Exec(ctx,
		"UPDATE schema_name.orders SET status = $1, courier_id = $2 WHERE id = $3",
		status,
		courierId,
		orderId,
	)
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}
	return nil
}
