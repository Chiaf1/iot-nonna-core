package repository

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	DbPool             *pgxpool.Pool
	QueryTimeout_read  time.Duration
	QueryTimeout_write time.Duration
}

func NewRepo(dbPool *pgxpool.Pool, queryTimeout_r, queryTimeout_w time.Duration) *Repository {
	return &Repository{
		DbPool:             dbPool,
		QueryTimeout_read:  queryTimeout_r,
		QueryTimeout_write: queryTimeout_w,
	}
}
