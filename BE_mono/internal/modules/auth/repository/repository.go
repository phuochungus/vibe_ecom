package repository

import (
	"time"

	"gorm.io/gorm"

	"golf-store/be-mono/internal/platform/db"
)

type Repository interface {
	FindUserByIdentifier(identifier string) (*db.UserEntity, error)
	UpdateUserFailedLogins(userID string, attempts int, lockedUntil *time.Time, now time.Time) error
	UpdateUserLoginSuccessTx(userID string, now time.Time, newToken *db.AuthTokenEntity) error
	FindUserByRefreshToken(refreshToken string, userID string, now time.Time) (*db.UserEntity, error)
	RefreshTokensTx(oldRefreshToken string, newToken *db.AuthTokenEntity, now time.Time) error
	DeleteTokenByRefreshToken(refreshToken string) error
	DeleteTokensByAccessToken(accessToken string) error
	ResolveAccessToken(accessToken string, now time.Time) (*db.UserEntity, error)
	FindUserByID(userID string) (*db.UserEntity, error)
}

type GormRepository struct {
	db *gorm.DB
}

func NewGorm(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) FindUserByIdentifier(identifier string) (*db.UserEntity, error) {
	user := &db.UserEntity{}
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
	return r.db.Model(&db.UserEntity{}).Where("id = ?", userID).Updates(updates).Error
}

func (r *GormRepository) UpdateUserLoginSuccessTx(userID string, now time.Time, newToken *db.AuthTokenEntity) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&db.UserEntity{}).
			Where("id = ?", userID).
			Updates(map[string]any{
				"failed_login_attempts": 0,
				"locked_until":          nil,
				"last_login_at":         now,
				"updated_at":            now,
			}).Error; err != nil {
			return err
		}

		if err := tx.Where("user_id = ?", userID).Delete(&db.AuthTokenEntity{}).Error; err != nil {
			return err
		}

		return tx.Create(newToken).Error
	})
}

func (r *GormRepository) FindUserByRefreshToken(refreshToken string, userID string, now time.Time) (*db.UserEntity, error) {
	user := &db.UserEntity{}
	err := r.db.
		Model(&db.UserEntity{}).
		Joins("JOIN auth_tokens t ON t.user_id = users.id").
		Where("t.refresh_token = ? AND t.user_id = ? AND t.expires_at > ?", refreshToken, userID, now).
		Take(user).Error
	return user, err
}

func (r *GormRepository) RefreshTokensTx(oldRefreshToken string, newToken *db.AuthTokenEntity, now time.Time) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("refresh_token = ?", oldRefreshToken).Delete(&db.AuthTokenEntity{}).Error; err != nil {
			return err
		}
		return tx.Create(newToken).Error
	})
}

func (r *GormRepository) DeleteTokenByRefreshToken(refreshToken string) error {
	return r.db.Where("refresh_token = ?", refreshToken).Delete(&db.AuthTokenEntity{}).Error
}

func (r *GormRepository) DeleteTokensByAccessToken(accessToken string) error {
	return r.db.Where("access_token = ?", accessToken).Delete(&db.AuthTokenEntity{}).Error
}

func (r *GormRepository) ResolveAccessToken(accessToken string, now time.Time) (*db.UserEntity, error) {
	var userEntity db.UserEntity
	err := r.db.Model(&db.UserEntity{}).
		Joins("JOIN auth_tokens t ON t.user_id = users.id").
		Where("t.access_token = ? AND t.expires_at > ?", accessToken, now).
		Take(&userEntity).Error
	return &userEntity, err
}

func (r *GormRepository) FindUserByID(userID string) (*db.UserEntity, error) {
	var userEntity db.UserEntity
	err := r.db.Where("id = ?", userID).Take(&userEntity).Error
	return &userEntity, err
}
