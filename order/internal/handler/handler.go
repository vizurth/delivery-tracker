package handler

import (
	"delivery-tracker/order/internal/models"
	"delivery-tracker/order/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type OrderHandler struct {
	service *service.OrderService
	router  *gin.Engine
}

func NewOrderHandler(service *service.OrderService, router *gin.Engine) *OrderHandler {
	return &OrderHandler{
		service: service,
		router:  router,
	}
}

func (h *OrderHandler) CreateOrderHandler(c *gin.Context) {
	ctx := c.Request.Context()

	var req models.OrderCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.service.CreateOrder(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "create order"})

}

func (h *OrderHandler) UpdateOrderHandler(c *gin.Context) {
	ctx := c.Request.Context()

	var req models.OrderUpdateRequest

	paramId := c.Param("id")
	orderId, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.service.UpdateOrder(ctx, req, orderId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "order update"})
}

func (h *OrderHandler) GetOrderHandler(c *gin.Context) {
	ctx := c.Request.Context()

	var order models.OrderGet

	paramId := c.Param("id")
	orderId, err := strconv.Atoi(paramId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.GetOrder(ctx, &order, orderId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) GetUserOrderHandler(c *gin.Context) {
	ctx := c.Request.Context()

	var orders []models.OrderGet

	paramUserId := c.Param("id")
	userId, err := strconv.Atoi(paramUserId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.service.GetUserOrder(ctx, &orders, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)

}

func (h *OrderHandler) DeleteOrderHandler(c *gin.Context) {
	ctx := c.Request.Context()

	paramOrderId := c.Param("id")
	orderId, err := strconv.Atoi(paramOrderId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.service.DeleteOrder(ctx, orderId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order delete"})
}

func (h *OrderHandler) AcceptOrderHandler(c *gin.Context) {
	ctx := c.Request.Context()
	paramOrderId := c.Param("id")
	courierId := c.GetInt("courier_id")
	orderId, err := strconv.Atoi(paramOrderId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = h.service.AcceptOrder(ctx, orderId, courierId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "order accepted"})
}

func (h *OrderHandler) RegisterRoutes(secret string) {
	order := h.router.Group("/orders")
	order.POST("/", h.CreateOrderHandler)
	order.PUT("/:id/status", h.UpdateOrderHandler)
	order.GET("/:id", h.GetOrderHandler)
	order.GET("/orders/user/:user_id", h.GetUserOrderHandler)
	order.DELETE("/orders/:id", h.DeleteOrderHandler)
	order.PUT("/:id/accept", h.AcceptOrderHandler)
	fmt.Println(secret)
}
