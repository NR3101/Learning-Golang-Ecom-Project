package server

import (
	"errors"
	"strconv"

	"github.com/NR3101/go-ecom-project/internal/dto"
	"github.com/NR3101/go-ecom-project/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

/*--------------------- Category Handlers ---------------------*/

// @Summary Create a new category
// @Description Create a new product category (Admin only)
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param createCategoryRequest body dto.CreateCategoryRequest true "Category creation data"
// @Success 201 {object} utils.Response{data=dto.CategoryResponse} "Category created successfully"
// @Failure 400 {object} utils.Response "Invalid request body"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden - Admin access required"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /categories [post]
func (s *Server) createCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body", err)
		return
	}

	category, err := s.productService.CreateCategory(&req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create category", err)
		return
	}

	utils.CreatedResponse(c, "Category created successfully", category)
}

// @Summary Get all categories
// @Description Retrieve all active product categories
// @Tags Categories
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=[]dto.CategoryResponse} "Categories retrieved successfully"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /categories [get]
func (s *Server) getCategories(c *gin.Context) {
	categories, err := s.productService.GetCategories()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve categories", err)
		return
	}

	utils.SuccessResponse(c, "Categories retrieved successfully", categories)
}

// @Summary Update a category
// @Description Update an existing product category (Admin only)
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Param updateCategoryRequest body dto.UpdateCategoryRequest true "Category update data"
// @Success 200 {object} utils.Response{data=dto.CategoryResponse} "Category updated successfully"
// @Failure 400 {object} utils.Response "Invalid request body or category ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden - Admin access required"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /categories/{id} [put]
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

	category, err := s.productService.UpdateCategory(uint(id), &req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update category", err)
		return
	}

	utils.SuccessResponse(c, "Category updated successfully", category)
}

// @Summary Delete a category
// @Description Delete a product category (Admin only)
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Success 200 {object} utils.Response "Category deleted successfully"
// @Failure 400 {object} utils.Response "Invalid category ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden - Admin access required"
// @Failure 404 {object} utils.Response "Category not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /categories/{id} [delete]
func (s *Server) deleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid category ID", err)
		return
	}

	if err := s.productService.DeleteCategory(uint(id)); err != nil {
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

// @Summary Create a new product
// @Description Create a new product (Admin only)
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param createProductRequest body dto.CreateProductRequest true "Product creation data"
// @Success 201 {object} utils.Response{data=dto.ProductResponse} "Product created successfully"
// @Failure 400 {object} utils.Response "Invalid request body"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden - Admin access required"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /products [post]
func (s *Server) createProduct(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body", err)
		return
	}

	product, err := s.productService.CreateProduct(&req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create product", err)
		return
	}

	utils.CreatedResponse(c, "Product created successfully", product)
}

// @Summary Get all products
// @Description Retrieve all products with pagination
// @Tags Products
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} utils.PaginatedResponse{data=[]dto.ProductResponse} "Products retrieved successfully"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /products [get]
func (s *Server) getProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	products, meta, err := s.productService.GetProducts(page, limit)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve products", err)
		return
	}

	utils.PaginatedSuccessResponse(c, "Products retrieved successfully", products, *meta)
}

// @Summary Get a product by ID
// @Description Retrieve a single product by its ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} utils.Response{data=dto.ProductResponse} "Product retrieved successfully"
// @Failure 400 {object} utils.Response "Invalid product ID"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /products/{id} [get]
func (s *Server) getProductByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid product ID", err)
		return
	}

	product, err := s.productService.GetProduct(uint(id))
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve product", err)
		return
	}

	utils.SuccessResponse(c, "Product retrieved successfully", product)
}

// @Summary Update a product
// @Description Update an existing product (Admin only)
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param updateProductRequest body dto.UpdateProductRequest true "Product update data"
// @Success 200 {object} utils.Response{data=dto.ProductResponse} "Product updated successfully"
// @Failure 400 {object} utils.Response "Invalid request body or product ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden - Admin access required"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /products/{id} [put]
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

	product, err := s.productService.UpdateProduct(uint(id), &req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update product", err)
		return
	}

	utils.SuccessResponse(c, "Product updated successfully", product)
}

// @Summary Delete a product
// @Description Delete a product (Admin only)
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 200 {object} utils.Response "Product deleted successfully"
// @Failure 400 {object} utils.Response "Invalid product ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden - Admin access required"
// @Failure 404 {object} utils.Response "Product not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /products/{id} [delete]
func (s *Server) deleteProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid product ID", err)
		return
	}

	if err := s.productService.DeleteProduct(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.NotFoundResponse(c, "Product not found")
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to delete product", err)
		return
	}

	utils.SuccessResponse(c, "Product deleted successfully", nil)
}

// @Summary Upload a product image
// @Description Upload an image for a product (Admin only)
// @Tags Products
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param image formData file true "Product image file"
// @Success 201 {object} utils.Response "Product image uploaded successfully"
// @Failure 400 {object} utils.Response "Invalid product ID or image file required"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden - Admin access required"
// @Failure 404 {object} utils.Response "Product not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /products/{id}/images [post]
func (s *Server) uploadProductImage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid product ID", err)
		return
	}

	// Check if product exists before uploading
	_, err = s.productService.GetProduct(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Product not found")
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		utils.BadRequestResponse(c, "Image file is required", err)
		return
	}

	imageURL, err := s.uploadService.UploadProductImage(uint(id), file)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to upload product image", err)
		return
	}

	if err := s.productService.AddProductImage(uint(id), imageURL, file.Filename); err != nil {
		utils.InternalServerErrorResponse(c, "Failed to associate image with product", err)
		return
	}

	utils.CreatedResponse(c, "Product image uploaded successfully", gin.H{"image_url": imageURL})
}
