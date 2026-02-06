package storage

import (
	"auth-go-skd/internal/domain"
	"context"
)

type UserStorage interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, id string) error
}

type SessionStorage interface {
	CreateSession(ctx context.Context, session *domain.Session) error
	GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*domain.Session, error)
	DeleteSession(ctx context.Context, id string) error
}

type IdentityStorage interface {
	CreateIdentity(ctx context.Context, identity *domain.Identity) error
	GetIdentityByProvider(ctx context.Context, provider, providerID string) (*domain.Identity, error)
}
