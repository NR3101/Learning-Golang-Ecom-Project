package server

import (
	"strconv"

	"github.com/NR3101/go-ecom-project/internal/dto"
	"github.com/NR3101/go-ecom-project/internal/utils"
	"github.com/gin-gonic/gin"
)

// @Summary Get user's cart
// @Description Retrieve the shopping cart of the authenticated user
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=dto.CartResponse} "Cart retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "Cart not found"
// @Router /cart [get]
func (s *Server) getCart(c *gin.Context) {
	userId := c.GetUint("user_id")

	cart, err := s.cartService.GetCart(userId)
	if err != nil {
		utils.NotFoundResponse(c, "Cart not found")
		return
	}

	utils.SuccessResponse(c, "Cart retrieved successfully", cart)
}

// @Summary Add item to cart
// @Description Add a product to the authenticated user's cart
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param addToCartRequest body dto.AddToCartRequest true "Add to cart data"
// @Success 200 {object} utils.Response{data=dto.CartResponse} "Product added to cart successfully"
// @Failure 400 {object} utils.Response "Invalid request body, product not found, or not enough stock"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Router /cart/items [post]
func (s *Server) addToCart(c *gin.Context) {
	userId := c.GetUint("user_id")

	var req dto.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body", err)
		return
	}

	cart, err := s.cartService.AddToCart(userId, &req)
	if err != nil {
		utils.BadRequestResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, "Product added to cart successfully", cart)
}

// @Summary Update cart item quantity
// @Description Update the quantity of a specific item in the cart
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param itemId path int true "Cart Item ID"
// @Param updateCartItemRequest body dto.UpdateCartItemRequest true "Update cart item data"
// @Success 200 {object} utils.Response{data=dto.CartResponse} "Cart item updated successfully"
// @Failure 400 {object} utils.Response "Invalid request body, cart item ID, or not enough stock"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Router /cart/items/{itemId} [put]
func (s *Server) updateCartItem(c *gin.Context) {
	userId := c.GetUint("user_id")

	itemId, err := strconv.ParseUint(c.Param("itemId"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid cart item ID", err)
		return
	}

	var req dto.UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body", err)
		return
	}

	cart, err := s.cartService.UpdateCart(userId, uint(itemId), &req)
	if err != nil {
		utils.BadRequestResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, "Cart item updated successfully", cart)
}

// @Summary Remove item from cart
// @Description Remove a specific item from the cart
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param itemId path int true "Cart Item ID"
// @Success 200 {object} utils.Response "Cart item removed successfully"
// @Failure 400 {object} utils.Response "Invalid cart item ID or cart item not found"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Router /cart/items/{itemId} [delete]
func (s *Server) removeCartItem(c *gin.Context) {
	userId := c.GetUint("user_id")

	itemId, err := strconv.ParseUint(c.Param("itemId"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid cart item ID", err)
		return
	}

	err = s.cartService.RemoveFromCart(userId, uint(itemId))
	if err != nil {
		utils.BadRequestResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, "Cart item removed successfully", nil)
}
