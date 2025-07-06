package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"sync"
	"time"
)

type AuthService struct {
	tokens        map[string]*TokenInfo
	refreshTokens map[string]*RefreshTokenInfo
	users         map[string]*UserCredentials
	mutex         sync.RWMutex
	jwtSecret     string
}

type TokenInfo struct {
	UserID    string
	Username  string
	ExpiresAt time.Time
}

type RefreshTokenInfo struct {
	UserID    string
	Username  string
	ExpiresAt time.Time
}

type UserCredentials struct {
	UserID   string
	Username string
	Password string // In production, this should be hashed
}

func NewAuthService(jwtSecret string) *AuthService {
	service := &AuthService{
		tokens:        make(map[string]*TokenInfo),
		refreshTokens: make(map[string]*RefreshTokenInfo),
		users:         make(map[string]*UserCredentials),
		jwtSecret:     jwtSecret,
	}

	// Add some demo users
	service.users["admin"] = &UserCredentials{
		UserID:   "1",
		Username: "admin",
		Password: "password", // In production, use bcrypt
	}
	service.users["user"] = &UserCredentials{
		UserID:   "2",
		Username: "user",
		Password: "password",
	}

	return service
}

func (s *AuthService) Login(ctx context.Context, username, password string) (*struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
	TokenType    string
}, error) {
	s.mutex.RLock()
	user, exists := s.users[username]
	s.mutex.RUnlock()

	if !exists || user.Password != password {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate access token
	accessToken, err := s.generateToken()
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := s.generateToken()
	if err != nil {
		return nil, err
	}

	expiresIn := int64(3600)         // 1 hour
	refreshExpiresIn := int64(86400) // 24 hours

	s.mutex.Lock()
	s.tokens[accessToken] = &TokenInfo{
		UserID:    user.UserID,
		Username:  user.Username,
		ExpiresAt: time.Now().Add(time.Duration(expiresIn) * time.Second),
	}
	s.refreshTokens[refreshToken] = &RefreshTokenInfo{
		UserID:    user.UserID,
		Username:  user.Username,
		ExpiresAt: time.Now().Add(time.Duration(refreshExpiresIn) * time.Second),
	}
	s.mutex.Unlock()

	return &struct {
		AccessToken  string
		RefreshToken string
		ExpiresIn    int64
		TokenType    string
	}{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
		TokenType:    "Bearer",
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.tokens, token)
	return nil
}

func (s *AuthService) ValidateToken(ctx context.Context, token string) (*struct {
	Valid     bool
	UserID    string
	Username  string
	ExpiresAt int64
}, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	tokenInfo, exists := s.tokens[token]
	if !exists {
		return &struct {
			Valid     bool
			UserID    string
			Username  string
			ExpiresAt int64
		}{Valid: false}, nil
	}

	if time.Now().After(tokenInfo.ExpiresAt) {
		return &struct {
			Valid     bool
			UserID    string
			Username  string
			ExpiresAt int64
		}{Valid: false}, nil
	}

	return &struct {
		Valid     bool
		UserID    string
		Username  string
		ExpiresAt int64
	}{
		Valid:     true,
		UserID:    tokenInfo.UserID,
		Username:  tokenInfo.Username,
		ExpiresAt: tokenInfo.ExpiresAt.Unix(),
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
}, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	refreshTokenInfo, exists := s.refreshTokens[refreshToken]
	if !exists {
		return nil, fmt.Errorf("invalid refresh token")
	}

	if time.Now().After(refreshTokenInfo.ExpiresAt) {
		delete(s.refreshTokens, refreshToken)
		return nil, fmt.Errorf("refresh token expired")
	}

	// Generate new access token
	newAccessToken, err := s.generateToken()
	if err != nil {
		return nil, err
	}

	// Generate new refresh token
	newRefreshToken, err := s.generateToken()
	if err != nil {
		return nil, err
	}

	expiresIn := int64(3600)         // 1 hour
	refreshExpiresIn := int64(86400) // 24 hours

	s.tokens[newAccessToken] = &TokenInfo{
		UserID:    refreshTokenInfo.UserID,
		Username:  refreshTokenInfo.Username,
		ExpiresAt: time.Now().Add(time.Duration(expiresIn) * time.Second),
	}

	s.refreshTokens[newRefreshToken] = &RefreshTokenInfo{
		UserID:    refreshTokenInfo.UserID,
		Username:  refreshTokenInfo.Username,
		ExpiresAt: time.Now().Add(time.Duration(refreshExpiresIn) * time.Second),
	}

	// Remove old refresh token
	delete(s.refreshTokens, refreshToken)

	return &struct {
		AccessToken  string
		RefreshToken string
		ExpiresIn    int64
	}{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}

func (s *AuthService) generateToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
