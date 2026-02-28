package services

import (
	"errors"

	"github.com/NR3101/go-ecom-project/internal/dto"
	"github.com/NR3101/go-ecom-project/internal/models"
	"gorm.io/gorm"
)

var _ CartServiceInterface = (*CartService)(nil)

type CartService struct {
	db *gorm.DB
}

func NewCartService(db *gorm.DB) *CartService {
	return &CartService{db: db}
}

// GetCart retrieves the cart for a given user ID, including cart items and their associated products and categories.
func (s *CartService) GetCart(userID uint) (*dto.CartResponse, error) {
	var cart models.Cart
	if err := s.db.Preload("CartItems.Product.Category").Where("user_id = ?", userID).First(&cart).Error; err != nil {
		return nil, err
	}

	return s.ConvertToCartResponse(&cart), nil
}

// AddToCart adds a product to the user's cart, checking for product existence and stock availability. If the product is already in the cart, it updates the quantity.
func (s *CartService) AddToCart(userID uint, req *dto.AddToCartRequest) (*dto.CartResponse, error) {
	// Check if the product exists
	var product models.Product
	if err := s.db.First(&product, req.ProductID).Error; err != nil {
		return nil, errors.New("product not found")
	}

	// Check if the product has enough stock
	if product.Stock < req.Quantity {
		return nil, errors.New("not enough stock available")
	}

	// Get or create the user's cart
	var cart models.Cart
	if err := s.db.Where("user_id = ?", userID).FirstOrCreate(&cart, models.Cart{UserID: userID}).Error; err != nil {
		return nil, err
	}

	// Check if the product is already in the cart (including soft-deleted items)
	var cartItem models.CartItem
	err := s.db.Unscoped().Where("cart_id = ? AND product_id = ?", cart.ID, req.ProductID).First(&cartItem).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Product is not in the cart, create a new cart item
			cartItem = models.CartItem{
				CartID:    cart.ID,
				ProductID: req.ProductID,
				Quantity:  req.Quantity,
			}
			if err := s.db.Create(&cartItem).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else if cartItem.DeletedAt.Valid {
		// Item was soft-deleted, restore it with new quantity
		cartItem.DeletedAt = gorm.DeletedAt{}
		cartItem.Quantity = req.Quantity
		if err := s.db.Unscoped().Save(&cartItem).Error; err != nil {
			return nil, err
		}
	} else {
		// Product is already in the cart, update the quantity
		cartItem.Quantity += req.Quantity
		if err := s.db.Save(&cartItem).Error; err != nil {
			return nil, err
		}
	}

	return s.GetCart(userID)
}

// UpdateCart updates the quantity of a specific cart item for a user, checking for product existence and stock availability.
func (s *CartService) UpdateCart(userID, itemID uint, req *dto.UpdateCartItemRequest) (*dto.CartResponse, error) {
	// First, get the user's cart
	var cart models.Cart
	if err := s.db.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		return nil, errors.New("cart not found")
	}

	// Find the cart item that belongs to this cart
	var cartItem models.CartItem
	if err := s.db.Where("id = ? AND cart_id = ?", itemID, cart.ID).First(&cartItem).Error; err != nil {
		return nil, errors.New("cart item not found")
	}

	// Check if the product has enough stock
	var product models.Product
	if err := s.db.First(&product, cartItem.ProductID).Error; err != nil {
		return nil, errors.New("product not found")
	}
	if product.Stock < req.Quantity {
		return nil, errors.New("not enough stock available")
	}

	cartItem.Quantity = req.Quantity
	if err := s.db.Save(&cartItem).Error; err != nil {
		return nil, err
	}

	return s.GetCart(userID)
}

// RemoveFromCart removes a specific cart item from the user's cart.
func (s *CartService) RemoveFromCart(userID, itemID uint) error {
	// First, get the user's cart
	var cart models.Cart
	if err := s.db.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		return errors.New("cart not found")
	}

	// Delete the cart item that belongs to this cart
	result := s.db.Where("id = ? AND cart_id = ?", itemID, cart.ID).Delete(&models.CartItem{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("cart item not found")
	}
	return nil
}

// ConvertToCartResponse converts a Cart model to a CartResponse DTO, calculating subtotals and total price.
func (s *CartService) ConvertToCartResponse(cart *models.Cart) *dto.CartResponse {
	cartItems := make([]dto.CartItemResponse, len(cart.CartItems))
	var total float64

	for i := range cart.CartItems {
		subtotal := float64(cart.CartItems[i].Quantity) * cart.CartItems[i].Product.Price
		total += subtotal

		cartItems[i] = dto.CartItemResponse{
			ID: cart.CartItems[i].ID,
			Product: dto.ProductResponse{
				ID:          cart.CartItems[i].Product.ID,
				CategoryID:  cart.CartItems[i].Product.CategoryID,
				Name:        cart.CartItems[i].Product.Name,
				Description: cart.CartItems[i].Product.Description,
				Price:       cart.CartItems[i].Product.Price,
				Stock:       cart.CartItems[i].Product.Stock,
				SKU:         cart.CartItems[i].Product.SKU,
				IsActive:    cart.CartItems[i].Product.IsActive,
				Category: dto.CategoryResponse{
					ID:          cart.CartItems[i].Product.Category.ID,
					Name:        cart.CartItems[i].Product.Category.Name,
					Description: cart.CartItems[i].Product.Category.Description,
					IsActive:    cart.CartItems[i].Product.Category.IsActive,
				},
			},
			Quantity:  cart.CartItems[i].Quantity,
			Subtotal:  subtotal,
			CreatedAt: cart.CartItems[i].CreatedAt,
			UpdatedAt: cart.CartItems[i].UpdatedAt,
		}
	}

	return &dto.CartResponse{
		ID:        cart.ID,
		UserID:    cart.UserID,
		CartItems: cartItems,
		Total:     total,
		CreatedAt: cart.CreatedAt,
		UpdatedAt: cart.UpdatedAt,
	}
}
