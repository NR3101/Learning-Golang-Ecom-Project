package server

import (
	"errors"
	"strconv"

	"github.com/NR3101/go-ecom-project/internal/dto"
	"github.com/NR3101/go-ecom-project/internal/services"
	"github.com/NR3101/go-ecom-project/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

/*--------------------- Category Handlers ---------------------*/

// createCategory handles the HTTP request for creating a new product category.
func (s *Server) createCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body", err)
		return
	}

	productService := services.NewProductService(s.db)
	category, err := productService.CreateCategory(&req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create category", err)
		return
	}

	utils.CreatedResponse(c, "Category created successfully", category)
}

// getCategories handles the HTTP request for retrieving all active product categories.
func (s *Server) getCategories(c *gin.Context) {
	productService := services.NewProductService(s.db)
	categories, err := productService.GetCategories()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve categories", err)
		return
	}

	utils.SuccessResponse(c, "Categories retrieved successfully", categories)
}

// updateCategory handles the HTTP request for updating an existing product category.
func (s *Server) updateCategory(c *gin.Context) {
	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body", err)
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid category ID", err)
		return
	}

	productService := services.NewProductService(s.db)
	category, err := productService.UpdateCategory(uint(id), &req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update category", err)
		return
	}

	utils.SuccessResponse(c, "Category updated successfully", category)
}

// deleteCategory handles the HTTP request for deleting a product category.
func (s *Server) deleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid category ID", err)
		return
	}

	productService := services.NewProductService(s.db)
	if err := productService.DeleteCategory(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.NotFoundResponse(c, "Category not found")
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to delete category", err)
		return
	}

	utils.SuccessResponse(c, "Category deleted successfully", nil)
}

/*--------------------- Product Handlers ---------------------*/

// createProduct handles the HTTP request for creating a new product.
func (s *Server) createProduct(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body", err)
		return
	}

	productService := services.NewProductService(s.db)
	product, err := productService.CreateProduct(&req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create product", err)
		return
	}

	utils.CreatedResponse(c, "Product created successfully", product)
}

// getProducts handles the HTTP request for retrieving all products with pagination and filtering.
func (s *Server) getProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	productService := services.NewProductService(s.db)
	products, meta, err := productService.GetProducts(page, limit)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve products", err)
		return
	}

	utils.PaginatedSuccessResponse(c, "Products retrieved successfully", products, *meta)
}

// getProductByID handles the HTTP request for retrieving a single product by its ID.
func (s *Server) getProductByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid product ID", err)
		return
	}

	productService := services.NewProductService(s.db)
	product, err := productService.GetProduct(uint(id))
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve product", err)
		return
	}

	utils.SuccessResponse(c, "Product retrieved successfully", product)
}

// updateProduct handles the HTTP request for updating an existing product.
func (s *Server) updateProduct(c *gin.Context) {
	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body", err)
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid product ID", err)
		return
	}

	productService := services.NewProductService(s.db)
	product, err := productService.UpdateProduct(uint(id), &req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update product", err)
		return
	}

	utils.SuccessResponse(c, "Product updated successfully", product)
}

// deleteProduct handles the HTTP request for deleting a product.
func (s *Server) deleteProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid product ID", err)
		return
	}

	productService := services.NewProductService(s.db)
	if err := productService.DeleteProduct(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.NotFoundResponse(c, "Product not found")
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to delete product", err)
		return
	}

	utils.SuccessResponse(c, "Product deleted successfully", nil)
}
