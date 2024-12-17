package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mauleyzaola/validation/domain"
)

type PostgresDB struct {
	pool *pgxpool.Pool
}

func NewPostgresDB(ctx context.Context, connectionString string) (*PostgresDB, error) {
	pool, err := pgxpool.Connect(ctx, connectionString)
	if err != nil {
		return nil, err
	}
	return &PostgresDB{pool: pool}, nil
}

func (db *PostgresDB) Close() {
	db.pool.Close()
}

func (db *PostgresDB) CreateUser(ctx context.Context, user domain.User) (string, bool, error) {
	var (
		id string
		ok bool
	)

	query := `SELECT * FROM create_user($1, $2);`

	// Execute the query
	err := db.pool.QueryRow(ctx, query, user.Email, user.Password).Scan(&id, &ok)
	if err != nil {
		return "", false, fmt.Errorf("failed to create user: %w", err)
	}

	return id, ok, nil
}
