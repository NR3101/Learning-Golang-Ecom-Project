package models

type Category struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active" gorm:"default:true"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
	DeletedAt   int64  `json:"-" gorm:"index"`

	// Associations
	Products []Product `json:"-"`
}

type Product struct {
	ID          uint    `json:"id" gorm:"primary_key"`
	Name        string  `json:"name" gorm:"not null"`
	Description string  `json:"description"`
	Price       float64 `json:"price" gorm:"not null"`
	Stock       int     `json:"stock" gorm:"default:0"`
	SKU         string  `json:"sku" gorm:"uniqueIndex;not null"`
	IsActive    bool    `json:"is_active" gorm:"default:true"`
	CreatedAt   int64   `json:"created_at"`
	UpdatedAt   int64   `json:"updated_at"`
	DeletedAt   int64   `json:"-" gorm:"index"`

	// Associations
	CategoryID uint           `json:"category_id" gorm:"not null;index"`
	Category   Category       `json:"category"`
	Images     []ProductImage `json:"-"`
	OrderItems []OrderItem    `json:"-"`
	CartItems  []CartItem     `json:"-"`
}

type ProductImage struct {
	ID        uint   `json:"id" gorm:"primary_key"`
	ProductID uint   `json:"product_id" gorm:"not null;index"`
	URL       string `json:"url" gorm:"not null"`
	AltText   string `json:"alt_text"`
	IsPrimary bool   `json:"is_primary" gorm:"default:false"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	DeletedAt int64  `json:"-" gorm:"index"`

	// Associations
	Product Product `json:"-"`
}
