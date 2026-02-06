package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"auth-go-skd/config"
	"auth-go-skd/internal/domain"
	"auth-go-skd/internal/providers"
	"auth-go-skd/internal/storage"
)

type AuthService struct {
	userStr     storage.UserStorage
	sessionStr  storage.SessionStorage
	identityStr storage.IdentityStorage
	cfg         *config.Config
	providers   map[string]providers.Provider // Generic Map
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Custom Claims
type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func NewAuthService(u storage.UserStorage, s storage.SessionStorage, i storage.IdentityStorage, p map[string]providers.Provider, cfg *config.Config) *AuthService {
	return &AuthService{
		userStr:     u,
		sessionStr:  s,
		identityStr: i,
		providers:   p,
		cfg:         cfg,
	}
}

// ... (Existing methods Register, Login, Refresh, Logout, GetProfile, UpdateProfile, ChangePassword, DeleteAccount remain unchanged) ...

// 12. Social Login (Generic)
func (s *AuthService) GetAuthURL(providerName string) (string, error) {
	p, ok := s.providers[providerName]
	if !ok {
		return "", errors.New("provider not supported")
	}
	return p.GetAuthURL("state-token"), nil // TODO: Generate random state
}

func (s *AuthService) SocialLogin(ctx context.Context, providerName string, code string, ua, ip string) (*TokenPair, error) {
	// 1. Get Provider
	p, ok := s.providers[providerName]
	if !ok {
		return nil, errors.New("provider not supported")
	}

	// 2. Fetch User from Provider
	providerInfo, err := p.FetchUser(ctx, code)
	if err != nil {
		return nil, err
	}

	// 3. Check if Identity exists
	identity, err := s.identityStr.GetIdentityByProvider(ctx, providerName, providerInfo.ID)

	var user *domain.User

	if err == nil {
		// Identity exists, get User
		user, err = s.userStr.GetUserByID(ctx, identity.UserID)
		if err != nil {
			return nil, errors.New("user linked to identity not found")
		}
	} else {
		// Identity does not exist. Check if user with same email exists.
		user, err = s.userStr.GetUserByEmail(ctx, providerInfo.Email)
		if err == nil {
			// User exists, link identity
			// (Pass through to create identity below)
		} else {
			// Create new user
			user = &domain.User{
				ID:         uuid.New().String(),
				Email:      providerInfo.Email,
				Name:       providerInfo.Name,
				Role:       "user",
				IsVerified: true, // Verified by Provider
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}
			if err := s.userStr.CreateUser(ctx, user); err != nil {
				return nil, err
			}
		}

		// Create Identity
		identity = &domain.Identity{
			ID:         uuid.New().String(),
			UserID:     user.ID,
			Provider:   providerName,
			ProviderID: providerInfo.ID,
			CreatedAt:  time.Now(),
			LastLogin:  time.Now(),
		}
		if err := s.identityStr.CreateIdentity(ctx, identity); err != nil {
			return nil, err
		}
	}

	// 4. Login (Generate Tokens)
	return s.generateTokens(ctx, user.ID, user.Role, ua, ip)
}

// 1. Register
func (s *AuthService) Register(ctx context.Context, req domain.RegisterRequest) error {
	if _, err := s.userStr.GetUserByEmail(ctx, req.Email); err == nil {
		return errors.New("user already exists")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &domain.User{
		ID:           uuid.New().String(),
		Email:        req.Email,
		PasswordHash: string(hashed),
		Name:         req.Name,
		Role:         "user",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return s.userStr.CreateUser(ctx, user)
}

// 2. Login
func (s *AuthService) Login(ctx context.Context, req domain.LoginRequest, userAgent, ip string) (*TokenPair, error) {
	user, err := s.userStr.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return s.generateTokens(ctx, user.ID, user.Role, userAgent, ip)
}

// 3. Refresh Token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error) {
	// Verify refresh token from DB
	session, err := s.sessionStr.GetSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if session.ExpiresAt.Before(time.Now()) {
		s.sessionStr.DeleteSession(ctx, session.ID)
		return nil, errors.New("refresh token expired")
	}

	if session.IsBlocked {
		return nil, errors.New("session blocked")
	}

	// Get User
	user, err := s.userStr.GetUserByID(ctx, session.UserID)
	if err != nil {
		return nil, err
	}

	// Rotate tokens (create new session, delete old one)
	// Or just issue new access token. Let's do rotation for better security.
	s.sessionStr.DeleteSession(ctx, session.ID)

	return s.generateTokens(ctx, user.ID, user.Role, session.UserAgent, session.ClientIP)
}

// 4. Logout (simple)
func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	// Find session by refresh token and delete it
	// In a real generic param, we might need ID, but refresh token is unique enough if indexed
	// Our Repo has GetSessionByRefreshToken. We need DeleteSessionByRefreshToken or get ID first.
	session, err := s.sessionStr.GetSessionByRefreshToken(ctx, refreshToken)
	if err == nil {
		return s.sessionStr.DeleteSession(ctx, session.ID)
	}
	return nil
}

// 5. Get Profile
func (s *AuthService) GetProfile(ctx context.Context, userID string) (*domain.User, error) {
	return s.userStr.GetUserByID(ctx, userID)
}

// 6. Update Profile
func (s *AuthService) UpdateProfile(ctx context.Context, userID string, req domain.UpdateProfileRequest) error {
	user, err := s.userStr.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	user.Name = req.Name
	user.UpdatedAt = time.Now()
	return s.userStr.UpdateUser(ctx, user)
}

// 11. Delete Account
func (s *AuthService) DeleteAccount(ctx context.Context, userID string) error {
	return s.userStr.DeleteUser(ctx, userID)
}

// 7. Change Password
func (s *AuthService) ChangePassword(ctx context.Context, userID string, req domain.ChangePasswordRequest) error {
	user, err := s.userStr.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
		return errors.New("invalid old password")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashed)
	user.UpdatedAt = time.Now()

	// Better: Revoke all sessions here!

	return s.userStr.UpdateUser(ctx, user)
}

// Helpers
func (s *AuthService) generateTokens(ctx context.Context, userID, role, ua, ip string) (*TokenPair, error) {
	// Access Token
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)), // Short lived
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte("secret_key_change_me")) // TODO: Use config
	if err != nil {
		return nil, err
	}

	// Refresh Token (Just a random string or long lived JWT)
	refreshToken := uuid.New().String()

	// Save Session
	session := &domain.Session{
		ID:           uuid.New().String(),
		UserID:       userID,
		RefreshToken: refreshToken,
		UserAgent:    ua,
		ClientIP:     ip,
		ExpiresAt:    time.Now().Add(24 * 7 * time.Hour), // 7 days
		CreatedAt:    time.Now(),
	}

	if err := s.sessionStr.CreateSession(ctx, session); err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
