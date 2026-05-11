package repository

import (
	"context"
	"errors"

	"github.com/chiaf1/iot-nonna-core/internal/domain"
	"github.com/jackc/pgx/v5"
)

// Get a slice of all device_type
func (r *Repository) GetAllDeviceType() ([]domain.Device_type, error) {
	// 1. Create context with timeout for query
	ctx, cancel := context.WithTimeout(context.Background(), r.QueryTimeout_read)
	defer cancel()

	// 2. Create and lounch the query
	rows, err := r.DbPool.Query(ctx, `SELECT 
	id, code, topic, description FROM device_type ORDER BY code`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 3. Scroll all rows and add them to the slice of sensor_type
	var device_types []domain.Device_type
	for rows.Next() {
		var dt domain.Device_type
		// Scan each row for column data
		err := rows.Scan(
			&dt.Id,
			&dt.Code,
			&dt.Topic,
			&dt.Description,
		)
		if err != nil {
			return nil, err
		}

		device_types = append(device_types, dt)
	}

	// Returns the slice of all sensor types
	return device_types, rows.Err()
}

// Get the device_type from id
func (r *Repository) GetDeviceTypeById(id string) (*domain.Device_type, error) {
	// 1. Create context with timeout for query
	ctx, cancel := context.WithTimeout(context.Background(), r.QueryTimeout_read)
	defer cancel()

	// 2. Create and lounch the query
	var dt domain.Device_type
	err := r.DbPool.QueryRow(ctx, `SELECT 
	id, code, topic, description FROM device_type WHERE id = $1`, id).Scan(
		&dt.Id,
		&dt.Code,
		&dt.Topic,
		&dt.Description,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil //not found
	}
	if err != nil {
		return nil, err
	}

	// Returns the slice of all sensor types
	return &dt, err
}

// Create device_type and returns the device_type with id
func (r *Repository) CreateDeviceType(newDt domain.Device_typeRequest) (*domain.Device_type, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.QueryTimeout_read)
	defer cancel()

	var dt domain.Device_type
	err := r.DbPool.QueryRow(ctx, `
	INSERT INTO device_type (
		code,
		topic,
		description
	)
	VALUES ($1,$2,$3)
	RETURNING
		id,
		code,
		topic,
		description
	`,
		newDt.Code,
		newDt.Topic,
		newDt.Description,
	).Scan(
		&dt.Id,
		&dt.Code,
		&dt.Topic,
		&dt.Description,
	)
	if err != nil {
		return nil, err
	}

	return &dt, nil
}

// Update device_type and return the new one
func (r *Repository) UpdateDeviceType(id string, newDt domain.Device_typeRequest) (*domain.Device_type, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.QueryTimeout_read)
	defer cancel()

	var dt domain.Device_type
	err := r.DbPool.QueryRow(ctx, `
	Update device_type 
	SET 
		code = $1,
		topic = $2,
		description = $3
	WHERE id = $4
	RETURNING
		id,
		code,
		topic,
		description
	`,
		newDt.Code,
		newDt.Topic,
		newDt.Description,
		id,
	).Scan(
		&dt.Id,
		&dt.Code,
		&dt.Topic,
		&dt.Description,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil //not found
	}
	if err != nil {
		return nil, err
	}

	return &dt, nil
}

// Delete device_type from id
func (r *Repository) DeleteDeviceType(id string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.QueryTimeout_read)
	defer cancel()

	result, err := r.DbPool.Exec(ctx,
		`DELETE FROM device_type WHERE id = $1`, id,
	)
	if err != nil {
		return false, err
	}

	return result.RowsAffected() > 0, nil
}
