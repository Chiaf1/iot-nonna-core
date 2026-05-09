package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type Room struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// Query the database to retrive a list of all rooms
func (r *Repository) GetAllRooms() ([]Room, error) {
	// 1. Create context with timeout for query
	ctx, cancel := context.WithTimeout(context.Background(), r.QueryTimeout_read)
	defer cancel()

	// 2. Create and lounch the query
	rows, err := r.DbPool.Query(ctx, `SELECT id, name FROM rooms ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 3. Scroll all raws and add them to the room slice
	var rooms []Room
	for rows.Next() {
		var (
			id   string
			name string
		)

		// Scan each row for columns data
		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}

		rooms = append(rooms, Room{
			Id:   id,
			Name: name,
		})
	}

	// Return the slice of room and the error of the query if any
	return rooms, rows.Err()
}

// Get single room data by id
func (r *Repository) GetRoomById(id string) (*Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.QueryTimeout_read)
	defer cancel()

	var room Room
	err := r.DbPool.QueryRow(ctx,
		`SELECT id, name FROM rooms WHERE id = $1`, id,
	).Scan(&room.Id, &room.Name)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil // nil nil = not found
	}
	if err != nil {
		return nil, err
	}
	return &room, err
}

// Create room and returns the room with the id created by the db
func (r *Repository) CreateRoom(name string) (*Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.QueryTimeout_read)
	defer cancel()

	var room Room
	err := r.DbPool.QueryRow(ctx,
		`INSERT INTO rooms(name) VALUES ($1) RETURNING id, name`, name,
	).Scan(&room.Id, &room.Name)
	if err != nil {
		return nil, err
	}
	return &room, nil
}

// Update the name and returns the update room
func (r *Repository) UpdateRoom(id, name string) (*Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.QueryTimeout_read)
	defer cancel()

	var room Room
	err := r.DbPool.QueryRow(ctx,
		`UPDATE rooms SET name = $1 WHERE id = $2 RETURNING id, name`, name, id,
	).Scan(&room.Id, &room.Name)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil //not found
	}
	if err != nil {
		return nil, err
	}
	return &room, nil
}

// Elimina la room in base all'id ritorna true se la trova e la elimina
func (r *Repository) DeleteRoom(id string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.QueryTimeout_read)
	defer cancel()

	result, err := r.DbPool.Exec(ctx,
		`DELETE FROM rooms WHERE id = $1`, id,
	)
	if err != nil {
		return false, err
	}

	return result.RowsAffected() > 0, nil
}
