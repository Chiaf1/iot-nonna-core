package db

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/chiaf1/iot-nonna-core/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Open pool for database connection
func OpenPool(dsn string, queryTimeout time.Duration) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	// loading configs drom connection url
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	// Adjusting some configs
	cfg.MaxConns = 10
	cfg.MinConns = 2
	cfg.MaxConnLifetime = 30 * time.Minute
	cfg.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	// Connection test
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}

func NewPgPool(cfg config.DbConfig) (*pgxpool.Pool, error) {
	delay := cfg.ConnectionInterval

	for attempt := 1; cfg.MaxRetry == 0 || attempt <= cfg.MaxRetry; attempt++ {
		dbPool, err := OpenPool(cfg.DbURL, cfg.Query_timeout_read)
		if err == nil {
			return dbPool, nil
		}
		log.Printf("[DB] Error while opening connection: %v \n", err)

		// Exponential backoff + jitter
		jitter := time.Duration(rand.Int63n(int64(delay / 2)))
		wait := delay + jitter

		if wait > cfg.MaxDelay {
			wait = cfg.MaxDelay
		}

		log.Printf("[DB] Waiting %v before connection retry...", wait)
		time.Sleep(wait)

		delay *= 2

		if delay > cfg.MaxDelay {
			delay = cfg.MaxDelay
		}
	}

	return nil, fmt.Errorf("max connection retry reached (%d)", cfg.MaxRetry)
}
