package repositories

import (
	"github.com/NR3101/go-ecom-project/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *UserRepository) GetByID(id uint) (*models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *UserRepository) GetByEmailAndIsActive(email string, isActive bool) (*models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *UserRepository) Create(user *models.User) error {
	//TODO implement me
	panic("implement me")
}

func (r *UserRepository) Update(user *models.User) error {
	//TODO implement me
	panic("implement me")
}

func (r *UserRepository) Delete(id uint) error {
	//TODO implement me
	panic("implement me")
}

func (r *UserRepository) CreateRefreshToken(token *models.RefreshToken) error {
	//TODO implement me
	panic("implement me")
}

func (r *UserRepository) GetValidRefreshToken(token string) (*models.RefreshToken, error) {
	//TODO implement me
	panic("implement me")
}

func (r *UserRepository) DeleteRefreshToken(token string) error {
	//TODO implement me
	panic("implement me")
}

func (r *UserRepository) DeleteRefreshTokenByID(id uint) error {
	//TODO implement me
	panic("implement me")
}
