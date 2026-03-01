package repository

import (
	"gorm.io/gorm"

	entities "golf-store/be-mono/internal/platform/entities"
)

type Repository interface {
	GetUserByID(id string) (*entities.User, error)
}

type GormRepository struct {
	db *gorm.DB
}

func NewGorm(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) GetUserByID(id string) (*entities.User, error) {
	var user entities.User
	// Exclude password_hash for security reasons
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	user.Password = "masked" // Mask the password hash before returning
	return &user, nil
}
