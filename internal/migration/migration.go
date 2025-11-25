package migration

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/vnchk1/subscription-aggregator/internal/config"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type Migrator struct {
	db  *sql.DB
	dir string
}

func NewMigrator(cfg config.DatabaseConfig) (*Migrator, error) {
	connConfig, err := pgx.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	connConfig.TLSConfig = nil

	db := stdlib.OpenDB(*connConfig)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Migrator{
		db:  db,
		dir: cfg.MigrationPath,
	}, nil
}

func (m *Migrator) Up(ctx context.Context) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := os.MkdirAll(m.dir, 0755); err != nil {
		return fmt.Errorf("failed to create migration directory: %w", err)
	}

	if err := goose.UpContext(ctx, m.db, m.dir); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func (m *Migrator) Close() error {
	if m.db != nil {
		return m.db.Close()
	}

	return nil
}
