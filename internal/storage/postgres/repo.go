package postgres

import (
	"auth-go-skd/internal/domain"
	"context"
)

// UserStorage implementation

func (p *Postgres) CreateUser(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (id, email, password_hash, name, role, is_verified, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := p.Pool.Exec(ctx, query, user.ID, user.Email, user.PasswordHash, user.Name, user.Role, user.IsVerified, user.CreatedAt, user.UpdatedAt)
	return err
}

func (p *Postgres) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, email, password_hash, name, role, is_verified, created_at, updated_at FROM users WHERE email = $1`
	var user domain.User
	err := p.Pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Role, &user.IsVerified, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (p *Postgres) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	query := `SELECT id, email, password_hash, name, role, is_verified, created_at, updated_at FROM users WHERE id = $1`
	var user domain.User
	err := p.Pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Role, &user.IsVerified, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (p *Postgres) UpdateUser(ctx context.Context, user *domain.User) error {
	query := `UPDATE users SET name=$1, password_hash=$2, updated_at=$3 WHERE id=$4`
	_, err := p.Pool.Exec(ctx, query, user.Name, user.PasswordHash, user.UpdatedAt, user.ID)
	return err
}

func (p *Postgres) DeleteUser(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id=$1`
	_, err := p.Pool.Exec(ctx, query, id)
	return err
}

// SessionStorage implementation

func (p *Postgres) CreateSession(ctx context.Context, session *domain.Session) error {
	query := `INSERT INTO sessions (id, user_id, refresh_token, user_agent, client_ip, is_blocked, expires_at, created_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := p.Pool.Exec(ctx, query,
		session.ID, session.UserID, session.RefreshToken, session.UserAgent, session.ClientIP, session.IsBlocked, session.ExpiresAt, session.CreatedAt)
	return err
}

func (p *Postgres) GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*domain.Session, error) {
	query := `SELECT id, user_id, refresh_token, user_agent, client_ip, is_blocked, expires_at, created_at FROM sessions WHERE refresh_token = $1`
	var s domain.Session
	err := p.Pool.QueryRow(ctx, query, refreshToken).Scan(
		&s.ID, &s.UserID, &s.RefreshToken, &s.UserAgent, &s.ClientIP, &s.IsBlocked, &s.ExpiresAt, &s.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (p *Postgres) DeleteSession(ctx context.Context, id string) error {
	query := `DELETE FROM sessions WHERE id=$1`
	_, err := p.Pool.Exec(ctx, query, id)
	return err
}

// IdentityStorage implementation

func (p *Postgres) CreateIdentity(ctx context.Context, identity *domain.Identity) error {
	query := `INSERT INTO identities (id, user_id, provider, provider_id, created_at, last_login) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := p.Pool.Exec(ctx, query, identity.ID, identity.UserID, identity.Provider, identity.ProviderID, identity.CreatedAt, identity.LastLogin)
	return err
}

func (p *Postgres) GetIdentityByProvider(ctx context.Context, provider, providerID string) (*domain.Identity, error) {
	query := `SELECT id, user_id, provider, provider_id, created_at, last_login FROM identities WHERE provider = $1 AND provider_id = $2`
	var identity domain.Identity
	err := p.Pool.QueryRow(ctx, query, provider, providerID).Scan(&identity.ID, &identity.UserID, &identity.Provider, &identity.ProviderID, &identity.CreatedAt, &identity.LastLogin)
	if err != nil {
		return nil, err
	}
	return &identity, nil
}
