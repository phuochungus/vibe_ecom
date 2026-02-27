package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"golf-store/be-mono/internal/modules/auth/repository"
	"golf-store/be-mono/internal/platform/db"
	apperrors "golf-store/be-mono/internal/shared/errors"
)

const (
	tokenTypeAccess  = "access"
	tokenTypeRefresh = "refresh"
)

type JWTConfig struct {
	Secret     string
	Issuer     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

type Service struct {
	repo       repository.Repository
	secret     []byte
	issuer     string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func New(repo repository.Repository, cfg JWTConfig) *Service {
	secret := strings.TrimSpace(cfg.Secret)
	if secret == "" {
		panic("JWT secret is not configured")
	}
	issuer := strings.TrimSpace(cfg.Issuer)
	if issuer == "" {
		panic("JWT issuer is not configured")
	}
	if cfg.AccessTTL <= 0 {
		panic("JWT access TTL must be positive")
	}
	if cfg.RefreshTTL <= 0 {
		panic("JWT refresh TTL must be positive")
	}

	return &Service{
		repo:       repo,
		secret:     []byte(secret),
		issuer:     issuer,
		accessTTL:  cfg.AccessTTL,
		refreshTTL: cfg.RefreshTTL,
	}
}

type authClaims struct {
	Issuer    string `json:"iss"`
	Subject   string `json:"sub"`
	IssuedAt  int64  `json:"iat"`
	NotBefore int64  `json:"nbf"`
	ExpiresAt int64  `json:"exp"`
	ID        string `json:"jti"`
	TokenType string `json:"token_type"`
	Role      string `json:"role,omitempty"`
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

	userEntity, err := s.repo.FindUserByIdentifier(identifier)
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
		var lockTime *time.Time
		if attempts >= 5 {
			t := now.Add(15 * time.Minute)
			lockTime = &t
		}
		_ = s.repo.UpdateUserFailedLogins(userEntity.ID, attempts, lockTime, now)
		return nil, &apperrors.APIError{Status: http.StatusUnauthorized, Code: "INVALID_CREDENTIALS", Message: "Invalid credentials"}
	}

	accessToken, err := s.signToken(userEntity, tokenTypeAccess, now, s.accessTTL)
	if err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to issue access token"}
	}
	refreshToken, err := s.signToken(userEntity, tokenTypeRefresh, now, s.refreshTTL)
	if err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to issue refresh token"}
	}
	expiresAt := now.Add(s.refreshTTL)

	newToken := &db.AuthTokenEntity{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       userEntity.ID,
		ExpiresAt:    expiresAt,
		CreatedAt:    now,
	}
	if err := s.repo.UpdateUserLoginSuccessTx(userEntity.ID, now, newToken); err != nil {
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
		ExpiresIn:    int(s.accessTTL.Seconds()),
		User:         sanitizeUser(userEntity),
	}, nil
}

func (s *Service) Refresh(refreshToken string) (*LoginResult, *apperrors.APIError) {
	now := time.Now().UTC()
	token := normalizeToken(refreshToken)
	claims, err := s.parseAndValidateJWT(token, tokenTypeRefresh)
	if err != nil || strings.TrimSpace(claims.Subject) == "" {
		return nil, &apperrors.APIError{Status: http.StatusUnauthorized, Code: "UNAUTHORIZED", Message: "Invalid refresh token"}
	}

	userEntity, err := s.findUserByRefreshToken(token, claims.Subject, now)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, &apperrors.APIError{Status: http.StatusUnauthorized, Code: "UNAUTHORIZED", Message: "Invalid refresh token"}
	}
	if err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to query refresh token"}
	}

	newAccessToken, err := s.signToken(userEntity, tokenTypeAccess, now, s.accessTTL)
	if err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to issue access token"}
	}
	newRefreshToken, err := s.signToken(userEntity, tokenTypeRefresh, now, s.refreshTTL)
	if err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to issue refresh token"}
	}
	expiresAt := now.Add(s.refreshTTL)

	newToken := &db.AuthTokenEntity{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		UserID:       userEntity.ID,
		ExpiresAt:    expiresAt,
		CreatedAt:    now,
	}
	if err := s.repo.RefreshTokensTx(token, newToken, now); err != nil {
		return nil, &apperrors.APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Failed to commit token refresh"}
	}

	return &LoginResult{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(s.accessTTL.Seconds()),
		User:         sanitizeUser(userEntity),
	}, nil
}

