package server

import (
	"strconv"

	"github.com/NR3101/go-ecom-project/internal/dto"
	"github.com/NR3101/go-ecom-project/internal/utils"
	"github.com/gin-gonic/gin"
)

func (s *Server) getCart(c *gin.Context) {
	userId := c.GetUint("user_id")

	cart, err := s.cartService.GetCart(userId)
	if err != nil {
		utils.NotFoundResponse(c, "Cart not found")
		return
	}

	utils.SuccessResponse(c, "Cart retrieved successfully", cart)
}

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
