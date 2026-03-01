package repository

import (
	db "golf-store/be-mono/internal/platform/db"

	"gorm.io/gorm"
)

type Repository interface {
	GetUserByID(id string) (*db.UserEntity, error)
}

type GormRepository struct {
	db *gorm.DB
}

func NewGorm(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) GetUserByID(id string) (*db.UserEntity, error) {
	var user db.UserEntity
	// Exclude password_hash for security reasons
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	user.Password = "masked" // Mask the password hash before returning
	return &user, nil
}
