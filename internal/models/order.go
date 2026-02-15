package models

type Order struct {
	ID          uint        `json:"id" gorm:"primary_key"`
	UserID      uint        `json:"user_id" gorm:"not null;index"`
	Status      OrderStatus `json:"status" gorm:"default:'pending'"`
	TotalAmount float64     `json:"total_amount" gorm:"not null"`
	CreatedAt   int64       `json:"created_at"`
	UpdatedAt   int64       `json:"updated_at"`
	DeletedAt   int64       `json:"-" gorm:"index"`

	// Associations
	User       User        `json:"user"`
	OrderItems []OrderItem `json:"order_items"`
}

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type OrderItem struct {
	ID        uint    `json:"id" gorm:"primary_key"`
	OrderID   uint    `json:"order_id" gorm:"not null;index"`
	ProductID uint    `json:"product_id" gorm:"not null;index"`
	Quantity  int     `json:"quantity" gorm:"not null"`
	Price     float64 `json:"price" gorm:"not null"`
	CreatedAt int64   `json:"created_at"`
	UpdatedAt int64   `json:"updated_at"`
	DeletedAt int64   `json:"-" gorm:"index"`

	// Associations
	Order   Order   `json:"-"`
	Product Product `json:"product"`
}

type Cart struct {
	ID        uint  `json:"id" gorm:"primary_key"`
	UserID    uint  `json:"user_id" gorm:"not null;index"`
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
	DeletedAt int64 `json:"-" gorm:"index"`

	// Associations
	CartItems []CartItem `json:"cart_items"`
}

type CartItem struct {
	ID        uint  `json:"id" gorm:"primary_key"`
	CartID    uint  `json:"cart_id" gorm:"not null;index"`
	ProductID uint  `json:"product_id" gorm:"not null;index"`
	Quantity  int   `json:"quantity" gorm:"not null"`
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
	DeletedAt int64 `json:"-" gorm:"index"`

	// Associations
	Cart    Cart    `json:"-"`
	Product Product `json:"product"`
}
