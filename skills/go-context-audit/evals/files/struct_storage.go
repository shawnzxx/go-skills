package auditfixture

import (
	"context"
	"database/sql"
)

type UserRepository struct {
	ctx context.Context
	db  *sql.DB
}

func NewUserRepository(ctx context.Context, db *sql.DB) *UserRepository {
	return &UserRepository{
		ctx: ctx,
		db:  db,
	}
}

func (r *UserRepository) LoadUser(id string) error {
	queryCtx := context.Background()
	return r.db.QueryRowContext(queryCtx, "SELECT id FROM users WHERE id = $1", id).Err()
}
