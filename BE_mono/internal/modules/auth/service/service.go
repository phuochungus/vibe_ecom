package service

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"golf-store/be-mono/internal/platform/db"
	apperrors "golf-store/be-mono/internal/shared/errors"
)

type Service struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Service {
	return &Service{db: db}
}

type LoginResult struct {
	AccessToken  string
	RefreshToken string
	TokenType    string
	ExpiresIn    int
	User         *db.UserEntity
}

func (s *Service) Login(identifier string, password string) (*LoginResult, *apperrors.APIError) {
	now := time.Now().UTC()

	userEntity, err := s.findUserByIdentifier(identifier)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, &apperrors.APIError{Status: http.StatusUnauthorized, Code: "INVALID_CREDENTIALS", Message: "Invalid credentials"}
	}
	if err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query user"}
	}

	if userEntity.LockedUntil != nil && userEntity.LockedUntil.After(now) {
		return nil, &apperrors.APIError{Status: http.StatusLocked, Code: "ACCOUNT_LOCKED", Message: "Account is temporarily locked"}
	}

	ok := s.verifyPassword(userEntity, password)
	if !ok {
		attempts := userEntity.FailedLoginAttempts + 1
		updates := map[string]any{
			"failed_login_attempts": attempts,
			"updated_at":            now,
		}
		if attempts >= 5 {
			lockTime := now.Add(15 * time.Minute)
			updates["locked_until"] = lockTime
		} else {
			updates["locked_until"] = nil
		}
		_ = s.db.Model(&db.UserEntity{}).Where("id = ?", userEntity.ID).Updates(updates).Error
		return nil, &apperrors.APIError{Status: http.StatusUnauthorized, Code: "INVALID_CREDENTIALS", Message: "Invalid credentials"}
	}

	accessToken := uuid.NewString()
	refreshToken := uuid.NewString()
	expiresAt := now.Add(24 * time.Hour)

	if err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&db.UserEntity{}).
			Where("id = ?", userEntity.ID).
			Updates(map[string]any{
				"failed_login_attempts": 0,
				"locked_until":          nil,
				"last_login_at":         now,
				"updated_at":            now,
			}).Error; err != nil {
			return err
		}

		if err := tx.Where("user_id = ?", userEntity.ID).Delete(&db.AuthTokenEntity{}).Error; err != nil {
			return err
		}

		return tx.Create(&db.AuthTokenEntity{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			UserID:       userEntity.ID,
			ExpiresAt:    expiresAt,
			CreatedAt:    now,
		}).Error
	}); err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to commit login"}
	}

	userEntity.FailedLoginAttempts = 0
	userEntity.LockedUntil = nil
	userEntity.LastLoginAt = &now
	userEntity.UpdatedAt = now

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		User:         sanitizeUser(userEntity),
	}, nil
}

func (s *Service) Refresh(refreshToken string) (*LoginResult, *apperrors.APIError) {
	now := time.Now().UTC()

	userEntity, err := s.findUserByRefreshToken(refreshToken, now)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, &apperrors.APIError{Status: http.StatusUnauthorized, Code: "UNAUTHORIZED", Message: "Invalid refresh token"}
	}
	if err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query refresh token"}
	}

	accessToken := uuid.NewString()
	newRefreshToken := uuid.NewString()
	expiresAt := now.Add(24 * time.Hour)

	if err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("refresh_token = ?", strings.TrimSpace(refreshToken)).
			Delete(&db.AuthTokenEntity{}).Error; err != nil {
			return err
		}

		return tx.Create(&db.AuthTokenEntity{
			AccessToken:  accessToken,
			RefreshToken: newRefreshToken,
			UserID:       userEntity.ID,
			ExpiresAt:    expiresAt,
			CreatedAt:    now,
		}).Error
	}); err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to commit token refresh"}
	}

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		User:         sanitizeUser(userEntity),
	}, nil
}

func (s *Service) Logout(accessToken string) {
	if strings.TrimSpace(accessToken) == "" {
		return
	}
	_ = s.db.Where("access_token = ?", strings.TrimSpace(accessToken)).Delete(&db.AuthTokenEntity{}).Error
}

func (s *Service) ResolveAccessToken(token string) (*db.UserEntity, bool) {
	var userEntity db.UserEntity
	err := s.db.Model(&db.UserEntity{}).
		Joins("JOIN auth_tokens t ON t.user_id = users.id").
		Where("t.access_token = ? AND t.expires_at > ?", strings.TrimSpace(token), time.Now().UTC()).
		Take(&userEntity).Error
	if err != nil {
		return nil, false
	}
	return sanitizeUser(&userEntity), true
}

func (s *Service) Profile(userID string) (*db.UserEntity, bool) {
	var userEntity db.UserEntity
	if err := s.db.Where("id = ?", userID).Take(&userEntity).Error; err != nil {
		return nil, false
	}
	return sanitizeUser(&userEntity), true
}

func (s *Service) verifyPassword(user *db.UserEntity, plain string) bool {
	if user == nil {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(plain))
	if err == nil {
		return true
	}

	return false
}

func (s *Service) findUserByIdentifier(identifier string) (*db.UserEntity, error) {
	lookup := strings.TrimSpace(identifier)
	user := &db.UserEntity{}
	if err := s.db.
		Where("LOWER(email) = LOWER(?) OR LOWER(phone) = LOWER(?)", lookup, lookup).
		Take(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) findUserByRefreshToken(refreshToken string, now time.Time) (*db.UserEntity, error) {
	user := &db.UserEntity{}
	if err := s.db.
		Model(&db.UserEntity{}).
		Joins("JOIN auth_tokens t ON t.user_id = users.id").
		Where("t.refresh_token = ? AND t.expires_at > ?", strings.TrimSpace(refreshToken), now).
		Take(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func sanitizeUser(user *db.UserEntity) *db.UserEntity {
	if user == nil {
		return nil
	}
	copy := *user
	copy.Password = ""
	if copy.LockedUntil != nil {
		t := copy.LockedUntil.UTC()
		copy.LockedUntil = &t
	}
	if copy.LastLoginAt != nil {
		t := copy.LastLoginAt.UTC()
		copy.LastLoginAt = &t
	}
	copy.CreatedAt = copy.CreatedAt.UTC()
	copy.UpdatedAt = copy.UpdatedAt.UTC()
	return &copy
}
