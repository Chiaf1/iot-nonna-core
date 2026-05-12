package repository

import (
	"context"
	"errors"

	"github.com/chiaf1/iot-nonna-core/internal/domain"
	"github.com/jackc/pgx/v5"
)

// Get all devices with nested device type and rooms
func (r *Repository) GetAllDevices() ([]domain.Device, error) {
	// 1. Create context with timeout for query
	ctx, cancel := context.WithTimeout(context.Background(), r.QueryTimeout_read)
	defer cancel()

	// 2. Create and lounch the query
	rows, err := r.DbPool.Query(ctx, `
	SELECT 
		d.id, d.code, d.created_at,
		dt.id, dt.code, dt.topic, dt.description,
		ro.id, ro.name
	FROM devices d
	JOIN device_type dt ON dt.id = d.device_type_id
	LEFT JOIN rooms ro ON ro.id = d.room_id
	ORDER BY d.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 3. Scroll all rows and add them to the slice of sensor_type
	var devices []domain.Device
	for rows.Next() {
		var d domain.Device
		var roomId *string
		var roomName *string
		// Scan each row for column data
		err := rows.Scan(
			&d.Id, &d.Code, &d.CreatedAt,
			&d.DeviceType.Id, &d.DeviceType.Code, &d.DeviceType.Topic, &d.DeviceType.Description,
			&roomId, &roomName,
		)
		if err != nil {
			return nil, err
		}

		// Reconstruct room only if necessary
		if roomId != nil {
			d.Room = &domain.Room{Id: *roomId, Name: *roomName}
		}

		devices = append(devices, d)
	}

	// Returns the slice of all sensor types
	return devices, rows.Err()
}

// Get device by id with nested device type and rooms
func (r *Repository) GetDeviceById(id string) (*domain.Device, error) {
	// 1. Create context with timeout for query
	ctx, cancel := context.WithTimeout(context.Background(), r.QueryTimeout_read)
	defer cancel()

	// 2. Create and lounch the query
	var d domain.Device
	var roomId *string
	var roomName *string
	err := r.DbPool.QueryRow(ctx, `
	SELECT 
		d.id, d.code, d.created_at,
		dt.id, dt.code, dt.topic, dt.description,
		ro.id, ro.name
	FROM devices d
	JOIN device_type dt ON dt.id = d.device_type_id
	LEFT JOIN rooms ro ON ro.id = d.room_id
	WHERE d.id = $1
	`, id).Scan(
		&d.Id, &d.Code, &d.CreatedAt,
		&d.DeviceType.Id, &d.DeviceType.Code, &d.DeviceType.Topic, &d.DeviceType.Description,
		&roomId, &roomName,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil //not found
	}
	if err != nil {
		return nil, err
	}

	// Reconstruct room only if necessary
	if roomId != nil {
		d.Room = &domain.Room{Id: *roomId, Name: *roomName}
	}

	// Returns the slice of all sensor types
	return &d, nil
}

// Create new device
func (r *Repository) CreateDevice(req domain.DeviceRequest) (*domain.Device, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.QueryTimeout_write)
	defer cancel()

	var d domain.Device
	var roomId *string
	var roomName *string
	err := r.DbPool.QueryRow(ctx, `
	WITH inserted AS (
		INSERT INTO devices (code, device_type_id, room_id)
		VALUES ($1, $2, $3)
		RETURNING id, code, device_type_id, room_id, created_at
	)
	SELECT 
		i.id, i.code, i.created_at,
		dt.id, dt.code, dt.topic, dt.description,
		ro.id, ro.name
	FROM inserted i
	JOIN device_type dt ON dt.id = i.device_type_id
	LEFT JOIN rooms ro ON ro.id = i.room_id 
	`,
		req.Code,
		req.DeviceTypeId,
		req.RoomId,
	).Scan(
		&d.Id, &d.Code, &d.CreatedAt,
		&d.DeviceType.Id, &d.DeviceType.Code, &d.DeviceType.Topic, &d.DeviceType.Description,
		&roomId, &roomName,
	)
	if err != nil {
		return nil, err
	}

	// Reconstruct room only if necessary
	if roomId != nil {
		d.Room = &domain.Room{Id: *roomId, Name: *roomName}
	}

	return &d, nil
}

// Update device from id
func (r *Repository) UpdateDevice(id string, req domain.DeviceRequest) (*domain.Device, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.QueryTimeout_write)
	defer cancel()

	var d domain.Device
	var roomId *string
	var roomName *string
	err := r.DbPool.QueryRow(ctx, `
	WITH updated AS (
		UPDATE devices
		SET
			code = $1,
			device_type_id = $2,
			room_id = $3
		WHERE id = $4
		RETURNING id, code, device_type_id, room_id, created_at
	)
	SELECT 
		u.id, u.code, u.created_at,
		dt.id, dt.code, dt.topic, dt.description,
		ro.id, ro.name
	FROM updated u
	JOIN device_type dt ON dt.id = u.device_type_id
	LEFT JOIN rooms ro ON ro.id = u.room_id 
	`,
		req.Code,
		req.DeviceTypeId,
		req.RoomId,
		id,
	).Scan(
		&d.Id, &d.Code, &d.CreatedAt,
		&d.DeviceType.Id, &d.DeviceType.Code, &d.DeviceType.Topic, &d.DeviceType.Description,
		&roomId, &roomName,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil //not found
	}
	if err != nil {
		return nil, err
	}

	// Reconstruct room only if necessary
	if roomId != nil {
		d.Room = &domain.Room{Id: *roomId, Name: *roomName}
	}

	return &d, nil
}

// Delete device from id
func (r *Repository) DeleteDevice(id string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.QueryTimeout_write)
	defer cancel()

	result, err := r.DbPool.Exec(ctx,
		`DELETE FROM devices WHERE id = $1`, id,
	)
	if err != nil {
		return false, err
	}

	return result.RowsAffected() > 0, nil
}
