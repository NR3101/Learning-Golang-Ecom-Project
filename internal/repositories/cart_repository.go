package repositories

import (
	"github.com/NR3101/go-ecom-project/internal/models"
	"gorm.io/gorm"
)

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{
		db: db,
	}
}

func (r *CartRepository) GetByUserID(userID uint) (*models.Cart, error) {
	//TODO implement me
	panic("implement me")
}

func (r *CartRepository) Create(cart *models.Cart) error {
	//TODO implement me
	panic("implement me")
}

func (r *CartRepository) Update(cart *models.Cart) error {
	//TODO implement me
	panic("implement me")
}

func (r *CartRepository) Delete(id uint) error {
	//TODO implement me
	panic("implement me")
}
