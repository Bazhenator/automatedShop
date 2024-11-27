package psql

import (
	"automatedShop/internal/dataprovider"
	customErr "automatedShop/internal/errors"
	"automatedShop/internal/repository/dto"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

const (
	_saveUserQuery   = `INSERT INTO "users"(login, pass_hash) VALUES ($1, $2) RETURNING id`
	_findUserQuery   = `SELECT id, login, pass_hash FROM "users" WHERE login = $1`
	_isRootUserQuery = `SELECT is_admin FROM "users" WHERE id = $1`
)

type AuthProvider struct {
	db *dataprovider.Provider
}

func NewAuthProvider(db *dataprovider.Provider) *AuthProvider {
	return &AuthProvider{db: db}
}

// SaveUser saves new record to users database with given login and password hash.
func (p *AuthProvider) SaveUser(ctx context.Context, login string, pwdHash []byte) error {
	const op = "AuthRepo.SaveUser"

	_, err := p.db.ExecContext(ctx, _saveUserQuery, login, pwdHash)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// FindUser returns user's data if user exists in database. Else returns error.
func (p *AuthProvider) FindUser(ctx context.Context, login string) (*dto.User, error) {
	const op = "AuthRepo.FindUser"

	var user dto.User

	err := p.db.GetContext(ctx, &user, _findUserQuery, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &dto.User{}, fmt.Errorf("%s: %w", op, customErr.ErrUserNotFound)
		}

		return &dto.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

// TODO: исправить mysql на psql нотацию

// IsRoot checks if user is admin.
func (p *AuthProvider) IsRoot(ctx context.Context, uid int64) (bool, error) {
	const op = "AuthRepo.IsRoot"

	stmt, err := p.db.Prepare("SELECT is_admin FROM user WHERE id = ?")
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	var isRoot bool
	row := stmt.QueryRowContext(ctx, uid)

	err = row.Scan(&isRoot)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, customErr.ErrUserNotFound)
		}

		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isRoot, nil
}