func (s *Service) Logout(token string) {
	t := normalizeToken(token)
	if t == "" {
		return
	}

	claims, err := s.parseAndValidateJWT(t, "")
	if err != nil {
		return
	}

	switch claims.TokenType {
	case tokenTypeRefresh:
		_ = s.repo.DeleteTokenByRefreshToken(t)
	default:
		_ = s.repo.DeleteTokensByAccessToken(t)
	}
}

func (s *Service) ResolveAccessToken(token string) (*db.UserEntity, bool) {
	t := normalizeToken(token)
	claims, err := s.parseAndValidateJWT(t, tokenTypeAccess)
	if err != nil || strings.TrimSpace(claims.Subject) == "" {
		return nil, false
	}

	userEntity, err := s.repo.ResolveAccessToken(t, time.Now().UTC())
	if err != nil {
		return nil, false
	}
	// ensure Subject matches
	if userEntity.ID != claims.Subject {
		return nil, false
	}
	return sanitizeUser(userEntity), true
}

func (s *Service) Profile(userID string) (*db.UserEntity, bool) {
	userEntity, err := s.repo.FindUserByID(userID)
	if err != nil {
		return nil, false
	}
	return sanitizeUser(userEntity), true
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
	return s.repo.FindUserByIdentifier(lookup)
}

func (s *Service) findUserByRefreshToken(refreshToken string, userID string, now time.Time) (*db.UserEntity, error) {
	return s.repo.FindUserByRefreshToken(refreshToken, userID, now)
}

func (s *Service) signToken(user *db.UserEntity, tokenType string, now time.Time, ttl time.Duration) (string, error) {
	headerJSON, err := json.Marshal(map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	})
	if err != nil {
		return "", err
	}

	claims := authClaims{
		Issuer:    s.issuer,
		Subject:   user.ID,
		IssuedAt:  now.Unix(),
		NotBefore: now.Unix(),
		ExpiresAt: now.Add(ttl).Unix(),
		ID:        uuid.NewString(),
		TokenType: tokenType,
		Role:      user.Role,
	}
	payloadJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	headerPart := base64.RawURLEncoding.EncodeToString(headerJSON)
	payloadPart := base64.RawURLEncoding.EncodeToString(payloadJSON)
	signingInput := headerPart + "." + payloadPart
	signature := signHMACSHA256(signingInput, s.secret)
	return signingInput + "." + base64.RawURLEncoding.EncodeToString(signature), nil
}

func (s *Service) parseAndValidateJWT(token string, expectedType string) (*authClaims, error) {
	if strings.TrimSpace(token) == "" {
		return nil, errors.New("empty token")
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}
	signingInput := parts[0] + "." + parts[1]

	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, errors.New("invalid token signature")
	}
	expectedSig := signHMACSHA256(signingInput, s.secret)
	if !hmac.Equal(signature, expectedSig) {
		return nil, errors.New("invalid token signature")
	}

	headerJSON, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, errors.New("invalid token header")
	}
	var header map[string]any
	if err := json.Unmarshal(headerJSON, &header); err != nil {
		return nil, errors.New("invalid token header")
	}
	if alg, _ := header["alg"].(string); strings.ToUpper(strings.TrimSpace(alg)) != "HS256" {
		return nil, errors.New("unsupported jwt alg")
	}

	payloadJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("invalid token payload")
	}
	claims := &authClaims{}
	if err := json.Unmarshal(payloadJSON, claims); err != nil {
		return nil, errors.New("invalid token payload")
	}

	now := time.Now().UTC().Unix()
	leeway := int64(3)
	if strings.TrimSpace(claims.Issuer) != s.issuer {
		return nil, errors.New("invalid token issuer")
	}
	if claims.Subject == "" {
		return nil, errors.New("invalid token subject")
	}
	if claims.NotBefore > now+leeway {
		return nil, errors.New("token not active")
	}
	if claims.IssuedAt > now+leeway {
		return nil, errors.New("token issued in the future")
	}
	if claims.ExpiresAt <= now-leeway {
		return nil, errors.New("token expired")
	}
	if expectedType != "" && claims.TokenType != expectedType {
		return nil, errors.New("invalid token type")
	}

	return claims, nil
}

func signHMACSHA256(input string, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write([]byte(input))
	return mac.Sum(nil)
}

func normalizeToken(token string) string {
	trimmed := strings.TrimSpace(token)
	if trimmed == "" {
		return ""
	}
	parts := strings.SplitN(trimmed, " ", 2)
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return strings.TrimSpace(parts[1])
	}
	return trimmed
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
