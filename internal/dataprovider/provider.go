package dataprovider

import (
	"automatedShop/configs"
	errorlib "automatedShop/internal/errors"
	"context"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"time"
)

type Provider struct {
	*sqlx.DB
}

func NewPsqlProvider(config *configs.DbConfig) (*Provider, error) {
	connectionFmt := "postgresql://%s:%s@%s/%s?sslmode=disable"
	db, err := sqlx.Open("pgx", fmt.Sprintf(connectionFmt, config.User, config.Password, config.Host, config.Name))
	if err != nil {
		return nil, fmt.Errorf("failed to add database to pool. Error: %w", errorlib.ErrHttpInternal)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database. Error: %w", err)
	}
	return &Provider{
		DB: db,
	}, nil
}
