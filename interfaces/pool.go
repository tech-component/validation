package interfaces

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type Pool interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}
