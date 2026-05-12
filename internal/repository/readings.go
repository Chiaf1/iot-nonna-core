package repository

import (
	"context"
	"errors"
	"time"

	"github.com/chiaf1/iot-nonna-core/internal/domain"
	"github.com/jackc/pgx/v5"
)

// Get readgings for dht from to timestamp and limit the number of entries
func (r *Repository) GetDhtReadings(deviceId string, from, to time.Time, limit int) ([]domain.DhtReadings, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.QueryTimeout_read)
	defer cancel()

	// 2. Create and lounch the query
	rows, err := r.DbPool.Query(ctx, `
	SELECT id, device_id, sensor_id, timestamp, temperature, humidity
	FROM dht_readings
	WHERE device_id = $1
		AND timestamp BETWEEN $2 AND $3
	ORDER BY timestamp DESC
	LIMIT $4
	`, deviceId, from, to, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var readings []domain.DhtReadings
	for rows.Next() {
		var rd domain.DhtReadings
		if err := rows.Scan(
			&rd.Id, &rd.DeviceId, &rd.SensorId, &rd.Timestamp, &rd.Temperature, &rd.Humidity,
		); err != nil {
			return nil, err
		}
		readings = append(readings, rd)
	}
	return readings, rows.Err()
}

// Get latest value for dht readings
func (r *Repository) GetLatestDhtReading(deviceId string) (*domain.DhtReadings, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.QueryTimeout_read)
	defer cancel()

	var rd domain.DhtReadings
	err := r.DbPool.QueryRow(ctx, `
	SELECT id, device_id, sensor_id, timestamp, temperature, humidity
	FROM dht_readings
	WHERE device_id = $1
	ORDER BY timestamp DESC
	LIMIT 1
	`, deviceId).Scan(
		&rd.Id, &rd.DeviceId, &rd.SensorId, &rd.Timestamp, &rd.Temperature, &rd.Humidity,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil //not found
	}
	if err != nil {
		return nil, err
	}

	return &rd, nil
}

// Get readgings for dht from to timestamp and limit the number of entries
func (r *Repository) GetStatusReadings(deviceId string, from, to time.Time, limit int) ([]domain.StatusReadings, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.QueryTimeout_read)
	defer cancel()

	// 2. Create and lounch the query
	rows, err := r.DbPool.Query(ctx, `
	SELECT id, device_id, sensor_id, timestamp, status
	FROM device_status
	WHERE device_id = $1
		AND timestamp BETWEEN $2 AND $3
	ORDER BY timestamp DESC
	LIMIT $4
	`, deviceId, from, to, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var readings []domain.StatusReadings
	for rows.Next() {
		var rd domain.StatusReadings
		if err := rows.Scan(
			&rd.Id, &rd.DeviceId, &rd.SensorId, &rd.Timestamp, &rd.Status,
		); err != nil {
			return nil, err
		}
		readings = append(readings, rd)
	}
	return readings, rows.Err()
}

// Get latest value for status readings
func (r *Repository) GetLatestStatusReading(deviceId string) (*domain.StatusReadings, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.QueryTimeout_read)
	defer cancel()

	var rd domain.StatusReadings
	err := r.DbPool.QueryRow(ctx, `
	SELECT id, device_id, sensor_id, timestamp, status
	FROM device_status
	WHERE device_id = $1
	ORDER BY timestamp DESC
	LIMIT 1
	`, deviceId).Scan(
		&rd.Id, &rd.DeviceId, &rd.SensorId, &rd.Timestamp, &rd.Status,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil //not found
	}
	if err != nil {
		return nil, err
	}

	return &rd, nil
}
