package repository

import (
	"time"

	"gorm.io/gorm"

	entities "golf-store/be-mono/internal/platform/entities"
)

type Repository interface {
	FindUserByIdentifier(identifier string) (*entities.User, error)
	UpdateUserFailedLogins(userID string, attempts int, lockedUntil *time.Time, now time.Time) error
	UpdateUserLoginSuccessTx(userID string, now time.Time, newToken *entities.AuthToken) error
	FindUserByRefreshToken(refreshToken string, userID string, now time.Time) (*entities.User, error)
	RefreshTokensTx(oldRefreshToken string, newToken *entities.AuthToken, now time.Time) error
	DeleteTokenByRefreshToken(refreshToken string) error
	DeleteTokensByAccessToken(accessToken string) error
	ResolveAccessToken(accessToken string, now time.Time) (*entities.User, error)
	FindUserByID(userID string) (*entities.User, error)
}

type GormRepository struct {
	db *gorm.DB
}

func NewGorm(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) FindUserByIdentifier(identifier string) (*entities.User, error) {
	user := &entities.User{}
	err := r.db.
		Where("LOWER(email) = LOWER(?) OR LOWER(phone) = LOWER(?)", identifier, identifier).
		Take(user).Error
	return user, err
}

func (r *GormRepository) UpdateUserFailedLogins(userID string, attempts int, lockedUntil *time.Time, now time.Time) error {
	updates := map[string]any{
		"failed_login_attempts": attempts,
		"updated_at":            now,
		"locked_until":          lockedUntil,
	}
	return r.db.Model(&entities.User{}).Where("id = ?", userID).Updates(updates).Error
}

func (r *GormRepository) UpdateUserLoginSuccessTx(userID string, now time.Time, newToken *entities.AuthToken) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&entities.User{}).
			Where("id = ?", userID).
			Updates(map[string]any{
				"failed_login_attempts": 0,
				"locked_until":          nil,
				"last_login_at":         now,
				"updated_at":            now,
			}).Error; err != nil {
			return err
		}

		if err := tx.Where("user_id = ?", userID).Delete(&entities.AuthToken{}).Error; err != nil {
			return err
		}

		return tx.Create(newToken).Error
	})
}

func (r *GormRepository) FindUserByRefreshToken(refreshToken string, userID string, now time.Time) (*entities.User, error) {
	user := &entities.User{}
	err := r.db.
		Model(&entities.User{}).
		Joins("JOIN auth_tokens t ON t.user_id = users.id").
		Where("t.refresh_token = ? AND t.user_id = ? AND t.expires_at > ?", refreshToken, userID, now).
		Take(user).Error
	return user, err
}

func (r *GormRepository) RefreshTokensTx(oldRefreshToken string, newToken *entities.AuthToken, now time.Time) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("refresh_token = ?", oldRefreshToken).Delete(&entities.AuthToken{}).Error; err != nil {
			return err
		}
		return tx.Create(newToken).Error
	})
}

func (r *GormRepository) DeleteTokenByRefreshToken(refreshToken string) error {
	return r.db.Where("refresh_token = ?", refreshToken).Delete(&entities.AuthToken{}).Error
}

func (r *GormRepository) DeleteTokensByAccessToken(accessToken string) error {
	return r.db.Where("access_token = ?", accessToken).Delete(&entities.AuthToken{}).Error
}

func (r *GormRepository) ResolveAccessToken(accessToken string, now time.Time) (*entities.User, error) {
	var userEntity entities.User
	err := r.db.Model(&entities.User{}).
		Joins("JOIN auth_tokens t ON t.user_id = users.id").
		Where("t.access_token = ? AND t.expires_at > ?", accessToken, now).
		Take(&userEntity).Error
	return &userEntity, err
}

func (r *GormRepository) FindUserByID(userID string) (*entities.User, error) {
	var userEntity entities.User
	err := r.db.Where("id = ?", userID).Take(&userEntity).Error
	return &userEntity, err
}
