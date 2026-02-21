package services

import (
	"fmt"

	"github.com/NR3101/go-ecom-project/internal/dto"
	"github.com/NR3101/go-ecom-project/internal/models"
	"github.com/NR3101/go-ecom-project/internal/utils"
	"gorm.io/gorm"
)

const (
	dateFormat = "2006-01-02T15:04:05Z"
)

type OrderService struct {
	db *gorm.DB
}

func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{db: db}
}

// CreateOrder creates a new order for the given user ID by processing the user's cart.
func (s *OrderService) CreateOrder(userId uint) (*dto.OrderResponse, error) {
	var orderResponse *dto.OrderResponse

	// Start a transaction to ensure atomicity of the order creation process
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Fetch the user's cart with associated products
		var cart models.Cart
		if err := tx.Preload("CartItems.Product").Where("user_id = ?", userId).First(&cart).Error; err != nil {
			return fmt.Errorf("could not find cart by user_id %d", userId)
		}

		// Check if the cart is empty
		if len(cart.CartItems) == 0 {
			return fmt.Errorf("cart is empty for user_id %d", userId)
		}

		var totalAmount float64
		var orderItems []models.OrderItem

		// Iterate through cart items to calculate total amount and prepare order items
		for i := range cart.CartItems {
			cartItem := &cart.CartItems[i]

			// Check if there is enough stock for the product
			if cartItem.Product.Stock < cartItem.Quantity {
				return fmt.Errorf("not enough stock for product %s", cartItem.Product.Name)
			}

			// Calculate total amount for the order
			itemTotal := float64(cartItem.Quantity) * cartItem.Product.Price
			totalAmount += itemTotal

			// Prepare order item for the order
			orderItems = append(orderItems, models.OrderItem{
				ProductID: cartItem.ProductID,
				Quantity:  cartItem.Quantity,
				Price:     cartItem.Product.Price,
			})

			// Update product stock
			cartItem.Product.Stock -= cartItem.Quantity
			if err := tx.Save(cartItem.Product).Error; err != nil {
				return fmt.Errorf("failed to update stock for product %s", cartItem.Product.Name)
			}

			// Create order
			order := models.Order{
				UserID:      userId,
				Status:      models.OrderStatusPending,
				TotalAmount: totalAmount,
				OrderItems:  orderItems,
			}

			if err := tx.Create(&order).Error; err != nil {
				return fmt.Errorf("failed to create order: %v", err)
			}

			// Clear the cart after creating the order
			if err := tx.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
				return fmt.Errorf("failed to clear cart items: %v", err)
			}

			// Prepare the order response to return after the transaction is committed
			response, err := s.getOrderResponse(tx, order.ID)
			if err != nil {
				return fmt.Errorf("failed to get order response: %v", err)
			}

			orderResponse = response
		}
		return nil // Return nil to commit the transaction
	})

	if err != nil {
		return nil, err
	}

	return orderResponse, nil
}

// GetAllOrders retrieves a paginated list of orders for the given user ID, along with pagination metadata.
func (s *OrderService) GetAllOrders(userId uint, page, limit int) ([]dto.OrderResponse, *utils.PaginationMeta, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit
	var orders []models.Order
	var total int64

	s.db.Model(&models.Order{}).Where("user_id = ?", userId).Count(&total)

	if err := s.db.Preload("OrderItems.Product.Category").
		Where("user_id = ?", userId).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&orders).Error; err != nil {
		return nil, nil, fmt.Errorf("could not retrieve orders for user_id %d: %v", userId, err)
	}

	response := make([]dto.OrderResponse, len(orders))
	for i := range orders {
		order := &orders[i]
		response[i] = s.convertToOrderResponse(order)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	paginationMeta := &utils.PaginationMeta{
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}

	return response, paginationMeta, nil
}

// GetOrderByID retrieves a specific order by its ID for the given user ID and returns it as an OrderResponse DTO.
func (s *OrderService) GetOrderByID(userId, orderId uint) (*dto.OrderResponse, error) {
	var order models.Order
	if err := s.db.Preload("OrderItems.Product.Category").
		Where("id = ? AND user_id = ?", orderId, userId).
		First(&order).Error; err != nil {
		return nil, fmt.Errorf("could not find order by id %d for user_id %d", orderId, userId)
	}

	orderResponse := s.convertToOrderResponse(&order)
	return &orderResponse, nil
}

// getOrderResponse retrieves an order by ID and converts it to an OrderResponse DTO
func (s *OrderService) getOrderResponse(tx *gorm.DB, orderId uint) (*dto.OrderResponse, error) {
	var order models.Order
	if err := tx.Preload("OrderItems.Product.Category").First(&order, orderId).Error; err != nil {
		return nil, fmt.Errorf("could not find order by id %d", orderId)
	}

	orderResponse := s.convertToOrderResponse(&order)
	return &orderResponse, nil
}

// convertToOrderResponse converts an Order model to an OrderResponse DTO
func (s *OrderService) convertToOrderResponse(order *models.Order) dto.OrderResponse {
	orderItems := make([]dto.OrderItemResponse, len(order.OrderItems))

	for i := range order.OrderItems {
		item := order.OrderItems[i]

		orderItems[i] = dto.OrderItemResponse{
			ID: item.ID,
			Product: dto.ProductResponse{
				ID:          item.Product.ID,
				CategoryID:  item.Product.CategoryID,
				Name:        item.Product.Name,
				Description: item.Product.Description,
				Price:       item.Product.Price,
				Stock:       item.Product.Stock,
				SKU:         item.Product.SKU,
				IsActive:    item.Product.IsActive,
				Category: dto.CategoryResponse{
					ID:          item.Product.Category.ID,
					Name:        item.Product.Category.Name,
					Description: item.Product.Category.Description,
					IsActive:    item.Product.IsActive,
				},
			},
			Quantity: item.Quantity,
			Price:    item.Price,
		}
	}

	return dto.OrderResponse{
		ID:          order.ID,
		UserID:      order.UserID,
		Status:      string(order.Status),
		TotalAmount: order.TotalAmount,
		OrderItems:  orderItems,
		CreatedAt:   order.CreatedAt.Format(dateFormat),
	}
}
