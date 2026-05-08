package repository

import "context"

// Pings the db and returns the error
func (r *Repository) GetDbHealth() error {
	ctx, cancel := context.WithTimeout(context.Background(), r.QueryTimeout_read)
	defer cancel()
	return r.DbPool.Ping(ctx)
}
