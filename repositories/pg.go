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

// buildQuery executes sql query with variable number of arguments.
// sql should return (string, boolean) otherwise an error will return.
func (pg *PGRepository) buildQuery(ctx context.Context, query string, args ...interface{}) (string, bool, error) {
	var (
		id string
		ok bool
	)

	err := pg.pool.QueryRow(ctx, query, args...).Scan(&id, &ok)
	if err != nil {
		return "", false, fmt.Errorf("failed to create entity: %w", err)
	}

	return id, ok, nil
}

// CreateUser creates a new user using the generic buildQuery function
func (pg *PGRepository) CreateUser(ctx context.Context, user domain.User) (string, bool, error) {
	query := `SELECT * FROM create_user($1, $2);`
	return pg.buildQuery(ctx, query, user.Email, user.Password)
}
