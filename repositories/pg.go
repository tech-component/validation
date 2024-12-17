package repositories

import (
	"context"
	"fmt"

	"github.com/tech-component/validation/domain"
	"github.com/tech-component/validation/interfaces"
)

type PGRepository struct {
	pool interfaces.Pool
}

func NewPGRepository(pool interfaces.Pool) *PGRepository {
	return &PGRepository{pool: pool}
}

func (pg *PGRepository) CreateUser(ctx context.Context, user domain.User) (string, bool, error) {
	var (
		id string
		ok bool
	)

	query := `SELECT * FROM create_user($1, $2);`

	// Execute the query
	err := pg.pool.QueryRow(ctx, query, user.Email, user.Password).Scan(&id, &ok)
	if err != nil {
		return "", false, fmt.Errorf("failed to create user: %w", err)
	}

	return id, ok, nil
}
