package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/smart-safety-hub/backend/shared"
)

type Ext interface {
	sqlx.Queryer
	sqlx.Execer
	GetContext(ctx context.Context, dest interface{}, query string, agrs ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type UserRepo struct {
	db Ext
}

func NewUserRepo(db Ext) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) ExecuteTransaction(ctx context.Context, fn func(repo *UserRepo) error) error {
	db, ok := r.db.(*sqlx.DB)
	if !ok {
		return fmt.Errorf("Transaction already in progress or invalid db type")
	}

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	txRepo := &UserRepo{db: tx}
	err = fn(txRepo)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

func (r *UserRepo) SaveUser(ctx context.Context, u *User) (*User, error) {
	var user User
	query := "INSERT INTO users(full_name, email, password, phone_number) VALUES ($1,$2,$3,$4) RETURNING id, full_name, email, phone_number, created_at, updated_at"
	if err := r.db.GetContext(ctx, &user, query, u.FullName, u.Email, u.Password, u.PhoneNumber); err != nil {
		return nil, shared.PostgresError(err)
	}
	return &user, nil
}

func (r *UserRepo) UpdatePassword(ctx context.Context, password, userId string) error {
	query := "Update users SET password = $1, updated_at = NOW() WHERE id = $2"
	if _, err := r.db.ExecContext(ctx, query, password, userId); err != nil {
		return shared.PostgresError(err)
	}

	return nil
}

func (r *UserRepo) GetUser(ctx context.Context, email string) (*User, error) {
	var user User
	query := "SELECT id, full_name, email, password, phone_number, created_at, updated_at FROM users WHERE email=$1"
	if err := r.db.GetContext(ctx, &user, query, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.ErrUserNotFound
		}
		return nil, shared.PostgresError(err)
	}
	return &user, nil
}

func (r *UserRepo) GetUserById(ctx context.Context, id string) (*User, error) {
	var user User
	query := "SELECT id, full_name, email, password, phone_number, created_at, updated_at FROM users WHERE id=$1"
	if err := r.db.GetContext(ctx, &user, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.ErrUserNotFound
		}
		return nil, shared.PostgresError(err)
	}
	return &user, nil
}

func (r *UserRepo) GetRole(ctx context.Context, userType string) (*Roles, error) {
	var roles Roles
	query := "SELECT * FROM roles WHERE name=$1"
	if err := r.db.GetContext(ctx, &roles, query, userType); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.ErrUserNotFound
		}
		return nil, shared.PostgresError(err)
	}
	return &roles, nil
}

func (r *UserRepo) SaveUserRoles(ctx context.Context, userId, RoleId string) error {
	query := "INSERT INTO user_roles(user_id, role_id) VALUES ($1,$2) ON CONFLICT DO NOTHING"
	if _, err := r.db.ExecContext(ctx, query, userId, RoleId); err != nil {
		return shared.PostgresError(err)
	}
	return nil
}

func (r *UserRepo) GetUserRoles(ctx context.Context, userId string) (*UserRoles, error) {
	var userRoles UserRoles
	query := "SELECT * FROM user_roles WHERE user_id=$1"
	if err := r.db.GetContext(ctx, &userRoles, query, userId); err != nil {
		return nil, shared.PostgresError(err)
	}
	return &userRoles, nil
}

func (r *UserRepo) GetRolesPermissions(ctx context.Context, roleId string) (*RolesPermissions, error) {
	var result RolesPermissions
	// query := `SELECT r.name AS role, p.name AS permission FROM roles r JOIN roles_permissions rp ON rp.role_id = r.id JOIN permissions p ON rp.permission_id = p.id WHERE r.id=$1`
	query := `SELECT r.name AS role, ARRAY_AGG(p.name) AS permissions FROM roles r JOIN roles_permissions rp ON rp.role_id = r.id JOIN permissions p ON rp.permission_id = p.id WHERE r.id = $1 GROUP BY r.name`
	if err := r.db.GetContext(ctx, &result, query, roleId); err != nil {
		return nil, shared.PostgresError(err)
	}
	return &result, nil
}

func (r *UserRepo) SaveRefreshToken(ctx context.Context, userId, refreshToken string, ttl time.Duration) error {
	query := "INSERT INTO refresh_tokens(user_id, token, expires_at) VALUES($1,$2,$3)"
	expiresAt := time.Now().UTC().Add(ttl)
	if _, err := r.db.ExecContext(ctx, query, userId, refreshToken, expiresAt); err != nil {
		return shared.PostgresError(err)
	}
	return nil
}

func (r *UserRepo) UpdateRefreshToken(ctx context.Context, token string) error {
	query := "UPDATE refresh_tokens SET revoked=TRUE WHERE token=$1"
	if _, err := r.db.ExecContext(ctx, query, token); err != nil {
		return shared.PostgresError(err)
	}
	return nil
}

func (r *UserRepo) RevokeRefreshToken(ctx context.Context, userId string) error {
	query := "UPDATE refresh_tokens SET revoked=TRUE WHERE user_id=$1"
	if _, err := r.db.ExecContext(ctx, query, userId); err != nil {
		return shared.PostgresError(err)
	}
	return nil
}

func (r *UserRepo) GetRefreshToken(ctx context.Context, token string) (*RefreshTokens, error) {
	var refreshTokens RefreshTokens
	query := "SELECT * FROM refresh_tokens WHERE token=$1 AND revoked=false"
	if err := r.db.GetContext(ctx, &refreshTokens, query, token); err != nil {
		return nil, shared.PostgresError(err)
	}
	return &refreshTokens, nil
}
