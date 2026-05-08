package repository

import "context"

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
