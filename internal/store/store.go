package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Store wraps the database connection pool.
type Store struct {
	pool *pgxpool.Pool
}

// New opens a connection pool and checks the database answers.
func New(ctx context.Context, databaseURL string) (*Store, error) {
	pool, err := pgxpool.New(ctx, databaseURL)

	if err != nil {
		return nil, fmt.Errorf("create database pool: %w", err)
	}

	s := &Store{pool: pool}
	if err := s.Ping(ctx); err != nil {
		s.Close()
		return nil, err
	}

	return s, nil
}

// Close closes the connection pool.
func (s *Store) Close() {
	s.pool.Close()
}

// Ping checks the database answers.
func (s *Store) Ping(ctx context.Context) error {
	if err := s.pool.Ping(ctx); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	return nil
}
