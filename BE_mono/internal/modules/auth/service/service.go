package service

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	apperrors "golf-store/be-mono/internal/shared/errors"
	"golf-store/be-mono/internal/shared/model"
)

type Service struct {
	db *sql.DB
}

func New(db *sql.DB) *Service {
	return &Service{db: db}
}

type LoginResult struct {
	AccessToken  string
	RefreshToken string
	TokenType    string
	ExpiresIn    int
	User         *model.User
}

func (s *Service) Login(identifier string, password string) (*LoginResult, *apperrors.APIError) {
	now := time.Now().UTC()

	user, err := s.findUserByIdentifier(identifier)
	if err == sql.ErrNoRows {
		return nil, &apperrors.APIError{Status: http.StatusUnauthorized, Code: "INVALID_CREDENTIALS", Message: "Invalid credentials"}
	}
	if err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query user"}
	}

	if user.LockedUntil != nil && user.LockedUntil.After(now) {
		return nil, &apperrors.APIError{Status: http.StatusLocked, Code: "ACCOUNT_LOCKED", Message: "Account is temporarily locked"}
	}

	ok := s.verifyPassword(user, password)
	if !ok {
		attempts := user.FailedLoginAttempts + 1
		var lockedUntil any = nil
		if attempts >= 5 {
			lockTime := now.Add(15 * time.Minute)
			lockedUntil = lockTime
		}
		_, _ = s.db.Exec(
			`UPDATE users SET failed_login_attempts = ?, locked_until = ?, updated_at = ? WHERE id = ?`,
			attempts, lockedUntil, now, user.ID,
		)
		return nil, &apperrors.APIError{Status: http.StatusUnauthorized, Code: "INVALID_CREDENTIALS", Message: "Invalid credentials"}
	}

	accessToken := uuid.NewString()
	refreshToken := uuid.NewString()
	expiresAt := now.Add(24 * time.Hour)

	tx, err := s.db.Begin()
	if err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to begin transaction"}
	}
	defer tx.Rollback()

	if _, err := tx.Exec(
		`UPDATE users
		 SET failed_login_attempts = 0, locked_until = NULL, last_login_at = ?, updated_at = ?
		 WHERE id = ?`,
		now, now, user.ID,
	); err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to update user login state"}
	}

	if _, err := tx.Exec(`DELETE FROM auth_tokens WHERE user_id = ?`, user.ID); err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to clear old auth tokens"}
	}

	if _, err := tx.Exec(
		`INSERT INTO auth_tokens (access_token, refresh_token, user_id, expires_at, created_at)
		 VALUES (?, ?, ?, ?, ?)`,
		accessToken, refreshToken, user.ID, expiresAt, now,
	); err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to issue auth token"}
	}

	if err := tx.Commit(); err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to commit login"}
	}

	user.FailedLoginAttempts = 0
	user.LockedUntil = nil
	user.LastLoginAt = &now
	user.UpdatedAt = now

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		User:         user,
	}, nil
}

func (s *Service) Refresh(refreshToken string) (*LoginResult, *apperrors.APIError) {
	now := time.Now().UTC()

	user, err := s.findUserByRefreshToken(refreshToken)
	if err == sql.ErrNoRows {
		return nil, &apperrors.APIError{Status: http.StatusUnauthorized, Code: "UNAUTHORIZED", Message: "Invalid refresh token"}
	}
	if err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query refresh token"}
	}

	accessToken := uuid.NewString()
	newRefreshToken := uuid.NewString()
	expiresAt := now.Add(24 * time.Hour)

	tx, err := s.db.Begin()
	if err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to begin transaction"}
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`DELETE FROM auth_tokens WHERE refresh_token = ?`, refreshToken); err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to revoke refresh token"}
	}

	if _, err := tx.Exec(
		`INSERT INTO auth_tokens (access_token, refresh_token, user_id, expires_at, created_at)
		 VALUES (?, ?, ?, ?, ?)`,
		accessToken, newRefreshToken, user.ID, expiresAt, now,
	); err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to issue auth token"}
	}

	if err := tx.Commit(); err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to commit token refresh"}
	}

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		User:         user,
	}, nil
}

func (s *Service) Logout(accessToken string) {
	if strings.TrimSpace(accessToken) == "" {
		return
	}
	_, _ = s.db.Exec(`DELETE FROM auth_tokens WHERE access_token = ?`, strings.TrimSpace(accessToken))
}

func (s *Service) ResolveAccessToken(token string) (*model.User, bool) {
	const query = `
SELECT u.id, u.email, u.phone, u.password, u.full_name, u.role, u.status,
       u.failed_login_attempts, u.locked_until, u.last_login_at, u.created_at, u.updated_at
FROM auth_tokens t
JOIN users u ON u.id = t.user_id
WHERE t.access_token = ? AND t.expires_at > UTC_TIMESTAMP(3)
LIMIT 1`

	row := s.db.QueryRow(query, strings.TrimSpace(token))
	user, err := scanUser(row)
	if err != nil {
		return nil, false
	}
	return user, true
}

func (s *Service) Profile(userID string) (*model.User, bool) {
	row := s.db.QueryRow(
		`SELECT id, email, phone, password, full_name, role, status,
		        failed_login_attempts, locked_until, last_login_at, created_at, updated_at
		   FROM users WHERE id = ? LIMIT 1`,
		userID,
	)
	user, err := scanUser(row)
	if err != nil {
		return nil, false
	}
	return user, true
}

func (s *Service) verifyPassword(user *model.User, plain string) bool {
	if user == nil {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(plain))
	if err == nil {
		return true
	}

	return false
}

func (s *Service) findUserByIdentifier(identifier string) (*model.User, error) {
	row := s.db.QueryRow(
		`SELECT id, email, phone, password, full_name, role, status,
		        failed_login_attempts, locked_until, last_login_at, created_at, updated_at
		   FROM users
		  WHERE LOWER(email) = LOWER(?) OR LOWER(phone) = LOWER(?)
		  LIMIT 1`,
		strings.TrimSpace(identifier), strings.TrimSpace(identifier),
	)
	return scanUser(row)
}

func (s *Service) findUserByRefreshToken(refreshToken string) (*model.User, error) {
	row := s.db.QueryRow(
		`SELECT u.id, u.email, u.phone, u.password, u.full_name, u.role, u.status,
		        u.failed_login_attempts, u.locked_until, u.last_login_at, u.created_at, u.updated_at
		   FROM auth_tokens t
		   JOIN users u ON u.id = t.user_id
		  WHERE t.refresh_token = ?
		  LIMIT 1`,
		strings.TrimSpace(refreshToken),
	)
	return scanUser(row)
}

func scanUser(scanner interface {
	Scan(dest ...any) error
}) (*model.User, error) {
	user := &model.User{}
	var role string
	var status string
	var lockedUntil sql.NullTime
	var lastLoginAt sql.NullTime

	err := scanner.Scan(
		&user.ID,
		&user.Email,
		&user.Phone,
		&user.PasswordHash,
		&user.FullName,
		&role,
		&status,
		&user.FailedLoginAttempts,
		&lockedUntil,
		&lastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	user.Role = model.UserRole(role)
	user.Status = model.UserStatus(status)
	if lockedUntil.Valid {
		t := lockedUntil.Time.UTC()
		user.LockedUntil = &t
	}
	if lastLoginAt.Valid {
		t := lastLoginAt.Time.UTC()
		user.LastLoginAt = &t
	}

	return user, nil
}
