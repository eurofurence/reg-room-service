package database

import "context"

type Repository interface {
	Open(ctx context.Context) error
	Close(ctx context.Context)
	Migrate(ctx context.Context) error
}
