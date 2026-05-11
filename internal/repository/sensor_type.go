package repository

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/chiaf1/iot-nonna-core/internal/domain"
	"github.com/jackc/pgx/v5"
)

// Get a list of all sensorsType
func (r *Repository) GetAllSensorType() ([]domain.SensorType, error) {
	// 1. Create context with timeout for query
	ctx, cancel := context.WithTimeout(context.Background(), r.QueryTimeout_read)
	defer cancel()

	// 2. Create and lounch the query
	rows, err := r.DbPool.Query(ctx, `SELECT 
	id, code, topic, description, readings_table_name, column_schema, value_mapping, payload_format, qos_mqtt 
	FROM sensor_type ORDER BY code`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 3. Scroll all rows and add them to the slice of sensor_type
	var Sensor_types []domain.SensorType
	for rows.Next() {
		var st domain.SensorType
		var rawColumnSchema json.RawMessage
		// Scan each row for column data
		err := rows.Scan(
			&st.Id,
			&st.Code,
			&st.Topic,
			&st.Description,
			&st.ReadingsTableName,
			&rawColumnSchema,
			&st.ValueMapping,
			&st.PayloadFormat,
			&st.QosMqtt,
		)
		if err != nil {
			return nil, err
		}

		// Unmarshal rawColumnSchema
		if err := json.Unmarshal(rawColumnSchema, &st.ColumnSchema); err != nil {
			return nil, err
		}

		Sensor_types = append(Sensor_types, st)
	}

	// Returns the slice of all sensor types
	return Sensor_types, rows.Err()
}

// Get the sensor_type based on id
func (r *Repository) GetSensorTypeById(id string) (*domain.SensorType, error) {
	// 1. Create context with timeout for query
	ctx, cancel := context.WithTimeout(context.Background(), r.QueryTimeout_read)
	defer cancel()

	// 2. Create and lounch the query
	var st domain.SensorType
	var rawColumnSchema json.RawMessage
	err := r.DbPool.QueryRow(ctx, `SELECT 
	id, code, topic, description, readings_table_name, column_schema, value_mapping, payload_format, qos_mqtt 
	FROM sensor_type WHERE id = $1`, id).Scan(
		&st.Id,
		&st.Code,
		&st.Topic,
		&st.Description,
		&st.ReadingsTableName,
		&rawColumnSchema,
		&st.ValueMapping,
		&st.PayloadFormat,
		&st.QosMqtt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil //not found
	}
	if err != nil {
		return nil, err
	}

	// Unmarshal rawColumnSchema
	if err := json.Unmarshal(rawColumnSchema, &st.ColumnSchema); err != nil {
		return nil, err
	}

	// Returns the slice of all sensor types
	return &st, err
}

// Create sensortype and returns the sensortype with the id created by the db
func (r *Repository) CreateSensorType(newSt domain.SensorTypeRequest) (*domain.SensorType, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.QueryTimeout_write)
	defer cancel()

	// Marshal column schema to add
	columnSchemaBytes_write, err := json.Marshal(newSt.ColumnSchema)
	if err != nil {
		return nil, err
	}

	var st domain.SensorType
	var columnSchemaBytes_read json.RawMessage
	err = r.DbPool.QueryRow(ctx, `
	INSERT INTO sensor_type (
		code,
		topic,
		description,
		readings_table_name,
		column_schema,
		value_mapping,
		payload_format,
		qos_mqtt
	)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	RETURNING
		id,
		code,
		topic,
		description,
		readings_table_name,
		column_schema,
		value_mapping,
		payload_format,
		qos_mqtt
	`,
		newSt.Code,
		newSt.Topic,
		newSt.Description,
		newSt.ReadingsTableName,
		columnSchemaBytes_write,
		newSt.ValueMapping,
		newSt.PayloadFormat,
		newSt.QosMqtt,
	).Scan(
		&st.Id,
		&st.Code,
		&st.Topic,
		&st.Description,
		&st.ReadingsTableName,
		&columnSchemaBytes_read,
		&st.ValueMapping,
		&st.PayloadFormat,
		&st.QosMqtt,
	)
	if err != nil {
		return nil, err
	}

	// Unmarshal rawColumnSchema
	if err := json.Unmarshal(columnSchemaBytes_read, &st.ColumnSchema); err != nil {
		return nil, err
	}

	return &st, nil
}

// Update the sensor type and returns the updated sensor type
func (r *Repository) UpdateSensorType(id string, newSt domain.SensorTypeRequest) (*domain.SensorType, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.QueryTimeout_write)
	defer cancel()

	// Marshal column schema to add
	columnSchemaBytes_write, err := json.Marshal(newSt.ColumnSchema)
	if err != nil {
		return nil, err
	}

	var st domain.SensorType
	var columnSchemaBytes_read json.RawMessage
	err = r.DbPool.QueryRow(ctx, `
	Update sensor_type 
	SET 
		code = $1,
		topic = $2,
		description = $3,
		readings_table_name = $4,
		column_schema = $5,
		value_mapping = $6,
		payload_format = $7,
		qos_mqtt = $8
	WHERE id = $9
	RETURNING
		id,
		code,
		topic,
		description,
		readings_table_name,
		column_schema,
		value_mapping,
		payload_format,
		qos_mqtt
	`,
		newSt.Code,
		newSt.Topic,
		newSt.Description,
		newSt.ReadingsTableName,
		columnSchemaBytes_write,
		newSt.ValueMapping,
		newSt.PayloadFormat,
		newSt.QosMqtt,
		id,
	).Scan(
		&st.Id,
		&st.Code,
		&st.Topic,
		&st.Description,
		&st.ReadingsTableName,
		&columnSchemaBytes_read,
		&st.ValueMapping,
		&st.PayloadFormat,
		&st.QosMqtt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil //not found
	}
	if err != nil {
		return nil, err
	}

	// Unmarshal rawColumnSchema
	if err := json.Unmarshal(columnSchemaBytes_read, &st.ColumnSchema); err != nil {
		return nil, err
	}

	return &st, nil
}

// Delete the sensor_type based on id returns true if found and eliminated
func (r *Repository) DeleteSensorType(id string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.QueryTimeout_read)
	defer cancel()

	result, err := r.DbPool.Exec(ctx,
		`DELETE FROM sensor_type WHERE id = $1`, id,
	)
	if err != nil {
		return false, err
	}

	return result.RowsAffected() > 0, nil
}
