package store

import (
	"auth-go-skd/data"
	"context"
)

type UserStorage interface {
	CreateUser(ctx context.Context, user *data.User) error
	GetUserByEmail(ctx context.Context, email string) (*data.User, error)
	GetUserByID(ctx context.Context, id string) (*data.User, error)
	UpdateUser(ctx context.Context, user *data.User) error
	DeleteUser(ctx context.Context, id string) error
}

type SessionStorage interface {
	CreateSession(ctx context.Context, session *data.Session) error
	GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*data.Session, error)
	DeleteSession(ctx context.Context, id string) error
}

type IdentityStorage interface {
	CreateIdentity(ctx context.Context, identity *data.Identity) error
	GetIdentityByProvider(ctx context.Context, provider, providerID string) (*data.Identity, error)
}
