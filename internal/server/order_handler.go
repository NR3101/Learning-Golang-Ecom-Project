package server

import (
	"strconv"

	"github.com/NR3101/go-ecom-project/internal/utils"
	"github.com/gin-gonic/gin"
)

// @Summary Create an order
// @Description Create a new order from the user's cart
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response "Order created successfully"
// @Failure 400 {object} utils.Response "Failed to create order (empty cart or insufficient stock)"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Router /orders [post]
func (s *Server) createOrder(c *gin.Context) {
	userID := c.GetUint("user_id")

	order, err := s.orderService.CreateOrder(userID)
	if err != nil {
		utils.BadRequestResponse(c, "Failed to create order", err)
		return
	}

	utils.SuccessResponse(c, "Order created successfully", order)
}

// @Summary Get all orders
// @Description Retrieve all orders of the authenticated user with pagination
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} utils.PaginatedResponse "Orders retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /orders [get]
func (s *Server) getOrders(c *gin.Context) {
	userID := c.GetUint("user_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	orders, meta, err := s.orderService.GetAllOrders(userID, page, limit)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve orders", err)
		return
	}

	utils.PaginatedSuccessResponse(c, "Orders retrieved successfully", orders, *meta)
}

// @Summary Get order by ID
// @Description Retrieve a specific order by its ID
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param orderId path int true "Order ID"
// @Success 200 {object} utils.Response "Order retrieved successfully"
// @Failure 400 {object} utils.Response "Invalid order ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "Order not found"
// @Router /orders/{orderId} [get]
func (s *Server) getOrder(c *gin.Context) {
	userID := c.GetUint("user_id")

	orderID, err := strconv.ParseUint(c.Param("orderId"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid order ID", err)
		return
	}

	order, err := s.orderService.GetOrderByID(userID, uint(orderID))
	if err != nil {
		utils.NotFoundResponse(c, "Order not found")
		return
	}

	utils.SuccessResponse(c, "Order retrieved successfully", order)
}
